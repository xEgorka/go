// Package service implements system facade.
package service

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xEgorka/project3/internal/app/auth"
	"github.com/xEgorka/project3/internal/app/config"
	"github.com/xEgorka/project3/internal/app/jobs"
	"github.com/xEgorka/project3/internal/app/logger"
	"github.com/xEgorka/project3/internal/app/models"
	"github.com/xEgorka/project3/internal/app/storage"
)

// Service provides business logic.
type Service struct {
	cfg *config.Config
	s   storage.Storage
	j   *jobs.Jobs
	a   auth.Auth
}

// New creates Service.
func New(config *config.Config, store storage.Storage, jobs *jobs.Jobs,
	auth auth.Auth) *Service {
	return &Service{cfg: config, s: store, j: jobs, a: auth}
}

// Register registers new user.
func (s *Service) Register(ctx context.Context, usr, pass string) error {
	if err := s.a.Register(ctx, usr, pass); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return status.Error(codes.AlreadyExists, "already exists")
		}
		logger.Log.Info("failed register usr", zap.Error(err))
		return status.Error(codes.Unknown, "unknown")
	}
	return nil
}

// Login logs user in.
func (s *Service) Login(ctx context.Context, usr, pass string) error {
	if err := s.a.Login(ctx, usr, pass); err != nil {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return nil
}

// GetToken provides token for user.
func (s *Service) GetToken(usr string) (string, error) {
	if token, err := auth.BuildJWTString(usr); err != nil {
		logger.Log.Info("internal server error", zap.Error(err))
		return "", status.Error(codes.Internal, "internal")
	} else {
		return token, nil
	}
}

// Ping checks storage availability.
func (s *Service) Ping() error { return s.s.Ping() }

// GetUserData provides user data merged after specified timestamp.
func (s *Service) GetUserData(ctx context.Context,
	usr string, timestamp int) ([]models.UserData, error) {
	res, err := s.s.GetUserData(ctx, usr, timestamp)
	if err != nil {
		logger.Log.Info("internal server error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal")
	}
	if len(res) == 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return res, nil
}

// MergeUserData calls jobs to merge user data.
func (s *Service) MergeUserData(r []models.UserData) { s.j.MergeUserData(r) }
