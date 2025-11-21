package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ContextKey is the type for context keys
type ContextKey string

const UserContextKey ContextKey = "user"

// UserContext represents the authenticated user
type UserContext struct {
	UserID string
	Email  string
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens and adds user context
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			
			// If no auth header, continue without user context
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from "Bearer <token>"
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				// No Bearer prefix found
				next.ServeHTTP(w, r)
				return
			}

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				// Invalid token, continue without user context
				next.ServeHTTP(w, r)
				return
			}

			// Extract claims
			if claims, ok := token.Claims.(*JWTClaims); ok {
				// Add user to context
				userCtx := &UserContext{
					UserID: claims.UserID,
					Email:  claims.Email,
				}
				ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext retrieves the user from context
func GetUserFromContext(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value(UserContextKey).(*UserContext)
	return user, ok
}