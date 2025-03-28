package database

import "time"

// Config represents the common database configuration
type Config struct {
	DSN                string        `yaml:"dsn"`
	DialTimeout        time.Duration `yaml:"dial_timeout"`
	ReadTimeout        time.Duration `yaml:"read_timeout"`
	WriteTimeout       time.Duration `yaml:"write_timeout"`
	MaxOpenConnections int           `yaml:"max_open_connections"`
	MaxIdleConnections int           `yaml:"max_idle_connections"`
	MaxLifeConnections time.Duration `yaml:"max_life_connections"`
	DebugSQL           bool          `yaml:"debug_sql"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		DialTimeout:        5 * time.Second,
		ReadTimeout:        50 * time.Second,
		WriteTimeout:       30 * time.Second,
		MaxOpenConnections: 256,
		MaxIdleConnections: 5,
		MaxLifeConnections: 30 * time.Second,
		DebugSQL:           true,
	}
}

// Database interface defines the common methods that all database implementations must provide
type Database interface {
	// Connect establishes a connection to the database
	Connect() error

	// Close closes the database connection
	Close() error

	// Ping checks if the database connection is alive
	Ping() error
}
