package storage

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"

	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/crypto"
	"github.com/xEgorka/project3/cmd/client/internal/app/models"
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
		want    Storage
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
				cfg: &config.Config{DBDriver: "sqlite3"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Open(tt.args.ctx, tt.args.cfg)
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
		want    Storage
		wantErr bool
	}{
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
			if tt.name == "negative test #1" {
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

func Test_db_GetTimestamp(t *testing.T) {
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
		want    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &db{conn: conn, cfg: &config.Config{Usr: "testUsr"}}
			if tt.name == "positive test #1" {
				mockRows := sqlmock.NewRows([]string{"timestamp"}).AddRow("0")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetTimestamp)).
					WithArgs(s.cfg.Usr).
					WillReturnRows(mockRows)
			}
			_, err := s.GetTimestamp(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.GetTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_db_UpdateUserData(t *testing.T) {
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
			name: "positive test #1",
			args: args{
				ctx: context.Background(),
				d:   []models.UserData{{ID: "testID"}}},
			wantErr: false,
		},
		{
			name: "negative test #1",
			args: args{
				ctx: context.Background(),
				d:   []models.UserData{{ID: "testID"}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &db{conn: conn, cfg: &config.Config{Usr: "testUsr"}}
			if tt.name == "positive test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryUpdateUserData)).
					WithArgs(tt.args.d[0].UserID, tt.args.d[0].ID,
						tt.args.d[0].Type, tt.args.d[0].Data,
						tt.args.d[0].Updated.Unix(), tt.args.d[0].Merged.Unix()).
					WillReturnResult(sqlmock.NewResult(0, 0))
			}
			if tt.name == "negative test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryUpdateUserData)).
					WithArgs(tt.args.d[0].UserID, tt.args.d[0].ID,
						tt.args.d[0].Type, tt.args.d[0].Data,
						tt.args.d[0].Updated.Unix(), tt.args.d[0].Merged.Unix()).
					WillReturnError(errors.New("test"))
			}
			if err := s.UpdateUserData(tt.args.ctx, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("db.UpdateUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_db_GetUnmergedUserData(t *testing.T) {
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
		usr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.UserData
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				ctx: context.Background(),
				usr: "testUsr",
			},
			wantErr: false,
		},
		{
			name: "negative test #1",
			args: args{
				ctx: context.Background(),
				usr: "testUsr",
			},
			wantErr: true,
		},
		{
			name: "negative test #2",
			args: args{
				ctx: context.Background(),
				usr: "testUsr",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &db{conn: conn}
			if tt.name == "positive test #1" {
				mockRows := sqlmock.
					NewRows([]string{"usr", "id", "type", "data", "updated"}).
					AddRow("testUsr", "testID", "1", "encrypted", "2024-06-15T07:09:43Z")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUnmergedUserData)).
					WithArgs(tt.args.usr).
					WillReturnRows(mockRows)
			}
			if tt.name == "negative test #1" {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUnmergedUserData)).
					WithArgs(tt.args.usr).
					WillReturnError(errors.New("test"))
			}
			if tt.name == "negative test #2" {
				mockRows := sqlmock.
					NewRows([]string{"usr", "id", "type", "data", "updated"}).
					AddRow("testUsr", "testID", "1", "encrypted", "bad2024-06-15T07:09:43Z")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUnmergedUserData)).
					WithArgs(tt.args.usr).
					WillReturnRows(mockRows)
			}
			_, err := s.GetUnmergedUserData(tt.args.ctx, tt.args.usr)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.GetUnmergedUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_db_IndexUserData(t *testing.T) {
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
		want    []models.UserData
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
		{
			name:    "negative test #1",
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
		{
			name:    "negative test #2",
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &db{conn: conn, cfg: &config.Config{Usr: "testUsr"}}
			if tt.name == "positive test #1" {
				mockRows := sqlmock.
					NewRows([]string{"id", "type", "updated"}).
					AddRow("testID", "PASS", "2024-06-15T07:09:43Z")
				mock.ExpectQuery(regexp.QuoteMeta(queryIndexUserData)).
					WithArgs(s.cfg.Usr).
					WillReturnRows(mockRows)
			}
			if tt.name == "negative test #1" {
				mock.ExpectQuery(regexp.QuoteMeta(queryIndexUserData)).
					WithArgs(s.cfg.Usr).
					WillReturnError(errors.New("test"))
			}
			if tt.name == "negative test #2" {
				mockRows := sqlmock.
					NewRows([]string{"id", "type", "updated"}).
					AddRow("testID", "PASS", "bad2024-06-15T07:09:43Z")
				mock.ExpectQuery(regexp.QuoteMeta(queryIndexUserData)).
					WithArgs(s.cfg.Usr).
					WillReturnRows(mockRows)
			}
			_, err := s.IndexUserData(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.IndexUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_db_PrintUserData(t *testing.T) {
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
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.UserData
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{ctx: context.Background(), id: "testID"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &db{conn: conn, cfg: &config.Config{Usr: "testUsr"}}
			if tt.name == "positive test #1" {
				mockRows := sqlmock.
					NewRows([]string{"type", "data"}).
					AddRow("1", "encrypted")
				mock.ExpectQuery(regexp.QuoteMeta(queryPrintUserData)).
					WithArgs(s.cfg.Usr, tt.args.id).
					WillReturnRows(mockRows)
			}
			_, err := s.PrintUserData(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.PrintUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_db_EnterUserData(t *testing.T) {
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
		d   *models.UserData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				ctx: context.Background(),
				d: &models.UserData{
					UserID: "testUsr",
					ID:     "testID"}},
			wantErr: false,
		},
		{
			name: "negative test #1",
			args: args{
				ctx: context.Background(),
				d: &models.UserData{
					UserID: "testUsr",
					ID:     "testID"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &db{conn: conn, cfg: &config.Config{Usr: "testUsr"}}
			if tt.name == "positive test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryEnterUserData)).
					WithArgs(tt.args.d.UserID, tt.args.d.ID, tt.args.d.Type, tt.args.d.Data).
					WillReturnResult(sqlmock.NewResult(0, 0))
			}
			if tt.name == "negative test #1" {
				mock.ExpectExec(regexp.QuoteMeta(queryEnterUserData)).
					WithArgs(tt.args.d.UserID, tt.args.d.ID, tt.args.d.Type, tt.args.d.Data).
					WillReturnError(errors.New("test"))
			}
			if err := s.EnterUserData(tt.args.ctx, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("db.EnterUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
