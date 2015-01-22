package postapp

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

func TestMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("test-header", "test/value")
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
