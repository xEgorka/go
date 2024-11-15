package cli

import (
	"bufio"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"text/tabwriter"

	"github.com/golang/mock/gomock"

	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/crypto"
	"github.com/xEgorka/project3/cmd/client/internal/app/mocks"
	"github.com/xEgorka/project3/cmd/client/internal/app/models"
	"github.com/xEgorka/project3/cmd/client/internal/app/requests"
	"github.com/xEgorka/project3/cmd/client/internal/app/service"
)

func TestApp_main(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)
					f1, err := os.Open("./testdata/case01.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f1),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					ms.EXPECT().PrintUserData(tt.args.ctx, "mm").
						Return(nil, errors.New("test"))
					a.main(tt.args.ctx)
					f1.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main1(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{
						AppAddr: srv.URL + "/",
						Token:   "testToken"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)
					f1, err := os.Open("./testdata/case01.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f1),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					ms.EXPECT().PrintUserData(tt.args.ctx, "mm").
						Return(nil, errors.New("test"))
					a.main(tt.args.ctx)
					f1.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main1")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main2(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)
					f1, err := os.Open("./testdata/case02.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f1),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					ms.EXPECT().PrintUserData(tt.args.ctx, "ff").
						Return(nil, errors.New("test"))
					a.main(tt.args.ctx)
					f1.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main2")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main3(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{
						AppAddr:       srv.URL + "/",
						FileSizeLimit: 10000000}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)
					f1, err := os.Open("./testdata/case03.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f1),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					ms.EXPECT().PrintUserData(tt.args.ctx, "ff").
						Return(nil, errors.New("test"))
					a.main(tt.args.ctx)
					f1.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main3")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main31(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{
						AppAddr:       srv.URL + "/",
						FileSizeLimit: 10000000}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)
					f1, err := os.Open("./testdata/case031.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f1),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					ms.EXPECT().PrintUserData(tt.args.ctx, "ff").
						Return(nil, errors.New("test"))
					a.main(tt.args.ctx)
					f1.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main31")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main32(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{
						AppAddr:       srv.URL + "/",
						FileSizeLimit: 1}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)
					f1, err := os.Open("./testdata/case03.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f1),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					ms.EXPECT().PrintUserData(tt.args.ctx, "ff").
						Return(nil, errors.New("test"))
					a.main(tt.args.ctx)
					f1.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main32")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main4(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)
					f1, err := os.Open("./testdata/case04.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f1),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					ms.EXPECT().PrintUserData(tt.args.ctx, "ff").
						Return(nil, errors.New("test"))
					a.main(tt.args.ctx)
					f1.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main4")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main41(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)
					f1, err := os.Open("./testdata/case041.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f1),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					ms.EXPECT().PrintUserData(tt.args.ctx, "ff").
						Return(nil, errors.New("test"))
					a.main(tt.args.ctx)
					f1.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main41")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main42(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)
					f1, err := os.Open("./testdata/case042.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f1),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					ms.EXPECT().PrintUserData(tt.args.ctx, "ff").
						Return(nil, errors.New("test"))
					a.main(tt.args.ctx)
					f1.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main42")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main43(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)
					f1, err := os.Open("./testdata/case043.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f1),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					ms.EXPECT().PrintUserData(tt.args.ctx, "ff").
						Return(nil, errors.New("test"))
					a.main(tt.args.ctx)
					f1.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main43")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main5(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)
					f1, err := os.Open("./testdata/case05.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f1),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					ms.EXPECT().PrintUserData(tt.args.ctx, "ff").
						Return(nil, errors.New("test"))
					a.main(tt.args.ctx)
					f1.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main5")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main6(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)

					f, err := os.Open("./testdata/case08.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					a.cfg.Key = a.c.Key("x0o1@ya.ru", "e")
					ms.EXPECT().IndexUserData(tt.args.ctx).
						Return([]models.UserData{{
							UserID: "x0o1@ya.ru",
							ID:     "ee", Type: 1}}, nil)
					ms.EXPECT().PrintUserData(tt.args.ctx, "ee").
						Return(&models.UserData{
							Type: 1,
							Data: "96763bd71acd3d00fe948d32f60e4a6dea3e154f04ccbf83fc10e1483cf9c5e88bae9a65210df86573eb2a8a7393d7bc79e22b5adf1e5d3ad54f76e8fa58f80577f34ab6ad6974"}, nil)
					a.main(tt.args.ctx)
					f.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main6")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main7(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)

					f, err := os.Open("./testdata/case08.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					a.cfg.Key = a.c.Key("x0o1@ya.ru", "e")
					ms.EXPECT().IndexUserData(tt.args.ctx).
						Return([]models.UserData{{
							UserID: "x0o1@ya.ru",
							ID:     "ee",
							Type:   1}}, nil)
					ms.EXPECT().PrintUserData(tt.args.ctx, "ee").
						Return(&models.UserData{
							Type: 2,
							Data: "438299998103c8b6e110d9a9a79667b6aba2808986b6d7bd707ba2bb4dfde7254c12202ffd4f527b290f5faf05291250540bc8"}, nil)
					a.main(tt.args.ctx)
					f.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main7")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main8(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)

					f, err := os.Open("./testdata/case08.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					a.cfg.Key = a.c.Key("x0o1@ya.ru", "e")
					ms.EXPECT().IndexUserData(tt.args.ctx).
						Return([]models.UserData{{
							UserID: "x0o1@ya.ru",
							ID:     "ee",
							Type:   1}}, nil)
					ms.EXPECT().PrintUserData(tt.args.ctx, "ee").
						Return(&models.UserData{
							Type: 3,
							Data: "1efdce06857aa968d5dd0f456cc69c65254d9bafa2ed608684d905e3a41e1f005a247ab2f73653a8416d22a87ff5d44f53de5ba64be4d274b2cec40690f662493aaaad1ea9d62c96ecda03193a87f6421ff07fa1c4cc70042c8f0d92763d3b2e121739a721aae691460afd2968d5b301d22bf06140dc9438e5f8491e486dee53bcb81cf94cfda5c8ccd29b8a50fcedc1baa13aae7b8412bd7888a987e7d71ecdaa3a652a83d1db5122afd93cf97e32b417c6a2d27c8189935213081d411267813c4c40c62667dda165c16a747338147dce79d8d6dc6b635098c2a30c604f7c1c716f386a2bec2a62707bc2e0aa0a0ab6eb9d57753f3dec09e7eacd96bb59554c6c318cfa0ba75163355723c0662f06a8bbd8c0519cfd00e5735e258cadef5d00099e1136180b07c6381778864793e636d73c499e6ce73886dc0b2b3d8d7604ba40f38266f3e59c739fa57595ce663f2c3f6e556044c7df2e5071e35659841783a8788b72222eaffb17b006f8ac3e5515f6cb5bdd1ad7d09a9bf8dd8cbd2fcbe95479a6622180be935900be3b537894633b1e69820632558042d06866b295ff17a65bfff1e42271c070885c8b37dce694bc576d928489425bdf6f119d4407fa94db8c2f0546632501ce0796d5ce350c555b139fc3aecc2ddf5b0e65d861041dbbb3155dace1cb876a3879aa9ebc09843461b61c5bcc31a62aab0b9bca56ee867f5e2c7360c5f24a37ee03ef73c3a5cbadb81b111ed53ccd8307d0e6eda79b8985a702874b7aaf047739fd0e4903fe0a531c534417e36236effdefce57fad097dc1009497c8f1b4ed1080ecd845c289e8bb4ab4e67e65125214aac0b60722e3eec90a362c3e525ff8928566af16569ac6bbe8226b5663dac197f1651eb11c7ae79a8d390f88a85773b5ebe2237af432cd45c345bf27349eb76409bc56027716ee4c1010f173b28fe33d5530f6da21b949385bf3b7ea33b86ae6c7178d212648b5e6e599dab8e02c6d6675c0b96f356c9ee4e70989eeaba0c079e1ae2b5c67095f9c8570754b119d8b7f204fd4602ccfc442d6edf9efdf989ed920b048e69a6ec3231110bf8c73971687023485b1acb322bfb216bf231edd3b946a4140faf1df3b9e117996336b877153e9b07137a93110ccabe6bc4c18cdf55d46ba4a75d5a60d8b9b1b0db94a22ff2f135039acc5464b9adec667e08db9a3204e2c97c0c283077059d65359d5117d21479dd87d997e8661a9d71b2a4379ba2cc364124c623cde84aecc35bb323838b462f27ff29ff41336c0ebf1bda2f9ef112c48b9f3ed81d2c7904257a98cff4bdfde7a06afd047b719157eadd86dd116b9e001b97393a8d5710d7995c283d07405e32a8beaa8e5d04acfa23123c13438f00b856ffe59d88737c6b979fe4ff7dabbc6a2449610c4186b9bba2c815c915ba29fc99f238dbb85b1513d7678ff2cf5da4b424ccccb4f3e93154d771375aacb27df259578e7ce0ceadc11620a28abb07b5b2d95381a08fd835ea28d2c51bf8e6a42993226aef66e17d9c8af7dfe95e9cc653ff594d2cd2430816a3fae023246d89e4e3d5132cb70f64eddcf1fa7b50b32ca26f45898037d38a177df2a6b84c19e72dcf2c6e15f908c5796f74a3e0ed7e72dda677fe4d3c42b5e3f94f020e97e0ddc64ee580880b6eae464ad54fac2445a63cb3e3675d41998cc79719cdafca41bf61578cb4e51fbcc328409556dd07cf1670e67d984037ce60e2b21afb2cadfd476ffcec46ef27c346cb49d41e49c3a712319ab129f8c2597ee790017d921d5d8329e5dff11da28ced42ec361fe4cbd228e938cd9c5c02bf0780734a5bbf35b69ede76b92b107062f92cbd9822787f1bf952d55d281cddfc06467f6534a1fdd6631f334f6bcf1fe4cb0a889157c5113ed5112694ff4999467a4a4b952589c7af8c3286570e64594786ef4ea4c52521ec446a70f02a850555f67f394ef7f396d9015c34fd4067bf365015673f9fbda823bd5049007ad5406d55326c66f3443393a8a487747edb6a7ca5a9"}, nil)
					a.main(tt.args.ctx)
					f.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main8")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main9(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)

					f, err := os.Open("./testdata/case08.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					a.cfg.Key = a.c.Key("x0o1@ya.ru", "e")
					ms.EXPECT().IndexUserData(tt.args.ctx).
						Return([]models.UserData{{
							UserID: "x0o1@ya.ru",
							ID:     "ee",
							Type:   1}}, nil)
					ms.EXPECT().PrintUserData(tt.args.ctx, "ee").
						Return(&models.UserData{
							Type: 4,
							Data: "36395a3581c75911337596dacbae133dbbaad303d00a1e8a86b803a74a2b7bb71fb603b1fbf32b2df655f9f7b2b9af90daa33cbda7c383ead5870602c85e0155a919c016b37b8063f437b64c4f5cf14391124ba2f85bdc"}, nil)
					a.main(tt.args.ctx)
					f.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main9")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestApp_main91(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
					cfg := &config.Config{AppAddr: srv.URL + "/"}
					r := requests.New(cfg, "")
					s := service.New(cfg, ms, ma, r)

					f, err := os.Open("./testdata/case08.txt")
					if err != nil {
						panic(err)
					}
					a := &App{
						cfg: cfg,
						s:   s,
						c:   crypto.New(),
						r:   bufio.NewScanner(f),
						w:   bufio.NewWriter(os.Stdout),
						tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
					}
					a.cfg.Key = a.c.Key("x0o1@ya.ru", "e")
					ms.EXPECT().IndexUserData(tt.args.ctx).
						Return([]models.UserData{{
							UserID: "x0o1@ya.ru",
							ID:     "ee",
							Type:   1}}, nil)
					ms.EXPECT().PrintUserData(tt.args.ctx, "ee").
						Return(&models.UserData{
							Type: 4,
							Data: "bad36395a3581c75911337596dacbae133dbbaad303d00a1e8a86b803a74a2b7bb71fb603b1fbf32b2df655f9f7b2b9af90daa33cbda7c383ead5870602c85e0155a919c016b37b8063f437b64c4f5cf14391124ba2f85bdc"}, nil)
					a.main(tt.args.ctx)
					f.Close()
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestApp_main91")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Run()
		})
	}
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "positive test #1",
		},
		{
			name: "negative test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "positive test #1" {
				x, _ := unmarshal[models.Text](`{"data": "user text"}`)
				if x.Text != "user text" {
					t.Errorf("unmarshal error")
				}
			}
			if tt.name == "negative test #1" {
				if _, err := unmarshal[models.Text](`bad{"data": "user text"}`); (err != nil) != true {
					t.Errorf("error = %v, wantErr %v", err, true)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	ma := mocks.NewMockAuth(ctrl)
	cfg := &config.Config{}
	r := requests.New(cfg, "")
	s := service.New(cfg, ms, ma, r)
	tests := []struct {
		name string
		want App
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(cfg, s, crypto.New())
			if reflect.TypeOf(got) == reflect.TypeOf((*App)(nil)).Elem() {
				t.Errorf("not app")
			}
		})
	}
}

func TestApp_auth(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
				cfg := &config.Config{AppAddr: srv.URL + "/"}
				r := requests.New(cfg, "")
				s := service.New(cfg, ms, ma, r)
				f, err := os.Open("./testdata/case06.txt")
				if err != nil {
					panic(err)
				}
				a := &App{
					cfg: cfg,
					s:   s,
					c:   crypto.New(),
					r:   bufio.NewScanner(f),
					w:   bufio.NewWriter(os.Stdout),
					tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
				}
				ma.EXPECT().LoginOnline(tt.args.ctx, "x0o1@ya.ru", "e").Return(nil)
				a.auth(tt.args.ctx)
				f.Close()
			}
		})
	}
}

func TestApp_auth2(t *testing.T) {
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
			if tt.name == "positive test #1" {
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
				cfg := &config.Config{AppAddr: srv.URL + "/"}
				r := requests.New(cfg, "")
				s := service.New(cfg, ms, ma, r)
				f, err := os.Open("./testdata/case07.txt")
				if err != nil {
					panic(err)
				}
				a := &App{
					cfg: cfg,
					s:   s,
					c:   crypto.New(),
					r:   bufio.NewScanner(f),
					w:   bufio.NewWriter(os.Stdout),
					tw:  tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
				}
				ma.EXPECT().LoginOnline(tt.args.ctx, "x0o1@ya.ru", "e").Return(nil)
				a.auth(tt.args.ctx)
				f.Close()
			}
		})
	}
}
