package core

import (
	"go.uber.org/dig"
)

// Engine is the start point of an app
type Engine struct {
	Container      *dig.Container
	ConfigFilePath string
}

// Fire initiates the core engine
func (e *Engine) Fire() {
	container := initContainer(e.ConfigFilePath)
	e.Container = container
}

// NewEngine creates a new engine instance
func NewEngine() *Engine {
	engine := Engine{}
	engine.Fire()
	return &engine
}

func initContainer(configFilePath string) *dig.Container {
	container := dig.New()
	container.Provide(func() *Config {
		return NewConfig(configFilePath)
	})
	container.Provide(NewRepository)
	container.Provide(NewCacheManager)
	container.Provide(NewRedisCache)
	container.Provide(NewEngine)
	return container
}
