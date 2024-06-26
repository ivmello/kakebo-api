package middlewares

import (
	"context"
	"net/http"

	"github.com/ivmello/kakebo-go-api/internal/provider"
	"github.com/ivmello/kakebo-go-api/internal/utils"
)

type authMiddleware struct {
	provider *provider.Provider
}

func NewAuthMiddleware(provider *provider.Provider) *authMiddleware {
	return &authMiddleware{
		provider,
	}
}

func (m *authMiddleware) Execute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tokenString = tokenString[len("Bearer "):]
		authService := m.provider.GetAuthService()
		err := authService.VerifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user, err := authService.GetUserFromToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), utils.USER_ID_KEY, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
