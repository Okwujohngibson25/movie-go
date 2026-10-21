package middleware

import (
	"net/http"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
)

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(claimsKey).(*CustomClaims)
		if !ok {
			apperrors.HandleError(w, apperrors.ErrInvalidToken)
			return
		}

		if claims.Role != "admin" {
			apperrors.HandleError(w, apperrors.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
