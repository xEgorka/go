// Package storage implements Storage interface.
package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/logger"
	"github.com/xEgorka/project3/cmd/client/internal/app/models"
)

// Storage describes methods required to implement storage.
type Storage interface {
	GetTimestamp(ctx context.Context) (string, error)
	UpdateUserData(ctx context.Context, d []models.UserData) error
	GetUnmergedUserData(ctx context.Context, usr string) ([]models.UserData, error)
	IndexUserData(ctx context.Context) ([]models.UserData, error)
	PrintUserData(ctx context.Context, id string) (*models.UserData, error)
	EnterUserData(ctx context.Context, d *models.UserData) error
}

// Open initializes Storage.
func Open(ctx context.Context, cfg *config.Config) (Storage, error) {
	logger.Log.Info("opening database...", zap.String("conninfo", cfg.DBPath))
	conn, err := sql.Open(cfg.DBDriver, cfg.DBPath)
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
		usr     text not null,
		id      text not null,
		type    integer(1) not null,
		data    text not null,
		updated integer(4) not null default (strftime('%s','now')),
		merged  integer(4) default null)
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
		logger.Log.Info("database bootstrap error", zap.Error(err))
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if _, err = tx.ExecContext(ctx, queryBootstrap2); err != nil {
		logger.Log.Info("database bootstrap error", zap.Error(err))
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	logger.Log.Info("bootstrapping database... done")
	return tx.Commit()
}

const queryGetTimestamp = `
		with dummy as (select ''),
		t as (select max(merged) ts from data where usr = $1)
		select case when t.ts is null then 0 else t.ts end ts
		from t cross join dummy
`

// GetTimestamp returns maximum user data timestamp.
func (s *db) GetTimestamp(ctx context.Context) (string, error) {
	var t string
	row := s.conn.QueryRowContext(ctx, queryGetTimestamp, s.cfg.Usr)
	if err := row.Scan(&t); err != nil {
		return "", err
	}
	return t, nil
}

const queryUpdateUserData = `
		insert into data (usr, id, type, data, updated, merged)
			values ($1, $2, $3, $4, $5, $6)
			on conflict(usr, id)
			do update set
				type = excluded.type,
				data = excluded.data,
				updated = excluded.updated,
				merged = excluded.merged
`

// UpdateUserData updates list of user data.
func (s *db) UpdateUserData(ctx context.Context, d []models.UserData) error {
	for _, v := range d {
		_, err := s.conn.ExecContext(ctx, queryUpdateUserData,
			v.UserID, v.ID, v.Type, v.Data, v.Updated.Unix(), v.Merged.Unix())
		if err != nil {
			return err
		}
	}
	return nil
}

const queryGetUnmergedUserData = `
		select usr, id, type, data,
			strftime('%Y-%m-%dT%H:%M:%SZ', datetime(updated, 'unixepoch')) updated
		from data where merged is null and usr = $1
`

// GetUnmergedUserData returns list of unmerged user data.
func (s *db) GetUnmergedUserData(ctx context.Context, usr string) (
	[]models.UserData, error) {
	var req []models.UserData
	rows, err := s.conn.QueryContext(ctx, queryGetUnmergedUserData, usr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Log.Info("failed rows close", zap.Error(err))
		}
	}()

	for rows.Next() {
		var usr, id, data, updatedStr string
		var t int
		if e := rows.Scan(&usr, &id, &t, &data, &updatedStr); e != nil {
			return nil, e
		}
		updated, ee := time.Parse(time.RFC3339, updatedStr)
		if ee != nil {
			return nil, ee
		}
		r := models.UserData{
			UserID: usr, ID: id, Type: t, Data: data, Updated: updated}
		req = append(req, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return req, nil
}

const queryIndexUserData = `
		select id,
				case
					when type = 1 then 'PASS'
					when type = 2 then 'TEXT'
					when type = 3 then 'FILE'
					when type = 4 then 'CARD'
				end as type,
			strftime('%Y-%m-%dT%H:%M:%SZ', datetime(updated, 'unixepoch')) updated
		from data where usr = $1 order by updated asc
`

// IndexUserData returns index of user data.
func (s *db) IndexUserData(ctx context.Context) ([]models.UserData, error) {
	var d []models.UserData
	rows, err := s.conn.QueryContext(ctx, queryIndexUserData, s.cfg.Usr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Log.Info("failed rows close", zap.Error(err))
		}
	}()

	for rows.Next() {
		var id, t, updatedStr string
		if e := rows.Scan(&id, &t, &updatedStr); e != nil {
			return nil, e
		}
		updated, ee := time.Parse(time.RFC3339, updatedStr)
		if ee != nil {
			return nil, ee
		}
		dd := models.UserData{ID: id, Data: t, Updated: updated}
		d = append(d, dd)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return d, nil
}

const queryPrintUserData = `select type, data from data where usr = $1 and id = $2`

// PrintUserData prints user data.
func (s *db) PrintUserData(ctx context.Context, id string) (*models.UserData, error) {
	row := s.conn.QueryRowContext(ctx, queryPrintUserData, s.cfg.Usr, id)
	var t int
	var d string
	if err := row.Scan(&t, &d); err != nil {
		return nil, err
	}
	return &models.UserData{UserID: s.cfg.Usr, ID: id, Type: t, Data: d}, nil
}

const queryEnterUserData = `
	insert into data (usr, id, type, data) values ($1, $2, $3, $4)
`

// EnterUserData enters user data.
func (s *db) EnterUserData(ctx context.Context, d *models.UserData) error {
	if _, err := s.conn.ExecContext(ctx, queryEnterUserData,
		d.UserID, d.ID, d.Type, d.Data); err != nil {
		return err
	}
	return nil
}
