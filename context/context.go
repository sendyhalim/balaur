package context

import (
	"net/http"
)

type ContextInterface interface {
	GetParam(string) string
	Get(*http.Request, interface{}) interface{}
	Set(*http.Request, interface{}, interface{})
}
