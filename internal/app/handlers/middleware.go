package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/xEgorka/project3/internal/app/auth"
	"github.com/xEgorka/project3/internal/app/logger"
)

type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Writes respose using original http.ResponseWriter and grabs a size.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// Writes status code using original http.ResponseWriter and grabs a
// status code.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

var sugar = logger.Log.Sugar()

// WithLogging embeds response data to the original http.ResponseWriter.
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sugar = logger.Log.Sugar()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		next.ServeHTTP(&lw, r)
		duration := time.Since(start)
		sugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)
	})
}

type userIDKeyType struct{}

// Auth authorizes request.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string
		cr, err := r.Cookie("token")
		if err != nil { // no cookies
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString = cr.Value
		if len(tokenString) == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		login, err := auth.GetUserID(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var userKey userIDKeyType
		ctx := context.WithValue(r.Context(), userKey, login)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
