// Package storage implements Storage interface.
package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/xEgorka/project3/internal/app/config"
	"github.com/xEgorka/project3/internal/app/logger"
	"github.com/xEgorka/project3/internal/app/models"
)

// Storage describes methods required to implement Storage.
type Storage interface {
	GetUserData(ctx context.Context, usr string, timestamp int) (
		[]models.UserData, error)
	MergeUserData(ctx context.Context, d []models.UserData) error
	Ping() error
	Close() error
}

// Open initializes Storage.
func Open(ctx context.Context, cfg *config.Config) (Storage, error) {
	logger.Log.Info("opening database...", zap.String("conninfo", cfg.DBURI))
	conn, err := sql.Open(cfg.DBDriver, cfg.DBURI)
	if err != nil {
		return nil, err
	}
	if e := conn.Ping(); e != nil {
		return nil, e
	}
	return open(ctx, cfg, conn)
}

func open(ctx context.Context, cfg *config.Config, conn *sql.DB) (Storage, error) {
	db := new(cfg, conn)
	if err := db.bootstrap(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

type db struct {
	conn *sql.DB
	cfg  *config.Config
}

func new(config *config.Config, conn *sql.DB) *db {
	return &db{cfg: config, conn: conn}
}

const (
	queryBootstrap1 = `
		create table if not exists data (
			usr   	varchar not null,
			id   	varchar not null,
			type   	smallint not null,
			data  	varchar not null,
			updated timestamp(0) not null,
			merged  timestamp(0) not null default now())
`
	queryBootstrap2 = `create unique index if not exists data_idx on data (usr, id)`
)

func (s *db) bootstrap(ctx context.Context) error {
	logger.Log.Info("bootstrapping database...")
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, queryBootstrap1); err != nil {
		if err = tx.Rollback(); err != nil {
			logger.Log.Info("failed tx rollback", zap.Error(err))
			return err
		}
		logger.Log.Info("database bootstrap error", zap.Error(err))
		return err
	}
	if _, err = tx.ExecContext(ctx, queryBootstrap2); err != nil {
		if err = tx.Rollback(); err != nil {
			logger.Log.Info("failed tx rollback", zap.Error(err))
			return err
		}
		logger.Log.Info("database bootstrap error", zap.Error(err))
		return err
	}
	logger.Log.Info("bootstrapping database... done")
	return tx.Commit()
}

const queryGetUserData = `
		select usr, id, type, data, updated, merged
		from data where usr = $1 and merged > to_timestamp($2)
`

// GetUserData returns user data merged after specified timestamp.
func (s *db) GetUserData(ctx context.Context, usr string, timestamp int) (
	[]models.UserData, error) {
	rows, err := s.conn.QueryContext(ctx, queryGetUserData, usr, timestamp)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger.Log.Error("failed close rows", zap.Error(err))
		}
	}()

	var d []models.UserData
	for rows.Next() {
		var usr, id, data, updatedStr, mergedStr string
		var t int
		if err = rows.Scan(&usr, &id, &t, &data,
			&updatedStr, &mergedStr); err != nil {
			return nil, err
		}
		updated, e := time.Parse(time.RFC3339, updatedStr)
		if e != nil {
			return nil, e
		}
		merged, ee := time.Parse(time.RFC3339, mergedStr)
		if ee != nil {
			return nil, ee
		}
		dd := models.UserData{
			UserID: usr, ID: id, Type: t, Data: data, Updated: updated, Merged: merged}
		d = append(d, dd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return d, nil
}

const queryMergeUserData = `
			insert into data (usr, id, type, data, updated)
			values ($1, $2, $3, $4, $5)
			on conflict(usr, id)
			do update set
				type = excluded.type,
				data = excluded.data,
				updated = excluded.updated,
				merged = now()
			where data.updated < excluded.updated
`

// MergeUserData merges list of user data.
func (s *db) MergeUserData(ctx context.Context, d []models.UserData) error {
	for _, v := range d {
		if _, err := s.conn.ExecContext(ctx, queryMergeUserData,
			v.UserID, v.ID, v.Type, v.Data, v.Updated); err != nil {
			return err
		}
	}
	return nil
}

// Ping checks database connection.
func (s *db) Ping() error { return s.conn.Ping() }

// Close closes database connection.
func (s *db) Close() error { return s.conn.Close() }
