package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/hashicorp/cap/jwt"
	"github.com/heikkilamarko/goutils"
)

func ErrorRecovery(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("error recovery", "panic", err)
					goutils.WriteInternalError(w, nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

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

type JWTConfig struct {
	Issuer   string
	Iss      string
	Aud      []string
	TokenKey any
	Logger   *slog.Logger
}

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
			token := goutils.TokenFromHeader(r)
			if token == "" {
				slog.Error("token is empty")
				goutils.WriteUnauthorized(w, nil)
				return
			}

			claims, err := validator.Validate(r.Context(), token, expected)
			if err != nil {
				slog.Error(err.Error())
				goutils.WriteUnauthorized(w, nil)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), config.TokenKey, claims))

			next.ServeHTTP(w, r)
		})
	}
}
