package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xEgorka/project3/cmd/client/internal/app/auth"
	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/logger"
	"github.com/xEgorka/project3/cmd/client/internal/app/mocks"
	"github.com/xEgorka/project3/cmd/client/internal/app/models"
	"github.com/xEgorka/project3/cmd/client/internal/app/requests"
	"github.com/xEgorka/project3/cmd/client/internal/app/storage"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)
	cfg := &config.Config{}
	r := requests.New(cfg, "")
	type args struct {
		config *config.Config
		store  storage.Storage
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
			got := New(cfg, ms, ma, r)
			if reflect.TypeOf(got) == reflect.TypeOf((*Service)(nil)).Elem() {
				t.Errorf("not service")
			}
		})
	}
}

func TestService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)
	cfg := &config.Config{}
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		a   auth.Auth
		r   *requests.HTTP
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
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			wantErr: false,
		},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *Service
			if tt.name == "positive test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ma.EXPECT().LoginOnline(tt.args.ctx, tt.args.usr, tt.args.pass).
					Return(nil)
			}
			if tt.name == "negative test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.WriteHeader(http.StatusBadRequest)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
			}
			if tt.name == "negative test #2" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ma.EXPECT().LoginOnline(tt.args.ctx, tt.args.usr, tt.args.pass).
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
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)
	cfg := &config.Config{}
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		a   auth.Auth
		r   *requests.HTTP
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
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			wantErr: false,
		},
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
			name: "negative test #3",
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			wantErr: true,
		},
		{
			name: "negative test #4",
			args: args{
				ctx:  context.Background(),
				usr:  "testUsr",
				pass: "testPass"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *Service
			if tt.name == "positive test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ma.EXPECT().LoginOnline(tt.args.ctx, tt.args.usr, tt.args.pass).
					Return(nil)
			}
			if tt.name == "negative test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusBadRequest)
						}
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
			}
			if tt.name == "negative test #2" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.WriteHeader(http.StatusBadRequest)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ma.EXPECT().LoginOffline(tt.args.ctx, tt.args.usr, tt.args.pass).
					Return(nil)
			}
			if tt.name == "negative test #3" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.WriteHeader(http.StatusBadRequest)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ma.EXPECT().LoginOffline(tt.args.ctx, tt.args.usr, tt.args.pass).
					Return(status.Error(codes.Unauthenticated, "unauthenticated"))
			}
			if tt.name == "negative test #4" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ma.EXPECT().LoginOnline(tt.args.ctx, tt.args.usr, tt.args.pass).
					Return(status.Error(codes.Unauthenticated, "unauthenticated"))
			}
			if err := s.Login(tt.args.ctx, tt.args.usr, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_DownloadUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)
	cfg := &config.Config{}
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		a   auth.Auth
		r   *requests.HTTP
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
		{
			name:    "negative test #3",
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *Service
			if tt.name == "positive test #1" {
				res := []models.UserData{{UserID: "testUsr"}}
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if err := json.NewEncoder(w).Encode(&res); err != nil {
							logger.Log.Info("JSON encode error", zap.Error(err))
							http.Error(w, "Internal server error", http.StatusInternalServerError)
							return
						}
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().GetTimestamp(tt.args.ctx).
					Return("0", nil)
				ms.EXPECT().UpdateUserData(tt.args.ctx, res).
					Return(nil)
			}
			if tt.name == "negative test #1" {
				res := []models.UserData{{UserID: "testUsr"}}
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if err := json.NewEncoder(w).Encode(&res); err != nil {
							logger.Log.Info("JSON encode error", zap.Error(err))
							http.Error(w, "Internal server error", http.StatusInternalServerError)
							return
						}
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().GetTimestamp(tt.args.ctx).
					Return("", errors.New("test"))
			}
			if tt.name == "negative test #2" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusBadRequest)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().GetTimestamp(tt.args.ctx).
					Return("0", nil)
			}
			if tt.name == "negative test #3" {
				res := []models.UserData{{UserID: "testUsr"}}
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if err := json.NewEncoder(w).Encode(&res); err != nil {
							logger.Log.Info("JSON encode error", zap.Error(err))
							http.Error(w, "Internal server error", http.StatusInternalServerError)
							return
						}
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().GetTimestamp(tt.args.ctx).
					Return("0", nil)
				ms.EXPECT().UpdateUserData(tt.args.ctx, res).
					Return(errors.New("test"))
			}
			if err := s.DownloadUserData(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Service.DownloadUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_UploadUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)
	cfg := &config.Config{Usr: "testUsr"}
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		a   auth.Auth
		r   *requests.HTTP
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
			var s *Service
			d := []models.UserData{{UserID: "testUsr"}}
			if tt.name == "positive test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusAccepted)
						}
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().GetUnmergedUserData(tt.args.ctx, s.cfg.Usr).
					Return(d, nil)
			}
			if tt.name == "negative test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusBadRequest)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().GetUnmergedUserData(tt.args.ctx, s.cfg.Usr).
					Return(nil, errors.New("test"))
			}
			if tt.name == "negative test #2" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusBadRequest)
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().GetUnmergedUserData(tt.args.ctx, s.cfg.Usr).
					Return(d, nil)
			}
			if err := s.UploadUserData(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Service.UploadUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_IndexUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)
	cfg := &config.Config{Usr: "testUsr"}
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		a   auth.Auth
		r   *requests.HTTP
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
			var s *Service
			d := []models.UserData{{UserID: "testUsr"}}
			if tt.name == "positive test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusAccepted)
						}
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().IndexUserData(tt.args.ctx).
					Return(d, nil)
			}
			if tt.name == "negative test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusAccepted)
						}
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().IndexUserData(tt.args.ctx).
					Return(nil, errors.New("test"))
			}
			if tt.name == "negative test #2" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusAccepted)
						}
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().IndexUserData(tt.args.ctx).
					Return([]models.UserData{}, nil)
			}
			_, err := s.IndexUserData(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.IndexUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_PrintUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)
	cfg := &config.Config{Usr: "testUsr"}
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		a   auth.Auth
		r   *requests.HTTP
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
		{
			name:    "negative test #1",
			args:    args{ctx: context.Background(), id: "testID"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *Service
			d := &models.UserData{UserID: "testUsr"}
			if tt.name == "positive test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusAccepted)
						}
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().PrintUserData(tt.args.ctx, tt.args.id).
					Return(d, nil)
			}
			if tt.name == "negative test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusAccepted)
						}
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				ms.EXPECT().PrintUserData(tt.args.ctx, tt.args.id).
					Return(nil, errors.New("test"))
			}
			_, err := s.PrintUserData(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.PrintUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_EnterUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)
	cfg := &config.Config{Usr: "testUsr"}
	type fields struct {
		cfg *config.Config
		s   storage.Storage
		a   auth.Auth
		r   *requests.HTTP
	}
	type args struct {
		ctx context.Context
		id  string
		t   int
		c   string
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
				id:  "testID",
				t:   1,
				c:   "encrypted",
			},
			wantErr: false,
		},
		{
			name: "negative test #1",
			args: args{
				ctx: context.Background(),
				id:  "testID",
				t:   1,
				c:   "encrypted",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *Service
			if tt.name == "positive test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusAccepted)
						}
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				d := &models.UserData{
					UserID: s.cfg.Usr,
					ID:     tt.args.id,
					Type:   tt.args.t,
					Data:   tt.args.c,
				}
				ms.EXPECT().EnterUserData(tt.args.ctx, d).
					Return(nil)
			}
			if tt.name == "negative test #1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusAccepted)
						}
					}))
				defer func() { srv.Close() }()
				r := requests.New(&config.Config{AppAddr: srv.URL + "/"}, "")
				s = New(cfg, ms, ma, r)
				d := &models.UserData{
					UserID: s.cfg.Usr,
					ID:     tt.args.id,
					Type:   tt.args.t,
					Data:   tt.args.c,
				}
				ms.EXPECT().EnterUserData(tt.args.ctx, d).
					Return(errors.New("test"))
			}
			if err := s.EnterUserData(tt.args.ctx, tt.args.id, tt.args.t, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Service.EnterUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
