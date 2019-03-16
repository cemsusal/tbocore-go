package core

//TheApp is the starting point of an app
type TheApp struct {
	Cache  *CacheManager
	Db     *Repository
	Engine *Engine
}

// NewApp initiates the app
func NewApp(cache CacheManager, db Repository) *TheApp {
	engine := NewEngine()
	return &TheApp{Cache: &cache, Db: &db, Engine: engine}
}
