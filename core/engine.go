package core

import "go.uber.org/dig"

type Engine struct {
	Container *dig.Container
}

// Fire initiates the core engine
func (e *Engine) Fire() *dig.Container {
	container := initContainer()
	return container
}

func initContainer() *dig.Container {
	container := dig.New()
	container.Provide(NewConfig)
	container.Provide(NewRepository)
	container.Provide(NewCacheManager)
	container.Provide(NewApp)
	container.Provide(NewRedisCache)
	return container
}

// NewEngine creates a new engine instance
func NewEngine() *Engine {
	engine := Engine{}
	engine.Container = engine.Fire()
	return &engine
}
