// Package jobs implements fan-in concurrency pattern for batch
// database update.
package jobs

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/xEgorka/project3/internal/app/logger"
	"github.com/xEgorka/project3/internal/app/models"
	"github.com/xEgorka/project3/internal/app/storage"
)

type config struct {
	Period    time.Duration
	BatchSize int
}

// Jobs provides entities required for batch job.
type Jobs struct {
	store         storage.Storage
	wg            *sync.WaitGroup
	dataToMergeCh chan models.UserData
	cfg           *config
	isRun         bool
}

// New creates Jobs.
func New(store storage.Storage) *Jobs {
	return &Jobs{
		store:         store,
		wg:            &sync.WaitGroup{},
		dataToMergeCh: make(chan models.UserData),
		cfg:           &config{Period: 1 * time.Second, BatchSize: 10},
	}
}

// Run starts and stops batch job.
func (j *Jobs) Run(ctx context.Context) error {
	if j.isRun {
		return nil
	}
	go func() {
		<-ctx.Done()
		j.stop()
	}()
	once := sync.Once{}
	once.Do(func() {
		j.jobServer(ctx)
		j.isRun = true
	})
	return nil
}

func (j *Jobs) stop() { j.wg.Wait(); close(j.dataToMergeCh) }

func (j *Jobs) jobServer(ctx context.Context) {
	go func() {
		buf := make([]models.UserData, 0, j.cfg.BatchSize)
		defer func() {
			logger.Log.Info("batch stopping...")
			if len(buf) > 0 {
				logger.Log.Info("merging remaining data...",
					zap.Int("count", len(buf)))
				if err := j.store.MergeUserData(ctx, buf); err != nil {
					logger.Log.Error("data merge failed", zap.Error(err))
				}
				logger.Log.Info("merging remaining data... done")
			}
			defer logger.Log.Info("batch stopping... done")
		}()
		ticker := time.NewTicker(j.cfg.Period)
		defer ticker.Stop()
		logger.Log.Info("running batch...")

		for {
			select {
			case <-ticker.C:
				if err := j.store.MergeUserData(ctx, buf); err != nil {
					continue
				}
				buf = buf[:0] // clear buf on success
			case dataUpdateReq, ok := <-j.dataToMergeCh:
				if !ok { // chan closed, exit
					return
				}
				buf = append(buf, dataUpdateReq)
				if len(buf) < j.cfg.BatchSize {
					continue
				}
				if err := j.store.MergeUserData(ctx, buf); err != nil {
					logger.Log.Error("data update failed", zap.Error(err))
					continue
				}
				buf = buf[:0] // clear buf on success
			}
		}
	}()
}

// MergeUserData adds data to channel for subsequent merge batch processing.
func (j *Jobs) MergeUserData(r []models.UserData) {
	j.wg.Add(1)
	go func() {
		defer j.wg.Done()
		for _, data := range r {
			j.dataToMergeCh <- data
		}
	}()
}
