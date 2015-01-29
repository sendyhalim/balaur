package balaur

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/golang/glog"
	"github.com/zenazn/goji/web"
)

var neededConfigs = []string{"app", "route", "middleware"}

func NewApp(dir string, configs map[string]string, rr RouteRegistrar, mr MiddlewareRegistrar) *App {
	app := &App{
		appConfig:              map[string]Config{},
		Controllers:            map[string]interface{}{},
		Middlewares:            map[string]interface{}{},
		Mux:                    web.New(),
		AppRouteRegistrar:      rr,
		AppMiddlewareRegistrar: mr,
	}

	if rr == nil {
		app.AppRouteRegistrar = NewBasicRouteRegistrar()
	}

	if mr == nil {
		app.AppMiddlewareRegistrar = NewBasicMiddlewareRegistrar()
	}

	// load all given configs
	checkNeededConfigs(configs)
	for k, p := range configs {
		app.appConfig[k] = NewConfig(fmt.Sprintf("%s/%s", dir, p))
	}

	app.Name = app.Config("app").Get("name", false)

	return app
}

func checkNeededConfigs(configs map[string]string) {
	for _, k := range neededConfigs {
		_, ok := configs[k]
		if !ok {
			glog.Fatalf("Please set config for %s", k)
		}
	}
}

type App struct {
	Name                    string
	Priority                int
	Controllers             map[string]interface{}
	Middlewares             map[string]interface{}
	ControllerMethodWrapper func(method ControllerMethod) web.HandlerFunc
	Mux                     *web.Mux
	AppRouteRegistrar       RouteRegistrar
	AppMiddlewareRegistrar  MiddlewareRegistrar

	appConfig map[string]Config
}

func (a *App) Boot() {
	color.Green("Booting app %s(priority %d)\n", a.Name, a.Priority)
	a.registerRoutes()
	a.registerMiddlewares()
	AddApp(a)
}

func (a *App) Config(section string) Config {
	conf, ok := a.appConfig[section]

	if !ok {
		glog.Fatalf("Config (%s) is not set", section)
	}

	return conf
}

func (a *App) registerRoutes() {
	conf := a.Config("route")
	a.AppRouteRegistrar.RegisterRoutes(a.Mux, conf.GetChildren("route", false), a.Controllers)
}

func (a *App) registerMiddlewares() {
	conf := a.Config("middleware")
	a.AppMiddlewareRegistrar.RegisterMiddlewares(a.Mux, conf.GetChildren("middleware", false), a.Middlewares)
}
