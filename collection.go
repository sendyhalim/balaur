package balaur

var collection = map[string]*App{}

func AddApp(app *App) {
	collection[app.Name] = app
}

func GetApp(name string) *App {
	return collection[name]
}

func RemoveApp(name string) {
	delete(collection, name)
}
