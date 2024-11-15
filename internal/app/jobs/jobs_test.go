package jobs

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/xEgorka/project3/internal/app/mocks"
	"github.com/xEgorka/project3/internal/app/models"
	"github.com/xEgorka/project3/internal/app/storage"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	type args struct {
		store storage.Storage
	}
	tests := []struct {
		name string
		args args
		want *Jobs
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(ms)
			if reflect.TypeOf(got) == reflect.TypeOf((*Jobs)(nil)).Elem() {
				t.Errorf("not jobs")
			}
		})
	}
}

func TestJobs_Run(t *testing.T) {
	type fields struct {
		store         storage.Storage
		wg            *sync.WaitGroup
		dataToMergeCh chan models.UserData
		cfg           *config
		isRun         bool
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
			fields:  fields{isRun: true},
			wantErr: false,
		},
		{
			name:    "positive test #2",
			fields:  fields{isRun: false, cfg: &config{Period: 1 * time.Second, BatchSize: 10}},
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Jobs{
				store:         tt.fields.store,
				wg:            tt.fields.wg,
				dataToMergeCh: tt.fields.dataToMergeCh,
				cfg:           tt.fields.cfg,
				isRun:         tt.fields.isRun,
			}
			if err := j.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Jobs.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJobs_jobServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStorage(ctrl)
	type fields struct {
		store         storage.Storage
		wg            *sync.WaitGroup
		dataToMergeCh chan models.UserData
		cfg           *config
		isRun         bool
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
			fields: fields{isRun: true, cfg: &config{
				Period: 1 * time.Second, BatchSize: 10},
				dataToMergeCh: make(chan models.UserData)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Jobs{
				store:         ms,
				wg:            &sync.WaitGroup{},
				dataToMergeCh: tt.fields.dataToMergeCh,
				cfg:           tt.fields.cfg,
				isRun:         tt.fields.isRun,
			}
			close(j.dataToMergeCh)
			j.jobServer(tt.args.ctx)
		})
	}
}

func TestJobs_MergeUserData(t *testing.T) {
	type fields struct {
		store         storage.Storage
		wg            *sync.WaitGroup
		dataToMergeCh chan models.UserData
		cfg           *config
		isRun         bool
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
			name:   "positive test #1",
			args:   args{r: []models.UserData{{ID: "testID"}}},
			fields: fields{isRun: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Jobs{
				store:         tt.fields.store,
				wg:            &sync.WaitGroup{},
				dataToMergeCh: make(chan models.UserData),
				cfg:           tt.fields.cfg,
				isRun:         tt.fields.isRun,
			}
			j.MergeUserData(tt.args.r)
		})
	}
}

func TestJobs_stop(t *testing.T) {
	type fields struct {
		store         storage.Storage
		wg            *sync.WaitGroup
		dataToMergeCh chan models.UserData
		cfg           *config
		isRun         bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "positive test #1",
			fields: fields{isRun: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Jobs{
				store:         tt.fields.store,
				wg:            &sync.WaitGroup{},
				dataToMergeCh: make(chan models.UserData),
				cfg:           tt.fields.cfg,
				isRun:         tt.fields.isRun,
			}
			j.stop()
		})
	}
}
