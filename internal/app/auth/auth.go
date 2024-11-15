package auth

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/xEgorka/project3/internal/app/config"
	"github.com/xEgorka/project3/internal/app/crypto"
	"github.com/xEgorka/project3/internal/app/logger"
)

// Auth describes auth interface.
type Auth interface {
	Register(ctx context.Context, usr, pass string) error
	Login(ctx context.Context, usr, pass string) error
}

type db struct {
	conn *sql.DB
	cfg  *config.Config
	c    crypto.Crypto
}

func new(config *config.Config, conn *sql.DB, crypto crypto.Crypto) *db {
	return &db{cfg: config, conn: conn, c: crypto}
}

// Open opens new auth database storage.
func Open(ctx context.Context, cfg *config.Config, c crypto.Crypto) (Auth, error) {
	logger.Log.Info("opening auth database...", zap.String("conninfo", cfg.DBURI))
	conn, err := sql.Open(cfg.DBDriver, cfg.DBURI)
	if err != nil {
		return nil, err
	}
	if e := conn.Ping(); e != nil {
		return nil, e
	}
	return open(ctx, cfg, c, conn)
}

func open(ctx context.Context, cfg *config.Config, c crypto.Crypto,
	conn *sql.DB) (Auth, error) {
	db := new(cfg, conn, c)
	if err := db.bootstrap(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

const queryBootstrap = `
		create table if not exists users (
			usr     varchar not null,
			hash    varchar not null,
			updated timestamp(0) not null default now(),
			primary key (usr))
`

func (a *db) bootstrap(ctx context.Context) error {
	logger.Log.Info("bootstrapping auth database...")
	if _, err := a.conn.ExecContext(ctx, queryBootstrap); err != nil {
		logger.Log.Info("auth database bootstrap error", zap.Error(err))
		return err
	} else {
		logger.Log.Info("bootstrapping auth database... done")
		return nil
	}
}

const queryRegister = `insert into users (usr, hash) values ($1, $2)`

// Register signs new user up.
func (a *db) Register(ctx context.Context, usr, pass string) error {
	h, err := a.c.Hash(pass)
	if err != nil {
		return err
	}
	if _, err := a.conn.ExecContext(ctx, queryRegister, usr, h); err != nil {
		return err
	}
	return nil
}

const queryLogin = `select hash from users where usr = $1`

// Login signs user in.
func (a *db) Login(ctx context.Context, usr, pass string) error {
	var hash string
	row := a.conn.QueryRowContext(ctx, queryLogin, usr)
	if err := row.Scan(&hash); err != nil {
		return err
	}
	if err := a.c.Verify(pass, hash); err != nil {
		return err
	}
	return nil
}
