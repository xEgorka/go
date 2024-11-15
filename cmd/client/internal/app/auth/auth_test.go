// Package auth provides interface to authentication database.
package auth

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"

	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/crypto"
)

func TestOpen(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *config.Config
		c   crypto.Crypto
	}
	tests := []struct {
		name    string
		args    args
		want    Auth
		wantErr bool
	}{
		{
			name: "negative test #1",
			args: args{
				ctx: context.Background(),
				cfg: &config.Config{},
				c:   crypto.New()},
			wantErr: true,
		},
		{
			name: "negative test #2",
			args: args{
				ctx: context.Background(),
				cfg: &config.Config{DBDriver: "sqlite3"},
				c:   crypto.New()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Open(tt.args.ctx, tt.args.cfg, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_open(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}
	defer conn.Close()
	type args struct {
		ctx  context.Context
		cfg  *config.Config
		c    crypto.Crypto
		conn *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    Auth
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				ctx:  context.Background(),
				cfg:  &config.Config{},
				c:    crypto.New(),
				conn: conn},
			wantErr: false,
		},
		{
			name: "negative test #1",
			args: args{ctx: context.Background(),
				cfg:  &config.Config{},
				c:    crypto.New(),
				conn: conn},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "positive test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			}
			if tt.name == "negative test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap)).
					WillReturnError(errors.New("test"))
			}
			_, err := open(tt.args.ctx, tt.args.cfg, tt.args.c, tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_new(t *testing.T) {
	conn, _ := sql.Open("sqlite3", "")
	type args struct {
		config *config.Config
		conn   *sql.DB
		crypto crypto.Crypto
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{config: &config.Config{}, conn: conn, crypto: crypto.New()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(tt.args.config, tt.args.conn, tt.args.crypto)
			if reflect.TypeOf(got) == reflect.TypeOf((*db)(nil)).Elem() {
				t.Errorf("not db")
			}
		})
	}
}

func Test_db_LoginOnline(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}
	defer conn.Close()
	type fields struct {
		c    crypto.Crypto
		conn *sql.DB
		cfg  *config.Config
	}
	type args struct {
		ctx  context.Context
		usr  string
		pass string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "negative test #1",
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			fields:  fields{c: crypto.New(), conn: conn, cfg: &config.Config{}},
			wantErr: true,
		},
		{
			name: "positive test #1",
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			fields:  fields{c: crypto.New(), conn: conn, cfg: &config.Config{}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &db{
				c:    tt.fields.c,
				conn: tt.fields.conn,
				cfg:  tt.fields.cfg,
			}
			if tt.name == "negative test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryLoginOnline)).
					WillReturnError(errors.New("test"))
			}
			if tt.name == "positive test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryLoginOnline)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			}
			if err := a.LoginOnline(tt.args.ctx, tt.args.usr, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("db.LoginOnline() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_db_LoginOffline(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}
	defer conn.Close()
	type fields struct {
		c    crypto.Crypto
		conn *sql.DB
		cfg  *config.Config
	}
	type args struct {
		ctx  context.Context
		usr  string
		pass string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "negative test #1",
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			fields:  fields{c: crypto.New(), conn: conn, cfg: &config.Config{}},
			wantErr: true,
		},
		{
			name: "positive test #1",
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			fields:  fields{c: crypto.New(), conn: conn, cfg: &config.Config{}},
			wantErr: false,
		},
		{
			name: "negative test #2",
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			fields:  fields{c: crypto.New(), conn: conn, cfg: &config.Config{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &db{
				c:    tt.fields.c,
				conn: tt.fields.conn,
				cfg:  tt.fields.cfg,
			}
			if tt.name == "negative test #1" {
				mock.ExpectQuery(regexp.QuoteMeta(queryLoginOffline)).
					WillReturnError(errors.New("test"))
			}
			if tt.name == "positive test #1" {
				mockRows := sqlmock.NewRows([]string{"hash"}).
					AddRow("7b982b524f859f466cbe42f3ab0f29bc377fd44be47e3d1ea1239249e7af9ac983432e")
				mock.ExpectQuery(regexp.QuoteMeta(queryLoginOffline)).
					WillReturnRows(mockRows)
			}
			if tt.name == "negative test #2" {
				mockRows := sqlmock.NewRows([]string{"hash"}).
					AddRow("bad7b982b524f859f466cbe42f3ab0f29bc377fd44be47e3d1ea1239249e7af9ac983432e")
				mock.ExpectQuery(regexp.QuoteMeta(queryLoginOffline)).
					WillReturnRows(mockRows)
			}

			if err := a.LoginOffline(tt.args.ctx, tt.args.usr, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("db.LoginOffline() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
