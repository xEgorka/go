package requests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/models"
)

func TestHTTP_ping(t *testing.T) {
	type fields struct {
		cfg *config.Config
		c   *http.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		srv     *httptest.Server
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "positive test #1" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { s.Close() }()
				h := New(&config.Config{AppAddr: s.URL + "/"}, "")
				if err := h.ping(tt.args.ctx); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.ping() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #1" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusBadRequest)
					}))
				defer func() { s.Close() }()
				h := New(&config.Config{AppAddr: s.URL + "/"}, "")
				if err := h.ping(tt.args.ctx); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.ping() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestHTTP_Register(t *testing.T) {
	type fields struct {
		cfg *config.Config
		c   *http.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		srv     *httptest.Server
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
		{
			name:    "negative test #4",
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "positive test #1" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.Header().Set("Content-type", "text/plain")
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				r := New(cfg, "")
				if err := r.Register(context.Background(), "x0o1@ya.ru", "password"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.Register() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #1" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.Header().Set("Content-type", "text/plain")
						w.WriteHeader(http.StatusBadRequest)
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.Register(context.Background(), "x0o1@ya.ru", "password"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.Register() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #2" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.Register(context.Background(), "x0o1@ya.ru", "password"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.Register() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #3" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusConflict)
						}
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.Register(context.Background(), "x0o1@ya.ru", "password"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.Register() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #4" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusBadRequest)
						}
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.Register(context.Background(), "x0o1@ya.ru", "password"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.Register() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestHTTP_Login(t *testing.T) {
	type fields struct {
		cfg *config.Config
		c   *http.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		srv     *httptest.Server
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
			if tt.name == "positive test #1" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.Header().Set("Content-type", "text/plain")
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.Login(context.Background(), "x0o1@ya.ru", "password"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.Login() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #1" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						c := &http.Cookie{
							Name:  "token",
							Value: "testToken",
						}
						http.SetCookie(w, c)
						w.Header().Set("Content-type", "text/plain")
						w.WriteHeader(http.StatusBadRequest)
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.Login(context.Background(), "x0o1@ya.ru", "password"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.Login() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #2" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						w.WriteHeader(http.StatusOK)
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.Login(context.Background(), "x0o1@ya.ru", "password"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.Login() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #3" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusConflict)
						}
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.Login(context.Background(), "x0o1@ya.ru", "password"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.Login() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestHTTP_DownloadUserData(t *testing.T) {
	type fields struct {
		cfg *config.Config
		c   *http.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		srv     *httptest.Server
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
			wantErr: true,
		},
		{
			name:    "negative test #5",
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "positive test #1" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						_, err := w.Write([]byte(`[{"ID": "secret pass"}]`))
						if err != nil {
							panic(err)
						}
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if _, err := h.DownloadUserData(context.Background(), "0"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.DownloadUserData() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #1" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						w.WriteHeader(http.StatusBadRequest)
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if _, err := h.DownloadUserData(context.Background(), "0"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.DownloadUserData() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #2" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusNoContent)
						}
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if _, err := h.DownloadUserData(context.Background(), "0"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.DownloadUserData() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #3" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusUnauthorized)
						}
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if _, err := h.DownloadUserData(context.Background(), "0"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.DownloadUserData() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #4" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusBadRequest)
						}
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if _, err := h.DownloadUserData(context.Background(), "0"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.DownloadUserData() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #5" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						_, err := w.Write([]byte(`[bad json]`))
						if err != nil {
							panic(err)
						}
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if _, err := h.DownloadUserData(context.Background(), "0"); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.DownloadUserData() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestHTTP_UploadUserData(t *testing.T) {
	type fields struct {
		cfg *config.Config
		c   *http.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		srv     *httptest.Server
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
			if tt.name == "positive test #1" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusAccepted)
						}
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.UploadUserData(context.Background(), []models.UserData{{ID: "secret pass"}}); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.UploadUserData() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #1" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						w.WriteHeader(http.StatusBadRequest)
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.UploadUserData(context.Background(), []models.UserData{{ID: "secret pass"}}); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.UploadUserData() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #2" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusUnauthorized)
						}
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.UploadUserData(context.Background(), []models.UserData{{ID: "secret pass"}}); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.UploadUserData() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.name == "negative test #3" {
				s := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-type", "text/plain")
						if r.RequestURI == "/ping" {
							w.WriteHeader(http.StatusOK)
						} else {
							w.WriteHeader(http.StatusInternalServerError)
						}
					}))
				defer func() { s.Close() }()
				cfg := &config.Config{AppAddr: s.URL + "/"}
				h := New(cfg, "")
				if err := h.UploadUserData(context.Background(), []models.UserData{{ID: "secret pass"}}); (err != nil) != tt.wantErr {
					t.Errorf("HTTP.UploadUserData() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
