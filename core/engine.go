package core

import "go.uber.org/dig"

// Engine is the start point of an app
type Engine struct {
	Container      *dig.Container
	configFilePath string
}

// Fire initiates the core engine
func (e *Engine) Fire() *dig.Container {
	container := initContainer(e.configFilePath)
	return container
}

// NewEngine creates a new engine instance
func NewEngine() *Engine {
	engine := Engine{}
	engine.Container = engine.Fire()
	return &engine
}

func initContainer(configFilePath string) *dig.Container {
	container := dig.New()
	container.Provide(NewConfig(configFilePath))
	container.Provide(NewRepository)
	container.Provide(NewCacheManager)
	container.Provide(NewRedisCache)
	container.Provide(NewApp)
	container.Provide(NewEngine)
	return container
}
