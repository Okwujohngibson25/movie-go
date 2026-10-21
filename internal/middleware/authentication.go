package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	Email   string    `json:"email"`
	User_ID uuid.UUID `json:"userid"`
	Role    string
	Type    string `json:"typ"`
	jwt.RegisteredClaims
}

type contextKey string

const claimsKey contextKey = "claims"

func RequireAuth(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get token
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				apperrors.HandleError(w, apperrors.ErrTokenMissing)
				return
			}

			token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
				// Ensure the signing method is HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(secretKey), nil
			})
			if err != nil {
				apperrors.HandleError(w, apperrors.ErrInvalidToken)
				return
			}

			claims, ok := token.Claims.(*CustomClaims)
			if !ok || !token.Valid {
				apperrors.HandleError(w, apperrors.ErrInvalidToken)
				return
			}

			// refresh tokens (and legacy tokens without a type) must not grant access
			if claims.Type != "access" {
				apperrors.HandleError(w, apperrors.ErrInvalidToken)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				claimsKey,
				claims,
			)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
