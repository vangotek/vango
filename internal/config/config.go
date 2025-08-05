package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
)

// Config represents the site configuration
type Config struct {
	Title       string            `toml:"title"`
	BaseURL     string            `toml:"baseURL"`
	Language    string            `toml:"language"`
	Description string            `toml:"description"`
	Author      string            `toml:"author"`
	Theme       string            `toml:"theme"`
	Params      map[string]interface{} `toml:"params"`
	
	// Directory paths
	ContentDir string `toml:"contentDir"`
	LayoutDir  string `toml:"layoutDir"`
	StaticDir  string `toml:"staticDir"`
	PublicDir  string `toml:"publicDir"`
	
	// Build settings
	BuildDrafts bool   `toml:"buildDrafts"`
	BuildFuture bool   `toml:"buildFuture"`
	CleanBuild  bool   `toml:"cleanBuild"`
	
	// Server settings
	Port       int    `toml:"port"`
	Host       string `toml:"host"`
	LiveReload bool   `toml:"liveReload"`
}

// Load reads and parses the configuration file
func Load(path string) (*Config, error) {
	// Set defaults
	cfg := &Config{
		Title:       "VanGo Site",
		BaseURL:     "http://localhost:1313/",
		Language:    "en",
		Description: "A static site built with VanGo",
		ContentDir:  "content",
		LayoutDir:   "layouts",
		StaticDir:   "static",
		PublicDir:   "public",
		BuildDrafts: false,
		BuildFuture: false,
		CleanBuild:  true,
		Port:        1313,
		Host:        "localhost",
		LiveReload:  true,
		Params:      make(map[string]interface{}),
	}

	// Check if config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil // Return defaults if no config file
	}

	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse TOML
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	if c.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	
	if c.BaseURL == "" {
		return fmt.Errorf("baseURL cannot be empty")
	}

	// Ensure directories exist or can be created
	dirs := []string{c.ContentDir, c.LayoutDir}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("required directory does not exist: %s", dir)
		}
	}

	return nil
}

// GetParam returns a parameter value by key
func (c *Config) GetParam(key string) interface{} {
	return c.Params[key]
}

// SetParam sets a parameter value
func (c *Config) SetParam(key string, value interface{}) {
	c.Params[key] = value
}
