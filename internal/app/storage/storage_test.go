package storage

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
	"github.com/xEgorka/project3/internal/app/models"
)

func TestOpen(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    Storage
		wantErr bool
	}{
		{
			name: "negative test #1",
			args: args{
				ctx: context.Background(),
				cfg: &config.Config{DBDriver: "bad"}},
			wantErr: true,
		},
		{
			name: "negative test #2",
			args: args{
				ctx: context.Background(),
				cfg: &config.Config{DBDriver: "pgx"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Open(tt.args.ctx, tt.args.cfg)
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
		conn *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    Storage
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				ctx:  context.Background(),
				cfg:  &config.Config{},
				conn: conn},
			wantErr: false,
		},
		{
			name: "negative test #1",
			args: args{
				ctx: context.Background(),
				cfg: &config.Config{}, conn: conn},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "positive test #1" {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap1)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap2)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			}
			if tt.name == "negative test #1" {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap1)).
					WillReturnError(errors.New("test"))
			}
			_, err := open(tt.args.ctx, tt.args.cfg, tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_db_bootstrap(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}
	defer conn.Close()
	type fields struct {
		conn *sql.DB
		cfg  *config.Config
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "negative test #1",
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
		{
			name:    "negative test #2",
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
		{
			name:    "negative test #3",
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
		{
			name:    "negative test #4",
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
		{
			name:    "negative test #5",
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
		{
			name:    "positive test #1",
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &db{conn: conn}
			if tt.name == "negative test #1" {
				mock.ExpectBegin().WillReturnError(errors.New("test"))
			}
			if tt.name == "negative test #2" {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap1)).
					WillReturnError(errors.New("test"))
				mock.ExpectRollback()
			}
			if tt.name == "negative test #3" {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap1)).
					WillReturnError(errors.New("test"))
				mock.ExpectRollback().WillReturnError(errors.New("test"))
			}
			if tt.name == "negative test #4" {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap1)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap2)).
					WillReturnError(errors.New("test"))
				mock.ExpectRollback()
			}
			if tt.name == "negative test #5" {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap1)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap2)).
					WillReturnError(errors.New("test"))
				mock.ExpectRollback().WillReturnError(errors.New("test"))
			}
			if tt.name == "positive test #1" {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap1)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(queryBootstrap2)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			}
			if err := s.bootstrap(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("db.bootstrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_db_GetUserData(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}
	defer conn.Close()

	type fields struct {
		conn *sql.DB
		cfg  *config.Config
	}
	type args struct {
		ctx       context.Context
		usr       string
		timestamp int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.UserData
		wantErr bool
	}{
		{
			name:    "negative test #1",
			args:    args{ctx: context.Background(), usr: "testUsr", timestamp: 0},
			fields:  fields{conn: conn, cfg: &config.Config{}},
			wantErr: true,
		},
		{
			name:    "positive test #1",
			args:    args{ctx: context.Background(), usr: "testUsr", timestamp: 0},
			fields:  fields{conn: conn, cfg: &config.Config{}},
			wantErr: false,
		},
		{
			name:    "negative test #2",
			args:    args{ctx: context.Background(), usr: "testUsr", timestamp: 0},
			fields:  fields{conn: conn, cfg: &config.Config{}},
			wantErr: true,
		},
		{
			name:    "negative test #3",
			args:    args{ctx: context.Background(), usr: "testUsr", timestamp: 0},
			fields:  fields{conn: conn, cfg: &config.Config{}},
			wantErr: true,
		},
		{
			name:    "negative test #4",
			args:    args{ctx: context.Background(), usr: "testUsr", timestamp: 0},
			fields:  fields{conn: conn, cfg: &config.Config{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := db{
				conn: tt.fields.conn,
				cfg:  tt.fields.cfg,
			}
			if tt.name == "negative test #1" {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUserData)).
					WillReturnError(errors.New("test"))
			}
			if tt.name == "positive test #1" {
				mockRows := sqlmock.NewRows(
					[]string{"usr", "id", "type", "data", "updatedStr", "mergedStr"}).
					AddRow("x0o1@ya.ru", "passID", 1, "encrypted", "2024-06-15T07:09:43Z", "2024-06-15T07:09:44Z")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUserData)).
					WithArgs(tt.args.usr, tt.args.timestamp).
					WillReturnRows(mockRows)
			}
			if tt.name == "negative test #2" {
				mockRows := sqlmock.NewRows(
					[]string{"usr", "id", "type", "data", "updatedStr", "mergedStr"}).
					AddRow("x0o1@ya.ru", "passID", 1, "encrypted", "2024-06-15T07:09:43Z", "2024-06-15T07:09:44Z")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUserData)).
					WithArgs(tt.args.usr, tt.args.timestamp).
					WillReturnRows(mockRows)
				mockRows.CloseError(errors.New("test"))
			}
			if tt.name == "negative test #3" {
				mockRows := sqlmock.NewRows(
					[]string{"usr", "id", "type", "data", "updatedStr", "mergedStr"}).
					AddRow("x0o1@ya.ru", "passID", 1, "encrypted", "bad2024-06-15T07:09:43Z", "2024-06-15T07:09:44Z")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUserData)).
					WithArgs(tt.args.usr, tt.args.timestamp).
					WillReturnRows(mockRows)
			}

			if tt.name == "negative test #4" {
				mockRows := sqlmock.NewRows(
					[]string{"usr", "id", "type", "data", "updatedStr", "mergedStr"}).
					AddRow("x0o1@ya.ru", "passID", 1, "encrypted", "2024-06-15T07:09:43Z", "bad2024-06-15T07:09:44Z")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUserData)).
					WithArgs(tt.args.usr, tt.args.timestamp).
					WillReturnRows(mockRows)
			}

			_, err := s.GetUserData(tt.args.ctx, tt.args.usr, tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.GetUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_db_MergeUserData(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}
	defer conn.Close()

	type fields struct {
		conn *sql.DB
		cfg  *config.Config
	}
	type args struct {
		ctx context.Context
		d   []models.UserData
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
				ctx: context.Background(),
				d:   []models.UserData{{ID: "testID"}}},
			wantErr: true,
		},
		{
			name: "positive test #1",
			args: args{
				ctx: context.Background(),
				d:   []models.UserData{{ID: "testID"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &db{conn: conn}
			if tt.name == "negative test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryMergeUserData)).
					WillReturnError(errors.New("test"))
			}
			if tt.name == "positive test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryMergeUserData)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			}

			if err := s.MergeUserData(tt.args.ctx, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("db.MergeUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_db_Ping(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}
	defer conn.Close()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &db{conn: conn}
			mock.ExpectPing()
			if err := s.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("db.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_db_Close(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}
	defer conn.Close()
	type fields struct {
		conn *sql.DB
		cfg  *config.Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "positive test #1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectClose()
			s := &db{conn: conn}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("db.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
