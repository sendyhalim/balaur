package postapp

import (
	"net/http"

	"github.com/sendyhalim/balaur/context"
)

type Ctrl struct {
}

func (t Ctrl) Index(c context.ContextInterface, r *http.Request) (string, int) {
	return "index", http.StatusOK
}

func (t Ctrl) Get(c context.ContextInterface, r *http.Request) (string, int) {
	return "get " + c.GetParam("id"), http.StatusOK
}
