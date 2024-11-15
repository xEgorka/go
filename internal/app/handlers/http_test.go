package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"

	"github.com/xEgorka/project3/internal/app/config"
	"github.com/xEgorka/project3/internal/app/jobs"
	"github.com/xEgorka/project3/internal/app/mocks"
	"github.com/xEgorka/project3/internal/app/models"
	"github.com/xEgorka/project3/internal/app/service"
)

func TestHTTP_PostUserRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := service.New(&config.Config{}, ms, j, ma)
	h := NewHTTP(s)
	type want struct {
		contentType string
		code        int
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "positive test #1",
			body: `{"usr": "testLogin","pass": "testPass"}`,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #1",
			body: `{"usr": "testLogin","pass": ""}`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #2",
			body: `bad json`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #3",
			body: `{"usr": "testLogin","pass": "testPass"}`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #4",
			body: `{"usr": "testLogin","pass": "testPass"}`,
			want: want{
				code:        http.StatusConflict,
				contentType: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost,
				"/api/user/register", strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			ctx := context.Background()
			if tt.want.code == http.StatusOK {
				ma.EXPECT().Register(ctx, "testLogin", "testPass").Return(nil)
			}
			if tt.name == "negative test #3" {
				ma.EXPECT().Register(ctx, "testLogin", "testPass").
					Return(errors.New("test"))
			}
			if tt.name == "negative test #4" {
				ma.EXPECT().Register(ctx, "testLogin", "testPass").Return(
					&pgconn.PgError{Code: pgerrcode.UniqueViolation})
			}
			h.PostUserRegister(w, r.WithContext(ctx))
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)
			if err := res.Body.Close(); err != nil {
				panic(err)
			}
		})
	}
}

func TestHTTP_PostUserLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := service.New(&config.Config{}, ms, j, ma)
	h := NewHTTP(s)
	type want struct {
		contentType string
		code        int
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "positive test #1",
			body: `{"usr": "testLogin","pass": "testPass"}`,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #2",
			body: `bad json`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #3",
			body: `{"usr": "testLogin","pass": ""}`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #4",
			body: `{"usr": "testLogin","pass": "testPass"}`,
			want: want{
				code:        http.StatusUnauthorized,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost,
				"/api/user/login", strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			ctx := context.Background()
			if tt.want.code == http.StatusOK {
				ma.EXPECT().Login(ctx, "testLogin", "testPass").Return(nil)
			}
			if tt.name == "negative test #4" {
				ma.EXPECT().Login(ctx, "testLogin", "testPass").Return(errors.New("test"))
			}
			h.PostUserLogin(w, r.WithContext(ctx))
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)
			if err := res.Body.Close(); err != nil {
				panic(err)
			}
		})
	}
}

func TestHandlers_GetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := service.New(&config.Config{}, ms, j, ma)
	h := NewHTTP(s)
	type want struct {
		contentType string
		url         string
		code        int
	}
	tests := []struct {
		name      string
		body      string
		userID    string
		want      want
		timestamp string
	}{
		{
			name:      "negative test #1",
			timestamp: "0",
			body:      "",
			userID:    "x0o1@ya.ru",
			want: want{
				code:        http.StatusNoContent,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:      "negative test #2",
			timestamp: "0",
			body:      "",
			userID:    "x0o1@ya.ru",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:      "positive test #1",
			timestamp: "0",
			body:      "",
			userID:    "x0o1@ya.ru",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/"+tt.timestamp,
				strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()
			var userKey userIDKeyType
			ctx := context.WithValue(r.Context(), userKey, tt.userID)
			if tt.want.code == http.StatusNoContent {
				timestamp, _ := strconv.Atoi(tt.timestamp)
				ms.EXPECT().GetUserData(ctx, tt.userID, timestamp).Return(nil, nil)
			}
			if tt.name == "negative test #2" {
				timestamp, _ := strconv.Atoi(tt.timestamp)
				ms.EXPECT().GetUserData(ctx, tt.userID, timestamp).Return(nil, errors.New("test"))
			}
			if tt.want.code == http.StatusOK {
				timestamp, _ := strconv.Atoi(tt.timestamp)
				ms.EXPECT().GetUserData(ctx, tt.userID, timestamp).Return(
					[]models.UserData{{ID: "testID"}}, nil)
			}
			h.GetUserData(w, r.WithContext(ctx))
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)
			if err := res.Body.Close(); err != nil {
				panic(err)
			}
		})
	}
}

func TestHTTP_PostUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := service.New(&config.Config{}, ms, j, ma)
	h := NewHTTP(s)
	type want struct {
		contentType string
		code        int
	}

	tests := []struct {
		name   string
		body   string
		userID string
		want   want
	}{
		{
			name:   "positive test #1",
			body:   `[{"ID": "secret pass"}]`,
			userID: "x0o1@ya.ru",
			want: want{
				code:        http.StatusAccepted,
				contentType: "text/plain",
			},
		},
		{
			name:   "negative test #1",
			body:   ``,
			userID: "x0o1@ya.ru",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #2",
			body:   `bad`,
			userID: "x0o1@ya.ru",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/user/data",
				strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			var userKey userIDKeyType
			ctx := context.WithValue(r.Context(), userKey, tt.userID)
			h.PostUserData(w, r.WithContext(ctx))
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)
			if err := res.Body.Close(); err != nil {
				panic(err)
			}
		})
	}
}

func TestHandlers_GetPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ma := mocks.NewMockAuth(ctrl)
	ms := mocks.NewMockStorage(ctrl)
	j := jobs.New(ms)
	s := service.New(&config.Config{}, ms, j, ma)
	h := NewHTTP(s)
	type want struct {
		contentType string
		url         string
		code        int
	}
	tests := []struct {
		name      string
		body      string
		userID    string
		want      want
		timestamp int
	}{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #1",
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/ping",
				strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()
			ctx := context.Background()
			if tt.want.code == http.StatusOK {
				ms.EXPECT().Ping().Return(nil)
			} else {
				ms.EXPECT().Ping().Return(errors.New("test"))
			}
			h.GetPing(w, r.WithContext(ctx))
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)
			if err := res.Body.Close(); err != nil {
				panic(err)
			}
		})
	}
}
