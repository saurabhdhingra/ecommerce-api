package handler

import (
	"context"
	"net/http"
	"strings"

	"ecommerce-api/service"
)

// AuthMiddleware wraps a JWT service middleware for use in handlers
func AuthMiddleware(jwtService service.JWTService, next http.HandlerFunc, requiredAdmin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			RespondError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			RespondError(w, http.StatusUnauthorized, "Invalid or expired token: "+err.Error())
			return
		}

		if requiredAdmin && !claims.IsAdmin {
			RespondError(w, http.StatusForbidden, "Access denied: Admin privilege required")
			return
		}

		// Attach the user ID and Admin status to the request context
		ctx := context.WithValue(r.Context(), service.UserContextKey, claims)
		next(w, r.WithContext(ctx))
	}
}
