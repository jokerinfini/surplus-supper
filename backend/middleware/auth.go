package middleware

import (
	"context"
	"net/http"
	"strings"

	"surplus-supper/backend/userService"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	authService *userService.AuthService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		authService: userService.NewAuthService(),
	}
}

// Authenticate middleware function
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if it's a Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// Validate token
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		ctx = context.WithValue(ctx, "user_claims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth middleware that doesn't require authentication but adds user info if token is present
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
				tokenString := tokenParts[1]
				claims, err := m.authService.ValidateToken(tokenString)
				if err == nil {
					ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
					ctx = context.WithValue(ctx, "user_email", claims.Email)
					ctx = context.WithValue(ctx, "user_claims", claims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext gets user ID from request context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("user_id").(int)
	return userID, ok
}

// GetUserEmailFromContext gets user email from request context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	userEmail, ok := ctx.Value("user_email").(string)
	return userEmail, ok
}

// GetUserClaimsFromContext gets user claims from request context
func GetUserClaimsFromContext(ctx context.Context) (*userService.JWTClaims, bool) {
	claims, ok := ctx.Value("user_claims").(*userService.JWTClaims)
	return claims, ok
}
