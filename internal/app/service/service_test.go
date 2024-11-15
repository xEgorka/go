package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/xEgorka/project3/internal/app/auth"
	"github.com/xEgorka/project3/internal/app/config"
	"github.com/xEgorka/project3/internal/app/jobs"
	"github.com/xEgorka/project3/internal/app/mocks"
	"github.com/xEgorka/project3/internal/app/models"
	"github.com/xEgorka/project3/internal/app/storage"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	type args struct {
		config *config.Config
		store  storage.Storage
		jobs   *jobs.Jobs
		auth   auth.Auth
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(&config.Config{}, ms, j, ma)
			if reflect.TypeOf(got) == reflect.TypeOf((*Service)(nil)).Elem() {
				t.Errorf("not service")
			}
		})
	}
}

func TestService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := New(&config.Config{}, ms, j, ma)
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		j   *jobs.Jobs
		a   auth.Auth
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
			name: "positive test #1",
		},
		{
			name:    "negative test #1",
			wantErr: true,
		},
		{
			name:    "negative test #2",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "positive test #1" {
				ma.EXPECT().Register(tt.args.ctx, tt.args.usr, tt.args.pass).
					Return(nil)
			}
			if tt.name == "negative test #1" {
				ma.EXPECT().Register(tt.args.ctx, tt.args.usr, tt.args.pass).
					Return(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
			}
			if tt.name == "negative test #2" {
				ma.EXPECT().Register(tt.args.ctx, tt.args.usr, tt.args.pass).
					Return(errors.New("test"))
			}
			if err := s.Register(tt.args.ctx, tt.args.usr, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("Service.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := New(&config.Config{}, ms, j, ma)
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		j   *jobs.Jobs
		a   auth.Auth
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
			name: "positive test #1",
		},
		{
			name:    "negative test #1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "positive test #1" {
				ma.EXPECT().Login(tt.args.ctx, tt.args.usr, tt.args.pass).
					Return(nil)
			}
			if tt.name == "negative test #1" {
				ma.EXPECT().Login(tt.args.ctx, tt.args.usr, tt.args.pass).
					Return(errors.New("test"))
			}
			if err := s.Login(tt.args.ctx, tt.args.usr, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := New(&config.Config{}, ms, j, ma)
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		j   *jobs.Jobs
		a   auth.Auth
	}
	type args struct {
		usr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetToken(tt.args.usr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != 121 {
				t.Errorf("Service.GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := New(&config.Config{}, ms, j, ma)
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		j   *jobs.Jobs
		a   auth.Auth
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
		{
			name:    "negative test #1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.name == "positive test #1" {
			ms.EXPECT().Ping().Return(nil)

		}
		if tt.name == "negative test #1" {
			ms.EXPECT().Ping().Return(errors.New("test"))
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Service.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := New(&config.Config{}, ms, j, ma)
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		j   *jobs.Jobs
		a   auth.Auth
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
			name:    "positive test #1",
			wantErr: true,
		},
		{
			name:    "positive test #2",
			want:    []models.UserData{{ID: "testID"}},
			wantErr: false,
		},
		{
			name:    "negative test #1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "positive test #1" {
				ms.EXPECT().GetUserData(tt.args.ctx, tt.args.usr, tt.args.timestamp).Return(nil, nil)
			}
			if tt.name == "positive test #2" {
				ms.EXPECT().GetUserData(tt.args.ctx, tt.args.usr, tt.args.timestamp).Return([]models.UserData{{ID: "testID"}}, nil)
			}
			if tt.name == "negative test #1" {
				ms.EXPECT().GetUserData(tt.args.ctx, tt.args.usr, tt.args.timestamp).Return(nil, errors.New("test"))
			}
			got, err := s.GetUserData(tt.args.ctx, tt.args.usr, tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetUserData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_MergeUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := New(&config.Config{}, ms, j, ma)
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		j   *jobs.Jobs
		a   auth.Auth
	}
	type args struct {
		r []models.UserData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.MergeUserData(tt.args.r)
		})
	}
}
