package context

import (
	"net/http"
)

type ContextInterface interface {
	GetParam(string) string
	Get(*http.Request, string) interface{}
	Set(*http.Request, string, interface{})
}
