package middleware

import (
	"context"
	"net/http"
	"strings"

	apierrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserEmailKey contextKey = "user_email"
const UserRoleKey contextKey = "user_role"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			apierrors.RespondWithError(w, http.StatusUnauthorized, apierrors.New(apierrors.ErrCodeUnauthorized, "Authorization header required"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			apierrors.RespondWithError(w, http.StatusUnauthorized, apierrors.New(apierrors.ErrCodeUnauthorized, "Invalid authorization header format"))
			return
		}

		token := parts[1]
		claims, err := utils.ValidateToken(token)
		if err != nil {
			apierrors.RespondWithError(w, http.StatusUnauthorized, apierrors.New(apierrors.ErrCodeUnauthorized, "Invalid or expired token"))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(UserRoleKey).(string)
		if !ok || role != "admin" {
			apierrors.RespondWithError(w, http.StatusForbidden, apierrors.New(apierrors.ErrCodeForbidden, "Admin access required"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) string {
	userID, _ := r.Context().Value(UserIDKey).(string)
	return userID
}

func GetUserRole(r *http.Request) string {
	role, _ := r.Context().Value(UserRoleKey).(string)
	return role
}
