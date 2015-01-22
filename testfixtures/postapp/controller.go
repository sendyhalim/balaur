package postapp

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type Ctrl struct {
}

func (t Ctrl) Index(c web.C, r *http.Request) (string, int) {
	return "index", http.StatusOK
}

func (t Ctrl) Get(c web.C, r *http.Request) (string, int) {
	return "get " + c.URLParams["id"], http.StatusOK
}
