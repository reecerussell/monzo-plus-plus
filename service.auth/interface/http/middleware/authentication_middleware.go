package middleware

import (
	"net/http"
	"strings"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"

	"github.com/reecerussell/monzo-plus-plus/service.auth/usecase"
)

// AnonymousRoutes is an array of url paths, that can bypass
// the authentication middleware.
var AnonymousRoutes = [1]string{"/token"}

type AuthenticationMiddleware struct {
	uu usecase.UserUsecase
}

func NewAuthenticationMiddleware(ctn *di.Container) *AuthenticationMiddleware {
	uu := ctn.Resolve("user_usecase").(usecase.UserUsecase)

	return &AuthenticationMiddleware{uu}
}

func (am *AuthenticationMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, rp := range AnonymousRoutes {
			if strings.Contains(
				strings.ToLower(r.URL.Path),
				strings.ToLower(rp)) {
				h.ServeHTTP(w, r)
				return
			}
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("no authorization header"))
			return
		}

		p := strings.Split(auth, " ")
		if len(p) < 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("malformed authorization header"))
			return
		}

		if p[0] == "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unsupported authorization scheme"))
			return
		}

		ctx, err := am.uu.WithUser(r.Context(), p[1])
		if err != nil {
			w.WriteHeader(err.ErrorCode())
			w.Write([]byte(err.Text()))
			return
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
