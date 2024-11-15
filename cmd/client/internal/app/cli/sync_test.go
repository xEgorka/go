package cli

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"text/tabwriter"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/crypto"
	"github.com/xEgorka/project3/cmd/client/internal/app/mocks"
	"github.com/xEgorka/project3/cmd/client/internal/app/models"
	"github.com/xEgorka/project3/cmd/client/internal/app/requests"
	"github.com/xEgorka/project3/cmd/client/internal/app/service"
)

func TestApp_download(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)

	type fields struct {
		cfg *config.Config
		s   *service.Service
		c   crypto.Crypto
		r   *bufio.Scanner
		w   *bufio.Writer
		tw  *tabwriter.Writer
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive test #1",
			args: args{ctx: context.Background()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					c := &http.Cookie{
						Name:  "token",
						Value: "testToken",
					}
					http.SetCookie(w, c)
					w.Header().Set("Content-type", "text/plain")
					w.WriteHeader(http.StatusOK)
				}))
			defer func() { srv.Close() }()
			cfg := &config.Config{AppAddr: srv.URL + "/", Usr: "x0o1@ya.ru"}
			r := requests.New(cfg, "")
			s := service.New(cfg, ms, ma, r)
			a := &App{
				cfg: cfg,
				s:   s,
				c:   crypto.New(),
				r:   bufio.NewScanner(os.Stdin),
				w:   bufio.NewWriter(os.Stdout),
				tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
			}
			ms.EXPECT().GetTimestamp(tt.args.ctx).Return("0", nil)
			a.download(tt.args.ctx)
		})
	}
}

func TestApp_download2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)

	type fields struct {
		cfg *config.Config
		s   *service.Service
		c   crypto.Crypto
		r   *bufio.Scanner
		w   *bufio.Writer
		tw  *tabwriter.Writer
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive test #1",
			args: args{ctx: context.Background()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if os.Getenv("BE_CRASHER") == "1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.Header().Set("Content-type", "text/plain")
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { srv.Close() }()
				cfg := &config.Config{AppAddr: srv.URL + "/", Usr: "x0o1@ya.ru"}
				r := requests.New(cfg, "")
				s := service.New(cfg, ms, ma, r)
				a := &App{
					cfg: cfg,
					s:   s,
					c:   crypto.New(),
					r:   bufio.NewScanner(os.Stdin),
					w:   bufio.NewWriter(os.Stdout),
					tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
				}
				ms.EXPECT().GetTimestamp(tt.args.ctx).
					Return("", status.Error(codes.Unauthenticated, "unauthenticated"))
				a.download(tt.args.ctx)
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_download2")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)

	type fields struct {
		cfg *config.Config
		s   *service.Service
		c   crypto.Crypto
		r   *bufio.Scanner
		w   *bufio.Writer
		tw  *tabwriter.Writer
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive test #1",
			args: args{ctx: context.Background()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					c := &http.Cookie{
						Name:  "token",
						Value: "testToken",
					}
					http.SetCookie(w, c)
					w.Header().Set("Content-type", "text/plain")
					w.WriteHeader(http.StatusOK)
				}))
			defer func() { srv.Close() }()
			cfg := &config.Config{AppAddr: srv.URL + "/", Usr: "x0o1@ya.ru"}
			r := requests.New(cfg, "")
			s := service.New(cfg, ms, ma, r)
			a := &App{
				cfg: cfg,
				s:   s,
				c:   crypto.New(),
				r:   bufio.NewScanner(os.Stdin),
				w:   bufio.NewWriter(os.Stdout),
				tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
			}
			ms.EXPECT().GetUnmergedUserData(tt.args.ctx, cfg.Usr).
				Return([]models.UserData{{ID: "x0o1@ya.ru"}}, nil)
			a.upload(tt.args.ctx)
		})
	}
}

func TestApp_upload2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)

	type fields struct {
		cfg *config.Config
		s   *service.Service
		c   crypto.Crypto
		r   *bufio.Scanner
		w   *bufio.Writer
		tw  *tabwriter.Writer
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive test #1",
			args: args{ctx: context.Background()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if os.Getenv("BE_CRASHER") == "1" {
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.Header().Set("Content-type", "text/plain")
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { srv.Close() }()
				cfg := &config.Config{AppAddr: srv.URL + "/", Usr: "x0o1@ya.ru"}
				r := requests.New(cfg, "")
				s := service.New(cfg, ms, ma, r)
				a := &App{
					cfg: cfg,
					s:   s,
					c:   crypto.New(),
					r:   bufio.NewScanner(os.Stdin),
					w:   bufio.NewWriter(os.Stdout),
					tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
				}
				ms.EXPECT().GetUnmergedUserData(tt.args.ctx, cfg.Usr).
					Return(nil, status.Error(codes.Unauthenticated, "unauthenticated"))
				a.upload(tt.args.ctx)
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_upload2")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}
