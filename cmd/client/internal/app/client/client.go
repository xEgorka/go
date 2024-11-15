// Package client starts client.
package client

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/xEgorka/project3/cmd/client/internal/app/auth"
	"github.com/xEgorka/project3/cmd/client/internal/app/cli"
	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/crypto"
	"github.com/xEgorka/project3/cmd/client/internal/app/logger"
	"github.com/xEgorka/project3/cmd/client/internal/app/requests"
	"github.com/xEgorka/project3/cmd/client/internal/app/service"
	"github.com/xEgorka/project3/cmd/client/internal/app/storage"
)

//go:embed server.crt
var crt string

var (
	buildVersion = "0.1"
	buildDate    = "2024-06-18"
)

// Start starts client.
func Start() error {
	w := bufio.NewWriter(os.Stdout)
	if _, err := fmt.Fprintf(w, "GOPHKEEPER PASSWORD MANAGER\nVersion: %s\t%s\n\n",
		buildVersion, buildDate); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	if e := logger.Initialize("error"); e != nil {
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
	c := crypto.New()
	a, err := auth.Open(ctx, cfg, c)
	if err != nil {
		return err
	}
	app := cli.New(cfg, service.New(cfg, s, a, requests.New(cfg, crt)), c)
	return app.Run(ctx)
}
