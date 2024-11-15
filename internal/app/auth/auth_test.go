package auth

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/xEgorka/project3/internal/app/config"
	"github.com/xEgorka/project3/internal/app/crypto"
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
				cfg: &config.Config{DBDriver: "pgx"},
				c:   crypto.New()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Open(tt.args.ctx, tt.args.cfg, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() = %v, want %v", got, tt.want)
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
			args: args{
				ctx:  context.Background(),
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
	conn, _ := sql.Open("pgx", "")
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

func Test_db_Register(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}
	defer conn.Close()
	type fields struct {
		conn *sql.DB
		cfg  *config.Config
		c    crypto.Crypto
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
			wantErr: true,
		},
		{
			name: "negative test #2",
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			wantErr: true,
		},
		{
			name: "positive test #1",
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &db{
				conn: conn,
				cfg:  tt.fields.cfg,
				c:    crypto.New(),
			}
			if tt.name == "negative test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryRegister)).
					WillReturnError(errors.New("test"))
			}
			if tt.name == "positive test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryRegister)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			}
			if err := a.Register(tt.args.ctx, tt.args.usr, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("db.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_db_Login(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}
	defer conn.Close()

	type fields struct {
		conn *sql.DB
		cfg  *config.Config
		c    crypto.Crypto
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
			wantErr: true,
		},
		{
			name: "negative test #2",
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			wantErr: true,
		},
		{
			name: "positive test #1",
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &db{
				conn: conn,
				cfg:  tt.fields.cfg,
				c:    crypto.New(),
			}
			if tt.name == "negative test #1" {
				mock.ExpectQuery(regexp.QuoteMeta(queryLogin)).
					WillReturnError(errors.New("test"))
			}
			if tt.name == "negative test #2" {
				mockRows := sqlmock.NewRows([]string{"hash"}).AddRow("bad")
				mock.ExpectQuery(regexp.QuoteMeta(queryLogin)).
					WillReturnRows(mockRows)
			}
			if tt.name == "positive test #1" {
				mockRows := sqlmock.NewRows([]string{"hash"}).
					AddRow("$2a$10$OA1MRqItnPyDQXzm3HY29uWJ7W7Q1kHX7fRnLqqxBiIDi1/fd0zgS")
				mock.ExpectQuery(regexp.QuoteMeta(queryLogin)).
					WillReturnRows(mockRows)
			}
			if err := a.Login(tt.args.ctx, tt.args.usr, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("db.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
