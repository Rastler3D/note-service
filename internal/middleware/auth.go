package middleware

import (
	"context"
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		userID, ok := validateToken(token)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func validateToken(token string) (int64, bool) {
	tokenToUserID := map[string]int64{
		"token1": 1,
		"token2": 2,
	}

	userID, ok := tokenToUserID[token]
	return userID, ok
}
