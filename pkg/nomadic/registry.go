package nomadic

import (
	"log"
)

var registry = map[string]Deployable{}

func Register(app Deployable){
	if app == nil {
		log.Fatal("Cannot register a nil application")
	}

	if app.Name() == "" {
		log.Fatal("Cannot register an application with no name")
	}

	registry[app.Name()] = app
}

func Registry() map[string]Deployable {
	return registry
}

func Find(name string) Deployable {
	if app, ok := registry[name]; ok {
		return app
	}
	return nil
}
