// Package service implements system facade.
package service

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xEgorka/project3/cmd/client/internal/app/auth"
	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/logger"
	"github.com/xEgorka/project3/cmd/client/internal/app/models"
	"github.com/xEgorka/project3/cmd/client/internal/app/requests"
	"github.com/xEgorka/project3/cmd/client/internal/app/storage"
)

// Service provides business logic.
type Service struct {
	cfg *config.Config
	s   storage.Storage
	a   auth.Auth
	r   *requests.HTTP
}

// New creates Service.
func New(config *config.Config, store storage.Storage, auth auth.Auth,
	req *requests.HTTP) *Service {
	return &Service{cfg: config, s: store, a: auth, r: req}
}

// Register registers new user.
func (s *Service) Register(ctx context.Context, usr, pass string) error {
	if err := s.r.Register(ctx, usr, pass); err != nil {
		return err
	}
	if err := s.a.LoginOnline(ctx, usr, pass); err != nil {
		return err
	}
	logger.Log.Info("new user registered", zap.String("usr", s.cfg.Usr))
	return nil
}

// Login signs user in.
func (s *Service) Login(ctx context.Context, usr, pass string) error {
	if err := s.r.Login(ctx, usr, pass); err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.Unavailable {
				if ee := s.a.LoginOffline(ctx, usr, pass); ee != nil {
					return status.Error(codes.Unauthenticated, "unauthenticated")
				}
				return err
			}
		}
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}
	if err := s.a.LoginOnline(ctx, usr, pass); err != nil {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}
	logger.Log.Info("user logged in", zap.String("usr", s.cfg.Usr))
	return nil
}

// DownloadUserData downloads merged data to client.
func (s *Service) DownloadUserData(ctx context.Context) error {
	timestamp, err := s.s.GetTimestamp(ctx)
	if err != nil {
		return err
	}
	d, err := s.r.DownloadUserData(ctx, timestamp)
	if err != nil {
		return err
	}
	if len(d) > 0 {
		logger.Log.Info("dowloading user data", zap.Int("count", len(d)))
		if err := s.s.UpdateUserData(ctx, d); err != nil {
			return err
		}
	}
	return nil
}

// UploadUserData uploads unmerged data to server.
func (s *Service) UploadUserData(ctx context.Context) error {
	d, err := s.s.GetUnmergedUserData(ctx, s.cfg.Usr)
	if err != nil {
		return err
	}
	if len(d) > 0 {
		logger.Log.Info("uploading user data", zap.Int("count", len(d)))
		if err := s.r.UploadUserData(ctx, d); err != nil {
			return err
		}
	}
	return nil
}

// IndexUserData provides index for user data.
func (s *Service) IndexUserData(ctx context.Context) ([]models.UserData, error) {
	i, err := s.s.IndexUserData(ctx)
	if err != nil {
		return nil, err
	}
	if len(i) == 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return i, nil
}

// PrintUserData prints user data from database.
func (s *Service) PrintUserData(ctx context.Context, id string) (
	*models.UserData, error) {
	d, err := s.s.PrintUserData(ctx, id)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// EnterUserData saves user data to database.
func (s *Service) EnterUserData(ctx context.Context, id string, t int, c string) error {
	if err := s.s.EnterUserData(ctx, &models.UserData{
		UserID: s.cfg.Usr, ID: id, Type: t, Data: c}); err != nil {
		return err
	}
	return nil
}
