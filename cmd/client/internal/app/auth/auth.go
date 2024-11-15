// Package auth provides interface to authentication database.
package auth

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/crypto"
	"github.com/xEgorka/project3/cmd/client/internal/app/logger"
)

// Auth describes auth interface.
type Auth interface {
	LoginOnline(ctx context.Context, usr, pass string) error
	LoginOffline(ctx context.Context, usr, pass string) error
}

type db struct {
	c    crypto.Crypto
	conn *sql.DB
	cfg  *config.Config
}

func new(config *config.Config, conn *sql.DB, crypto crypto.Crypto) *db {
	return &db{cfg: config, conn: conn, c: crypto}
}

// Open opens new auth database storage.
func Open(ctx context.Context, cfg *config.Config, c crypto.Crypto) (Auth, error) {
	logger.Log.Info("opening auth database...", zap.String("conninfo", cfg.DBPath))
	conn, err := sql.Open(cfg.DBDriver, cfg.DBPath)
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
			updated integer(4) not null default (strftime('%s','now')),
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

const queryLoginOnline = `
		insert into users (usr, hash) values ($1, $2)
		on conflict(usr)
		do update set
			hash = excluded.hash,
			updated = strftime('%s','now')
`

// LoginOnline signs user in online.
func (a *db) LoginOnline(ctx context.Context, usr, pass string) error {
	a.cfg.Usr = usr
	a.cfg.Key = a.c.Key(usr, pass)
	hash, err := a.c.Enc(a.cfg.Usr, a.cfg.Key)
	if err != nil {
		return err
	}
	if _, err = a.conn.ExecContext(ctx, queryLoginOnline, a.cfg.Usr, hash); err != nil {
		return err
	}
	return nil
}

const delay = 500 * time.Millisecond

const queryLoginOffline = `select hash from users where usr = $1`

// LoginOffline signs user in offline.
func (a *db) LoginOffline(ctx context.Context, usr, pass string) error {
	a.cfg.Usr = usr
	a.cfg.Key = a.c.Key(usr, pass)
	var hash string
	row := a.conn.QueryRowContext(ctx, queryLoginOffline, usr)
	if err := row.Scan(&hash); err != nil {
		return err
	}
	time.Sleep(delay) // prevents brute-force pass offline
	if _, err := a.c.Dec(hash, a.cfg.Key); err != nil {
		return err
	}
	return nil
}
