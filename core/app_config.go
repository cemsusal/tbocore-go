package core

// Config configures the application
type Config interface {
	ConnectionString() string
	Port() string
	NewConfig() Config
}
