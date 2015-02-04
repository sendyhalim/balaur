package context

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type GojiContext struct {
	*web.C
}

func (q *GojiContext) GetParam(key string) string {
	return q.C.URLParams[key]
}

func (q *GojiContext) Get(r *http.Request, key interface{}) interface{} {
	return q.C.Env[key]
}

func (q *GojiContext) Set(r *http.Request, key interface{}, value interface{}) {
	q.C.Env[key] = value
}
