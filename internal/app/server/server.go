// Package server starts storage, batch job and defines HTTP router.
package server

import (
	"bufio"
	"context"
	"crypto/tls"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/xEgorka/project3/internal/app/auth"
	"github.com/xEgorka/project3/internal/app/config"
	"github.com/xEgorka/project3/internal/app/crypto"
	"github.com/xEgorka/project3/internal/app/handlers"
	"github.com/xEgorka/project3/internal/app/jobs"
	"github.com/xEgorka/project3/internal/app/logger"
	"github.com/xEgorka/project3/internal/app/service"
	"github.com/xEgorka/project3/internal/app/storage"
)

//go:embed server.key
var key string

//go:embed server.crt
var crt string

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

// Start starts server.
func Start() error {
	w := bufio.NewWriter(os.Stdout)
	if _, err := fmt.Fprintf(w, "GOPHKEEPER PASSWORD MANAGER SERVER\nVersion: %s\t%s\nCommit: %s\n\n",
		buildVersion, buildDate, buildCommit); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	if e := logger.Initialize("info"); e != nil {
		return e
	}
	cfg, err := config.Setup()
	if err != nil {
		return err
	}
	ctx := context.Background()
	s, err := storage.Open(ctx, cfg)
	if err != nil {
		return err
	}
	a, err := auth.Open(ctx, cfg, crypto.New())
	if err != nil {
		return err
	}
	j := jobs.New(s)
	ctxBatch, cancelBatch := context.WithCancel(ctx)
	if err := j.Run(ctxBatch); err != nil {
		logger.Log.Error("failed run batch job", zap.Error(err))
		cancelBatch()
		return err
	}
	srv := http.Server{
		Addr:    cfg.URI,
		Handler: routes(handlers.NewHTTP(service.New(cfg, s, j, a))),
	}

	go func() {
		logger.Log.Info("running https server...", zap.String("uri", cfg.URI))
		if pair, err := tls.X509KeyPair([]byte(crt), []byte(key)); err != nil {
			logger.Log.Error("failed run https server", zap.Error(err))
			sigint <- syscall.SIGQUIT
		} else {
			srv.TLSConfig = &tls.Config{Certificates: []tls.Certificate{pair}}
			if err := srv.ListenAndServeTLS("", ""); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					logger.Log.Info("https server stopping... done")
				} else {
					logger.Log.Error("failed run https server", zap.Error(err))
				}
			}
		}
	}()
	return stop(cancelBatch, &srv)
}

var sigint = make(chan os.Signal, 1)

const timeout = 5 * time.Second

func stop(cancelBatch context.CancelFunc, srv *http.Server) error {
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-sigint
	logger.Log.Info("signal received", zap.String("sig", sig.String()))

	cancelBatch()
	logger.Log.Info("server stopping...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("failed server stop", zap.Error(err))
		return err
	}
	return nil
}

func routes(h handlers.HTTP) *chi.Mux {
	r := chi.NewRouter()
	r.Use(handlers.WithLogging)
	r.Get("/ping", h.GetPing)
	r.Post("/api/user/register", h.PostUserRegister)
	r.Post("/api/user/login", h.PostUserLogin)
	r.Group(func(r chi.Router) {
		r.Use(handlers.Auth)
		r.Get("/api/user/data/{timestamp}", h.GetUserData)
		r.Post("/api/user/data", h.PostUserData)
	})
	return r
}
