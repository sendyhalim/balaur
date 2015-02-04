package balaur

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sendyhalim/balaur/testfixtures"
	"github.com/sendyhalim/balaur/testfixtures/postapp"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/zenazn/goji/web"
)

func TestApp(t *testing.T) {
	Convey("Test App", t, func() {
		ctrl := postapp.Ctrl{}
		app := NewApp(testfixtures.TestAppPath, testfixtures.TestAppConfigPath, nil, nil)

		app.Controllers["test-controller"] = ctrl
		app.Middlewares["test-middleware"] = postapp.TestMiddleware

		app.registerRoutes()
		app.registerMiddlewares()

		Convey("Test registerRoutes()", func() {
			w := testServe(app.Mux, "GET", "/test/path")
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldEqual, "index")

			w = testServe(app.Mux, "GET", "/test/path/18")
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldEqual, "get 18")
		})

		Convey("Test registerMiddlewares", func() {
			w := testServe(app.Mux, "GET", "test/path")
			So(w.Header().Get("test-header"), ShouldEqual, "test/value")
		})
	})
}

func newContext() web.C {
	return web.C{
		Env:       map[interface{}]interface{}{},
		URLParams: map[string]string{},
	}
}

func testServe(handler web.Handler, verb string, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	request, _ := http.NewRequest(verb, path, nil)
	handler.ServeHTTPC(newContext(), w, request)

	return w
}
