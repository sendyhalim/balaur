package balaur

var appCollection = map[string]*App{}

func AddApp(app *App) {
	appCollection[app.Name] = app
}

func GetApp(name string) *App {
	return appCollection[name]
}

func RemoveApp(name string) {
	delete(appCollection, name)
}

func TotalApps() int {
	return len(appCollection)
}
