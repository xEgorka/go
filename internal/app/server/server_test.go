package server

import (
	"context"
	_ "embed"
	"net/http"
	"reflect"
	"syscall"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	"github.com/xEgorka/project3/internal/app/config"
	"github.com/xEgorka/project3/internal/app/handlers"
	"github.com/xEgorka/project3/internal/app/jobs"
	"github.com/xEgorka/project3/internal/app/mocks"
	"github.com/xEgorka/project3/internal/app/service"
)

func Test_stop(t *testing.T) {
	_, cancelBatch := context.WithCancel(context.Background())
	srv := http.Server{}
	go srv.ListenAndServe()
	type args struct {
		cancelBatch context.CancelFunc
		srv         *http.Server
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{cancelBatch: cancelBatch, srv: &srv},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				if err := stop(tt.args.cancelBatch, tt.args.srv); (err != nil) != tt.wantErr {
					t.Errorf("stop() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			sigint <- syscall.SIGQUIT
		})
	}
}

func Test_routes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := service.New(&config.Config{}, ms, j, ma)
	h := handlers.NewHTTP(s)
	type args struct {
		h handlers.HTTP
	}
	tests := []struct {
		name string
		args args
		want *chi.Mux
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := routes(h)
			if reflect.TypeOf(got) == reflect.TypeOf((*chi.Mux)(nil)).Elem() {
				t.Errorf("not chi mux")
			}
		})
	}
}

func TestStart(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "negative test #1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
