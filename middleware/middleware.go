package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/cap/jwt"
	"github.com/heikkilamarko/goutils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// NotFoundHandler func
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	goutils.WriteNotFound(w, nil)
}

// Logger middleware
func Logger(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return hlog.NewHandler(*logger)
}

// RequestLogger middleware
func RequestLogger() func(next http.Handler) http.Handler {
	return hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Send()
	})
}

// ErrorRecovery middleware
func ErrorRecovery() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					hlog.FromRequest(r).Error().Msgf("%s", err)
					goutils.WriteInternalError(w, nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// APIKey middleware
func APIKey(apiKey, headerKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ak := r.Header.Get(headerKey); ak != apiKey {
				goutils.WriteUnauthorized(w, nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Timeout middleware
func Timeout(duration time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// JWTConfig struct
type JWTConfig struct {
	Issuer   string
	Iss      string
	Aud      []string
	TokenKey interface{}
}

// JWT middleware
func JWT(ctx context.Context, config *JWTConfig) func(next http.Handler) http.Handler {

	keySet, err := jwt.NewOIDCDiscoveryKeySet(ctx, config.Issuer, "")
	if err != nil {
		panic(err)
	}

	validator, err := jwt.NewValidator(keySet)
	if err != nil {
		panic(err)
	}

	expected := jwt.Expected{
		Issuer:    config.Iss,
		Audiences: config.Aud,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			logger := hlog.FromRequest(r)

			token, err := extractToken(r)
			if err != nil {
				logger.Error().Err(err).Send()
				goutils.WriteUnauthorized(w, nil)
				return
			}

			claims, err := validator.Validate(r.Context(), token, expected)
			if err != nil {
				logger.Error().Err(err).Send()
				goutils.WriteUnauthorized(w, nil)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), config.TokenKey, claims))

			next.ServeHTTP(w, r)
		})
	}
}

func extractToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	parts := strings.Split(auth, " ")
	if len(parts) == 2 {
		return parts[1], nil
	}
	return "", errors.New("token not found")
}
