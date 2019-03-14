package core

import "go.uber.org/dig"

// Fire initiates the core engine
func Fire() *dig.Container {
	container := initContainer()

	return container
}

func initContainer() *dig.Container {
	container := dig.New()
	container.Provide(Config.NewConfig)
	container.Provide(NewRepository)
	return container
}
