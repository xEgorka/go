package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xEgorka/project3/internal/app/logger"
	"github.com/xEgorka/project3/internal/app/models"
	"github.com/xEgorka/project3/internal/app/service"
)

// HTTP provides methods for http server.
type HTTP struct{ s *service.Service }

// NewHTTP creates HTTP.
func NewHTTP(service *service.Service) HTTP { return HTTP{s: service} }

// GetPing checks service availability.
func (h *HTTP) GetPing(w http.ResponseWriter, r *http.Request) {
	if err := h.s.Ping(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "text/plain")
}

// PostUserRegister handels user register requests.
func (h *HTTP) PostUserRegister(w http.ResponseWriter, r *http.Request) {
	var req models.RequestAuth
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("JSON decode error", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	userID, pass := req.Usr, req.Pass
	if len(userID) == 0 || len(pass) == 0 {
		logger.Log.Info("empty userID or pass")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.s.Register(r.Context(), userID, pass); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unknown:
				logger.Log.Info("unknown error", zap.Error(err))
				http.Error(w, "Bad Request", http.StatusBadRequest)
			case codes.AlreadyExists:
				w.WriteHeader(http.StatusConflict)
			}
			return
		}
	}
	token, err := h.s.GetToken(userID)
	if err != nil {
		logger.Log.Info("unable to get token", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	c := &http.Cookie{
		Name:  "token",
		Value: token,
	}
	w.Header().Set("Content-type", "text/plain")
	http.SetCookie(w, c)
}

// PostUserLogin handels user login requests.
func (h *HTTP) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	var req models.RequestAuth
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("JSON decode error", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	userID, pass := req.Usr, req.Pass
	if len(userID) == 0 || len(pass) == 0 {
		logger.Log.Info("empty userID or pass")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if err := h.s.Login(r.Context(), userID, pass); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := h.s.GetToken(userID)
	if err != nil {
		logger.Log.Info("unable to get token", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	c := &http.Cookie{
		Name:  "token",
		Value: token,
	}
	http.SetCookie(w, c)
	w.Header().Set("Content-type", "text/plain")
}

var userKey userIDKeyType

// GetUserData handels requests of user data.
func (h *HTTP) GetUserData(w http.ResponseWriter, r *http.Request) {
	timestamp, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	ctx := r.Context()
	res, err := h.s.GetUserData(ctx, ctx.Value(userKey).(string), timestamp)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				http.Error(w, "No content", http.StatusNoContent)
			case codes.Internal:
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
			return
		}
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		logger.Log.Info("JSON encode error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}

// PostUserData runs batch update requests of user data.
func (h *HTTP) PostUserData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Info("internal error", zap.Error(err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		logger.Log.Info("internal error", zap.Error(err))
		http.Error(w, "Empty batch", http.StatusBadRequest)
		return
	}
	var req []models.UserData
	if err := json.Unmarshal(body, &req); err != nil {
		logger.Log.Info("internal error", zap.Error(err))
		http.Error(w, "JSON decode error", http.StatusBadRequest)
		return
	}
	h.s.MergeUserData(req)
	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
}
