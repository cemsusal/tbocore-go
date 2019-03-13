package core

type theCache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) (bool, error)
	IsSet(key string) (bool, error)
	Remove(key string) (bool, error)
	RemoveByPattern(pattern string) (bool, error)
}
