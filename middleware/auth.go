package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/hrusfandi/sb-task-management/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := bearerToken[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (*utils.JWTClaims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*utils.JWTClaims)
	return claims, ok
}