package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

type AuthMiddleware struct {
	AuthClient *auth.Client
}

func NewAuthMiddleware(authClient *auth.Client) *AuthMiddleware {
	return &AuthMiddleware{AuthClient: authClient}
}

func (am *AuthMiddleware) FirebaseAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		idToken := tokenParts[1]

		token, err := am.AuthClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "firebaseUser", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetFirebaseUser(r *http.Request) (*auth.Token, error) {
	user, ok := r.Context().Value("firebaseUser").(*auth.Token)
	if !ok {
		return nil, errors.New("user not authenticated")
	}
	return user, nil
}
