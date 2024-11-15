package cli

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xEgorka/project3/cmd/client/internal/app/logger"
)

const period = 1 * time.Minute

func (a *App) sync(ctx context.Context) {
	for {
		a.download(ctx)
		a.upload(ctx)
		time.Sleep(period)
	}
}

func (a *App) download(ctx context.Context) {
	if err := a.s.DownloadUserData(ctx); err != nil {
		logger.Log.Info("error download data", zap.Error(err))
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.Unauthenticated {
				a.exitTimeout()
			}
		}
	}
}

func (a *App) upload(ctx context.Context) {
	if err := a.s.UploadUserData(ctx); err != nil {
		logger.Log.Info("error upload data", zap.Error(err))
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.Unauthenticated {
				a.exitTimeout()
			}
		}
	}
}
