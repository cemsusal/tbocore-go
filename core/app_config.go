package core

// Config configures the application
type Config struct {
	ConnectionString string
	DbPort           string
	RedisURI         string
	RedisPass        string
	RedisDb          int
}

// NewConfig initiates new config instance
func NewConfig(connectionString string, port string) *Config {
	return &Config{ConnectionString: connectionString, DbPort: port}
}
