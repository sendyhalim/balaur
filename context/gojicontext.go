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

func (q *GojiContext) Get(r *http.Request, key string) interface{} {
	return q.C.Env[key]
}

func (q *GojiContext) Set(r *http.Request, key string, value interface{}) {
	q.C.Env[key] = value
}
