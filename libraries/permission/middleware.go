package permission

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/reecerussell/monzo-plus-plus/libraries/util"

	"github.com/reecerussell/monzo-plus-plus/libraries/jwt"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
)

var (
	mu   = sync.RWMutex{}
	urls = []string{}

	errNoAuthHeader          = errors.Unauthorised("no authorization header")
	errMalformedAuthHeader   = errors.Unauthorised("malformed authorization header")
	errUnsupportedAuthScheme = errors.Unauthorised("unsupported authorization scheme")
)

// IgnoreURL adds a string to the ignore list.
func IgnoreURL(substr string) {
	mu.RLock()
	defer mu.RUnlock()

	urls = append(urls, substr)
}

// Middleware provides an authentication middleware method to ensure a user
// has the givenn permission.
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isIgnored(r.URL.Path) {
			h.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			errors.HandleHTTPError(w, r, errNoAuthHeader)
			return
		}

		p := strings.Split(auth, " ")
		if len(p) < 2 {
			errors.HandleHTTPError(w, r, errMalformedAuthHeader)
			return
		}

		if p[0] != "Bearer" {
			errors.HandleHTTPError(w, r, errUnsupportedAuthScheme)
			return
		}

		// populates request context with jwt claim values
		r = r.WithContext(populateContext(r.Context(), p[1]))

		h.ServeHTTP(w, r)
	})
}

func isIgnored(path string) bool {
	mu.Lock()
	defer mu.Unlock()

	for _, url := range urls {
		if strings.Contains(
			strings.ToLower(path),
			strings.ToLower(url)) {
			return true
		}
	}

	return false
}

func populateContext(ctx context.Context, token string) context.Context {
	t, tErr := jwt.FromToken([]byte(token))
	if tErr != nil {
		return ctx
	}

	ctx = context.WithValue(ctx, util.ContextKey("token"), token)

	if userID, ok := t.Claims.String(jwt.ClaimUserID); ok {
		ctx = context.WithValue(ctx, util.ContextKey("user_id"), userID)
	}

	return ctx
}
