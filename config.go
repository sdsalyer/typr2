package main

// Application configuration
type Config struct {
	CommandKey string // Key to enter command mode (default ":")
	SearchKey  string // Key to enter search mode (default "/")
}

// Default configuration
func DefaultConfig() Config {
	return Config{
		CommandKey: ":",
		SearchKey:  "/",
	}
}

// TODO: Add configuration loading/saving methods here
// func LoadConfig(filename string) (Config, error) { ... }
// func (c Config) Save(filename string) error { ... }
// func (c Config) Validate() error { ... }
