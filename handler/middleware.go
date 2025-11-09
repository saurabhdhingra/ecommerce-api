package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (s *jwtAuthService) Middleware(next http.HandlerFunc, requiredAdmin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		tokenString := parts[1]
		claims, err := s.ValidateToken(tokenString)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid or expired token: "+err.Error())
			return
		}

		if requiredAdmin && !claims.IsAdmin {
			respondError(w, http.StatusForbidden, "Access denied: Admin privilege required")
			return
		}

		// Attach the user ID and Admin status to the request context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next(w, r.WithContext(ctx))
	}
}