// Package requests makes requests to server for user data syncronization.
package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/logger"
	"github.com/xEgorka/project3/cmd/client/internal/app/models"
)

// HTTP provides methods for http requests.
type HTTP struct {
	cfg *config.Config
	c   *http.Client
}

// New creates HTTP.
func New(config *config.Config, caCert string) *HTTP {
	return &HTTP{cfg: config, c: newClient([]byte(caCert))}
}

func newClient(caCert []byte) *http.Client {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
}

func (h *HTTP) ping(ctx context.Context) error {
	url := h.cfg.AppAddr + "ping"
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "text/plain")
	res, err := h.c.Do(r)
	if err != nil {
		return err
	}
	defer func() {
		if e := res.Body.Close(); e != nil {
			logger.Log.Info("failed body close", zap.Error(e))
		}
	}()
	if res.StatusCode != http.StatusOK {
		return status.Error(codes.Unavailable, "unavailable")
	}
	return nil
}

// Register requests register new user.
func (h *HTTP) Register(ctx context.Context, usr, pass string) error {
	if err := h.ping(ctx); err != nil {
		return status.Error(codes.Unavailable, "unavailable")
	}
	url := h.cfg.AppAddr + "api/user/register"
	d := models.RequestAuth{
		Usr:  usr,
		Pass: pass,
	}
	j, err := json.Marshal(d)
	if err != nil {
		return err
	}
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	res, err := h.c.Do(r)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Log.Info("failed body close", zap.Error(err))
		}
	}()

	switch {
	case res.StatusCode == http.StatusConflict:
		return status.Error(codes.AlreadyExists, "already exists")
	case res.StatusCode != http.StatusOK:
		return status.Error(codes.InvalidArgument, "invalid argument")
	}
	for _, v := range res.Cookies() {
		if v.Name == "token" {
			h.cfg.Token = v.Value
		}
	}
	if len(h.cfg.Token) == 0 {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return nil
}

// Login requests authentication on server.
func (h *HTTP) Login(ctx context.Context, usr, pass string) error {
	if err := h.ping(ctx); err != nil {
		return status.Error(codes.Unavailable, "unavailable")
	}
	url := h.cfg.AppAddr + "api/user/login"
	d := models.RequestAuth{
		Usr:  usr,
		Pass: pass,
	}
	j, err := json.Marshal(d)
	if err != nil {
		return err
	}
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	res, err := h.c.Do(r)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Log.Info("failed body close", zap.Error(err))
		}
	}()

	if res.StatusCode != http.StatusOK {
		return status.Error(codes.InvalidArgument, "invalid argument")
	}
	for _, v := range res.Cookies() {
		if v.Name == "token" {
			h.cfg.Token = v.Value
		}
	}
	if len(h.cfg.Token) == 0 {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return nil
}

// DownloadUserData downloads user data from server.
func (h *HTTP) DownloadUserData(ctx context.Context, timestamp string) (
	[]models.UserData, error) {
	if err := h.ping(ctx); err != nil {
		return nil, err
	}
	url := h.cfg.AppAddr + "api/user/data/" + timestamp
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "text/plain")
	r.AddCookie(&http.Cookie{Name: "token", Value: h.cfg.Token})
	res, err := h.c.Do(r)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := res.Body.Close(); e != nil {
			logger.Log.Info("failed body close", zap.Error(e))
		}
	}()
	switch {
	case res.StatusCode == http.StatusNoContent:
		return nil, nil
	case res.StatusCode == http.StatusUnauthorized:
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	case res.StatusCode != http.StatusOK:
		return nil, errors.New("not ok")
	}

	var d []models.UserData
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// UploadUserData uploads user data to server.
func (h *HTTP) UploadUserData(ctx context.Context, d []models.UserData) error {
	if err := h.ping(ctx); err != nil {
		return err
	}
	j, err := json.Marshal(d)
	if err != nil {
		return err
	}
	url := h.cfg.AppAddr + "api/user/data"
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json")
	r.AddCookie(&http.Cookie{Name: "token", Value: h.cfg.Token})
	res, err := h.c.Do(r)
	if err != nil {
		return err
	}
	defer func() {
		if e := res.Body.Close(); e != nil {
			logger.Log.Info("failed body close", zap.Error(e))
		}
	}()
	switch {
	case res.StatusCode == http.StatusUnauthorized:
		return status.Error(codes.Unauthenticated, "unauthenticated")
	case res.StatusCode != http.StatusAccepted:
		return errors.New("not accepted")
	}
	return nil
}
