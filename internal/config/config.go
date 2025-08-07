package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

// Enhanced Config with advanced features
type Config struct {
	// Basic site information
	Title         string            `toml:"title" yaml:"title"`
	BaseURL       string            `toml:"baseURL" yaml:"baseURL"`
	Language      string            `toml:"language" yaml:"language"`
	Description   string            `toml:"description" yaml:"description"`
	Author        string            `toml:"author" yaml:"author"`
	Theme         string            `toml:"theme" yaml:"theme"`
	Params        map[string]interface{} `toml:"params" yaml:"params"`
	
	// Directory configuration
	ContentDir    string `toml:"contentDir" yaml:"contentDir"`
	LayoutDir     string `toml:"layoutDir" yaml:"layoutDir"`
	StaticDir     string `toml:"staticDir" yaml:"staticDir"`
	PublicDir     string `toml:"publicDir" yaml:"publicDir"`
	ThemesDir     string `toml:"themesDir" yaml:"themesDir"`
	DataDir       string `toml:"dataDir" yaml:"dataDir"`
	AssetsDir     string `toml:"assetsDir" yaml:"assetsDir"`
	
	// Build configuration
	BuildDrafts   bool     `toml:"buildDrafts" yaml:"buildDrafts"`
	BuildFuture   bool     `toml:"buildFuture" yaml:"buildFuture"`
	BuildExpired  bool     `toml:"buildExpired" yaml:"buildExpired"`
	CleanBuild    bool     `toml:"cleanBuild" yaml:"cleanBuild"`
	Watch         bool     `toml:"watch" yaml:"watch"`
	Workers       int      `toml:"workers" yaml:"workers"`
	
	// Server configuration
	Port          int      `toml:"port" yaml:"port"`
	Host          string   `toml:"host" yaml:"host"`
	LiveReload    bool     `toml:"liveReload" yaml:"liveReload"`
	DevMode       bool     `toml:"devMode" yaml:"devMode"`
	
	// Content processing
	DefaultContentType string   `toml:"defaultContentType" yaml:"defaultContentType"`
	DefaultLayout      string   `toml:"defaultLayout" yaml:"defaultLayout"`
	SummaryLength      int      `toml:"summaryLength" yaml:"summaryLength"`
	
	// URL configuration
	PrettyURLs        bool              `toml:"prettyURLs" yaml:"prettyURLs"`
	CanonicalifyURLs  bool              `toml:"canonicalifyURLs" yaml:"canonicalifyURLs"`
	RelativeURLs      bool              `toml:"relativeURLs" yaml:"relativeURLs"`
	UglyURLs          bool              `toml:"uglyURLs" yaml:"uglyURLs"`
	
	// Markup configuration
	Markup            MarkupConfig      `toml:"markup" yaml:"markup"`
	
	// Multilingual support
	Languages         map[string]Language `toml:"languages" yaml:"languages"`
	DefaultContentLanguage string         `toml:"defaultContentLanguage" yaml:"defaultContentLanguage"`
	
	// SEO and social
	SEO               SEOConfig         `toml:"seo" yaml:"seo"`
	Social            SocialConfig      `toml:"social" yaml:"social"`
	
	// Performance and optimization
	Performance       PerformanceConfig `toml:"performance" yaml:"performance"`
	
	// Security
	Security          SecurityConfig    `toml:"security" yaml:"security"`
	
	// Plugin system
	Plugins           []PluginConfig    `toml:"plugins" yaml:"plugins"`
	
	// Advanced features
	Features          FeatureFlags      `toml:"features" yaml:"features"`
	
	// Environment-specific overrides
	Environment       string            `toml:"environment" yaml:"environment"`
	Environments      map[string]EnvConfig `toml:"environments" yaml:"environments"`
}
// MarkupConfig configures markdown processing
type MarkupConfig struct {
	Goldmark          GoldmarkConfig    `toml:"goldmark" yaml:"goldmark"`
	TableOfContents   TOCConfig         `toml:"tableOfContents" yaml:"tableOfContents"`
	Highlight         HighlightConfig   `toml:"highlight" yaml:"highlight"`
}

// GoldmarkConfig configures the Goldmark markdown processor
type GoldmarkConfig struct {
	Renderer          RendererConfig    `toml:"renderer" yaml:"renderer"`
	Parser            ParserConfig      `toml:"parser" yaml:"parser"`
	Extensions        ExtensionsConfig  `toml:"extensions" yaml:"extensions"`
}

type RendererConfig struct {
	Unsafe            bool `toml:"unsafe" yaml:"unsafe"`
	HardWraps         bool `toml:"hardWraps" yaml:"hardWraps"`
	XHTML             bool `toml:"xhtml" yaml:"xhtml"`
}

type ParserConfig struct {
	AutoHeadingID     bool `toml:"autoHeadingID" yaml:"autoHeadingID"`
	AutoHeadingIDType string `toml:"autoHeadingIDType" yaml:"autoHeadingIDType"`
	Attribute         bool `toml:"attribute" yaml:"attribute"`
}

type ExtensionsConfig struct {
	Table             bool `toml:"table" yaml:"table"`
	Strikethrough     bool `toml:"strikethrough" yaml:"strikethrough"`
	Linkify           bool `toml:"linkify" yaml:"linkify"`
	TaskList          bool `toml:"taskList" yaml:"taskList"`
	Footnote          bool `toml:"footnote" yaml:"footnote"`
	DefinitionList    bool `toml:"definitionList" yaml:"definitionList"`
	Typographer       bool `toml:"typographer" yaml:"typographer"`
}

// TOCConfig configures table of contents generation
type TOCConfig struct {
	StartLevel        int    `toml:"startLevel" yaml:"startLevel"`
	EndLevel          int    `toml:"endLevel" yaml:"endLevel"`
	Ordered           bool   `toml:"ordered" yaml:"ordered"`
}

// HighlightConfig configures syntax highlighting
type HighlightConfig struct {
	Style             string `toml:"style" yaml:"style"`
	LineNos           bool   `toml:"lineNos" yaml:"lineNos"`
	LineNumbersInTable bool  `toml:"lineNumbersInTable" yaml:"lineNumbersInTable"`
	TabWidth          int    `toml:"tabWidth" yaml:"tabWidth"`
	Guesslang         bool   `toml:"guesslang" yaml:"guesslang"`
}

// Language configuration for multilingual sites
type Language struct {
	LanguageName      string `toml:"languageName" yaml:"languageName"`
	ContentDir        string `toml:"contentDir" yaml:"contentDir"`
	Weight            int    `toml:"weight" yaml:"weight"`
	Title             string `toml:"title" yaml:"title"`
	Params            map[string]interface{} `toml:"params" yaml:"params"`
}


// SEOConfig configures SEO features
type SEOConfig struct {
	EnableRobotsTXT   bool     `toml:"enableRobotsTXT" yaml:"enableRobotsTXT"`
	EnableSitemap     bool     `toml:"enableSitemap" yaml:"enableSitemap"`
	EnableRSSFeed     bool     `toml:"enableRSSFeed" yaml:"enableRSSFeed"`
	EnableJSONFeed    bool     `toml:"enableJSONFeed" yaml:"enableJSONFeed"`
	SitemapFilename   string   `toml:"sitemapFilename" yaml:"sitemapFilename"`
	RSSFilename       string   `toml:"rssFilename" yaml:"rssFilename"`
	JSONFeedFilename  string   `toml:"jsonFeedFilename" yaml:"jsonFeedFilename"`
	MetaGenerator     bool     `toml:"metaGenerator" yaml:"metaGenerator"`
}

// SocialConfig configures social media integration
type SocialConfig struct {
	Twitter           string   `toml:"twitter" yaml:"twitter"`
	Facebook          string   `toml:"facebook" yaml:"facebook"`
	GitHub            string   `toml:"github" yaml:"github"`
	LinkedIn          string   `toml:"linkedin" yaml:"linkedin"`
	Instagram         string   `toml:"instagram" yaml:"instagram"`
	YouTube           string   `toml:"youtube" yaml:"youtube"`
	OpenGraph         OpenGraphConfig `toml:"openGraph" yaml:"openGraph"`
	TwitterCard       TwitterCardConfig `toml:"twitterCard" yaml:"twitterCard"`
}

type OpenGraphConfig struct {
	Enable            bool   `toml:"enable" yaml:"enable"`
	DefaultImage      string `toml:"defaultImage" yaml:"defaultImage"`
	SiteName          string `toml:"siteName" yaml:"siteName"`
}


type TwitterCardConfig struct {
	Enable            bool   `toml:"enable" yaml:"enable"`
	Site              string `toml:"site" yaml:"site"`
	Creator           string `toml:"creator" yaml:"creator"`
	DefaultImage      string `toml:"defaultImage" yaml:"defaultImage"`
}

// PerformanceConfig configures performance optimizations
type PerformanceConfig struct {
	EnableCompression bool     `toml:"enableCompression" yaml:"enableCompression"`
	EnableMinification bool    `toml:"enableMinification" yaml:"enableMinification"`
	EnableCaching     bool     `toml:"enableCaching" yaml:"enableCaching"`
	CacheDir          string   `toml:"cacheDir" yaml:"cacheDir"`
	ImageOptimization ImageOptConfig `toml:"imageOptimization" yaml:"imageOptimization"`
	AssetBundling     AssetBundlingConfig `toml:"assetBundling" yaml:"assetBundling"`
}

type ImageOptConfig struct {
	Enable            bool     `toml:"enable" yaml:"enable"`
	Quality           int      `toml:"quality" yaml:"quality"`
	Formats           []string `toml:"formats" yaml:"formats"`
	Responsive        bool     `toml:"responsive" yaml:"responsive"`
	LazyLoading       bool     `toml:"lazyLoading" yaml:"lazyLoading"`
}

type AssetBundlingConfig struct {
	Enable            bool     `toml:"enable" yaml:"enable"`
	CSS               bool     `toml:"css" yaml:"css"`
	JS                bool     `toml:"js" yaml:"js"`
	Fingerprinting    bool     `toml:"fingerprinting" yaml:"fingerprinting"`
}

// SecurityConfig configures security features
type SecurityConfig struct {
	ContentSecurityPolicy CSPConfig `toml:"contentSecurityPolicy" yaml:"contentSecurityPolicy"`
	HTTPS                HTTPSConfig `toml:"https" yaml:"https"`
	Headers              map[string]string `toml:"headers" yaml:"headers"`
}

type CSPConfig struct {
	Enable            bool   `toml:"enable" yaml:"enable"`
	DefaultSrc        string `toml:"defaultSrc" yaml:"defaultSrc"`
	ScriptSrc         string `toml:"scriptSrc" yaml:"scriptSrc"`
	StyleSrc          string `toml:"styleSrc" yaml:"styleSrc"`
	ImgSrc            string `toml:"imgSrc" yaml:"imgSrc"`
}

type HTTPSConfig struct {
	Enable            bool `toml:"enable" yaml:"enable"`
	RedirectHTTP      bool `toml:"redirectHTTP" yaml:"redirectHTTP"`
	HSTS              bool `toml:"hsts" yaml:"hsts"`
}

// PluginConfig configures individual plugins
type PluginConfig struct {
	Name              string                 `toml:"name" yaml:"name"`
	Version           string                 `toml:"version" yaml:"version"`
	Enabled           bool                   `toml:"enabled" yaml:"enabled"`
	Config            map[string]interface{} `toml:"config" yaml:"config"`
}

// FeatureFlags enables/disables experimental features
type FeatureFlags struct {
	ExperimentalMode  bool `toml:"experimentalMode" yaml:"experimentalMode"`
	BetaFeatures      bool `toml:"betaFeatures" yaml:"betaFeatures"`
	DebugMode         bool `toml:"debugMode" yaml:"debugMode"`
	ProfileMode       bool `toml:"profileMode" yaml:"profileMode"`
}

// EnvConfig allows environment-specific overrides
type EnvConfig struct {
	BaseURL           string                 `toml:"baseURL" yaml:"baseURL"`
	BuildDrafts       *bool                  `toml:"buildDrafts" yaml:"buildDrafts"`
	Minify            *bool                  `toml:"minify" yaml:"minify"`
	DevMode           *bool                  `toml:"devMode" yaml:"devMode"`
	Params            map[string]interface{} `toml:"params" yaml:"params"`
}

// ConfigLoader handles loading and validating configuration
type ConfigLoader struct {
	searchPaths []string
	envOverrides map[string]string
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		searchPaths: []string{
			"config.toml",
			"config.yaml",
			"config.yml",
			"vango.toml",
			"vango.yaml",
			"vango.yml",
			"config/config.toml",
			"config/config.yaml",
			"config/config.yml",
		},
		envOverrides: make(map[string]string),
	}
}

// AddSearchPath adds a path to search for configuration files
func (cl *ConfigLoader) AddSearchPath(path string) {
	cl.searchPaths = append(cl.searchPaths, path)
}

// SetEnvOverride sets an environment variable override
func (cl *ConfigLoader) SetEnvOverride(key, value string) {
	cl.envOverrides[key] = value
}


// Load reads and parses the configuration with enhanced features
func Load(configPath string) (*Config, error) {
	loader := NewConfigLoader()
	return loader.LoadConfig(configPath)
}

// LoadConfig loads configuration with full feature support
func (cl *ConfigLoader) LoadConfig(configPath string) (*Config, error) {
	// Set defaults
	cfg := cl.getDefaultConfig()

	// Determine config file to use
	var configFile string
	if configPath != "" {
		configFile = configPath
	} else {
		var err error
		configFile, err = cl.findConfigFile()
		if err != nil {
			return cfg, nil // Return defaults if no config file found
		}
	}

	// Load main config file
	if err := cl.loadConfigFile(configFile, cfg); err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configFile, err)
	}

	// Load environment-specific config
	if err := cl.loadEnvironmentConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to load environment config: %w", err)
	}

	// Apply environment variable overrides
	cl.applyEnvironmentOverrides(cfg)

	// Validate configuration
	if err := cl.validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Post-process configuration
	cl.postProcessConfig(cfg)

	return cfg, nil
}

// getDefaultConfig returns a configuration with sensible defaults
func (cl *ConfigLoader) getDefaultConfig() *Config {
	return &Config{
		Title:                  "VanGo Site",
		BaseURL:                "http://localhost:1313/",
		Language:               "en",
		Description:            "A static site built with VanGo",
		ContentDir:             "content",
		LayoutDir:              "layouts",
		StaticDir:              "static",
		PublicDir:              "public",
		ThemesDir:              "themes",
		DataDir:                "data",
		AssetsDir:              "assets",
		BuildDrafts:            false,
		BuildFuture:            false,
		BuildExpired:           false,
		CleanBuild:             true,
		Watch:                  false,
		Workers:                0, // Auto-detect
		Port:                   1313,
		Host:                   "localhost",
		LiveReload:             true,
		DevMode:                false,
		DefaultContentType:     "page",
		DefaultLayout:          "single",
		SummaryLength:          70,
		PrettyURLs:             true,
		CanonicalifyURLs:       false,
		RelativeURLs:           false,
		UglyURLs:               false,
		DefaultContentLanguage: "en",
		Environment:            "development",
		Params:                 make(map[string]interface{}),
		Languages:              make(map[string]Language),
		Environments:           make(map[string]EnvConfig),
		
		// Markup defaults
		Markup: MarkupConfig{
			Goldmark: GoldmarkConfig{
				Renderer: RendererConfig{
					Unsafe:    false,
					HardWraps: false,
					XHTML:     false,
				},
				Parser: ParserConfig{
					AutoHeadingID:     true,
					AutoHeadingIDType: "github",
					Attribute:         true,
				},
				Extensions: ExtensionsConfig{
					Table:          true,
					Strikethrough:  true,
					Linkify:        true,
					TaskList:       true,
					Footnote:       true,
					DefinitionList: true,
					Typographer:    true,
				},
			},
			TableOfContents: TOCConfig{
				StartLevel: 2,
				EndLevel:   6,
				Ordered:    false,
			},
			Highlight: HighlightConfig{
				Style:              "github",
				LineNos:            false,
				LineNumbersInTable: false,
				TabWidth:           4,
				Guesslang:          true,
			},
		},
		
		// SEO defaults
		SEO: SEOConfig{
			EnableRobotsTXT:  true,
			EnableSitemap:    true,
			EnableRSSFeed:    true,
			EnableJSONFeed:   false,
			SitemapFilename:  "sitemap.xml",
			RSSFilename:      "feed.xml",
			JSONFeedFilename: "feed.json",
			MetaGenerator:    true,
		},
		
		// Social defaults
		Social: SocialConfig{
			OpenGraph: OpenGraphConfig{
				Enable: true,
			},
			TwitterCard: TwitterCardConfig{
				Enable: true,
			},
		},
		
		// Performance defaults
		Performance: PerformanceConfig{
			EnableCompression:  true,
			EnableMinification: false,
			EnableCaching:      true,
			CacheDir:           ".cache",
			ImageOptimization: ImageOptConfig{
				Enable:      false,
				Quality:     85,
				Formats:     []string{"webp", "jpeg"},
				Responsive:  true,
				LazyLoading: true,
			},
			AssetBundling: AssetBundlingConfig{
				Enable:         false,
				CSS:            true,
				JS:             true,
				Fingerprinting: true,
			},
		},
		
		// Security defaults
		Security: SecurityConfig{
			ContentSecurityPolicy: CSPConfig{
				Enable: false,
			},
			HTTPS: HTTPSConfig{
				Enable:       false,
				RedirectHTTP: false,
				HSTS:         false,
			},
			Headers: make(map[string]string),
		},
		
		// Feature flags
		Features: FeatureFlags{
			ExperimentalMode: false,
			BetaFeatures:     false,
			DebugMode:        false,
			ProfileMode:      false,
		},
	}
}

// findConfigFile searches for a configuration file
func (cl *ConfigLoader) findConfigFile() (string, error) {
	for _, path := range cl.searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no configuration file found")
}

// loadConfigFile loads a specific configuration file
func (cl *ConfigLoader) loadConfigFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".toml":
		return toml.Unmarshal(data, cfg)
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, cfg)
	default:
		// Try to detect format
		if strings.Contains(string(data), "=") {
			return toml.Unmarshal(data, cfg)
		} else {
			return yaml.Unmarshal(data, cfg)
		}
	}
}

// loadEnvironmentConfig loads environment-specific configuration
func (cl *ConfigLoader) loadEnvironmentConfig(cfg *Config) error {
	if cfg.Environment == "" {
		return nil
	}

	envConfigPath := fmt.Sprintf("config/%s.toml", cfg.Environment)
	if _, err := os.Stat(envConfigPath); os.IsNotExist(err) {
		envConfigPath = fmt.Sprintf("config/%s.yaml", cfg.Environment)
		if _, err := os.Stat(envConfigPath); os.IsNotExist(err) {
			return nil // No environment-specific config
		}
	}

	envCfg := &Config{}
	if err := cl.loadConfigFile(envConfigPath, envCfg); err != nil {
		return err
	}

	// Merge environment config into main config
	cl.mergeConfigs(cfg, envCfg)

	return nil
}

// applyEnvironmentOverrides applies environment variable overrides
func (cl *ConfigLoader) applyEnvironmentOverrides(cfg *Config) {
	// Check for common environment variables
	envVars := map[string]func(string){
		"VANGO_BASE_URL": func(v string) { cfg.BaseURL = v },
		"VANGO_TITLE":    func(v string) { cfg.Title = v },
		"VANGO_THEME":    func(v string) { cfg.Theme = v },
		"VANGO_PORT":     func(v string) { 
			if port := parseInt(v); port > 0 { 
				cfg.Port = port 
			} 
		},
		"VANGO_HOST":     func(v string) { cfg.Host = v },
		"VANGO_ENV":      func(v string) { cfg.Environment = v },
	}

	for envVar, setter := range envVars {
		if value := os.Getenv(envVar); value != "" {
			setter(value)
		}
	}

	// Apply custom overrides
	for key, value := range cl.envOverrides {
		cl.setConfigValue(cfg, key, value)
	}
}

// validateConfig validates the configuration
func (cl *ConfigLoader) validateConfig(cfg *Config) error {
	if cfg.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	
	if cfg.BaseURL == "" {
		return fmt.Errorf("baseURL cannot be empty")
	}

	// Validate URLs
	if !cl.isValidURL(cfg.BaseURL) {
		return fmt.Errorf("invalid baseURL: %s", cfg.BaseURL)
	}

	// Ensure required directories exist
	requiredDirs := []string{cfg.ContentDir, cfg.LayoutDir}
	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("required directory does not exist: %s", dir)
		}
	}

	// Validate port range
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port: %d", cfg.Port)
	}

	// Validate markup configuration
	if err := cl.validateMarkupConfig(&cfg.Markup); err != nil {
		return fmt.Errorf("invalid markup config: %w", err)
	}

	return nil
}

// validateMarkupConfig validates markup-specific configuration
func (cl *ConfigLoader) validateMarkupConfig(markup *MarkupConfig) error {
	toc := &markup.TableOfContents
	if toc.StartLevel < 1 || toc.StartLevel > 6 {
		return fmt.Errorf("tableOfContents.startLevel must be between 1 and 6")
	}
	if toc.EndLevel < 1 || toc.EndLevel > 6 {
		return fmt.Errorf("tableOfContents.endLevel must be between 1 and 6")
	}
	if toc.StartLevel > toc.EndLevel {
		return fmt.Errorf("tableOfContents.startLevel cannot be greater than endLevel")
	}

	highlight := &markup.Highlight
	if highlight.TabWidth < 1 || highlight.TabWidth > 16 {
		return fmt.Errorf("highlight.tabWidth must be between 1 and 16")
	}

	return nil
}

// postProcessConfig performs post-processing on the configuration
func (cl *ConfigLoader) postProcessConfig(cfg *Config) {
	// Normalize URLs
	cfg.BaseURL = strings.TrimSuffix(cfg.BaseURL, "/") + "/"

	// Set worker count if not specified
	if cfg.Workers <= 0 {
		cfg.Workers = max(1, min(8, getNumCPU()))
	}

	// Ensure cache directory exists
	if cfg.Performance.EnableCaching && cfg.Performance.CacheDir != "" {
		os.MkdirAll(cfg.Performance.CacheDir, 0755)
	}

	// Set environment-specific defaults
	switch cfg.Environment {
	case "production":
		if cfg.Performance.EnableMinification == false {
			cfg.Performance.EnableMinification = true
		}
		cfg.DevMode = false
	case "development":
		cfg.DevMode = true
		cfg.Features.DebugMode = true
	}
}

// Helper methods
func (cl *ConfigLoader) mergeConfigs(base, override *Config) {
	// This would implement deep merging of configurations
	// For brevity, showing key fields only
	if override.BaseURL != "" {
		base.BaseURL = override.BaseURL
	}
	if override.Title != "" {
		base.Title = override.Title
	}
	// ... continue for all fields
}

func (cl *ConfigLoader) setConfigValue(cfg *Config, key, value string) {
	// Implement setting nested configuration values using dot notation
	// Example: "performance.enableMinification" = "true"
}

func (cl *ConfigLoader) isValidURL(urlStr string) bool {
	urlPattern := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlPattern.MatchString(urlStr) || strings.HasPrefix(urlStr, "http://localhost")
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

// Enhanced Config methods
func (c *Config) GetParam(key string) interface{} {
	return c.Params[key]
}

func (c *Config) SetParam(key string, value interface{}) {
	c.Params[key] = value
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) GetLanguage(code string) (Language, bool) {
	lang, exists := c.Languages[code]
	return lang, exists
}

func (c *Config) GetDefaultLanguage() Language {
	if lang, exists := c.Languages[c.DefaultContentLanguage]; exists {
		return lang
	}
	return Language{
		LanguageName: "English",
		ContentDir:   c.ContentDir,
		Weight:       0,
		Title:        c.Title,
		Params:       make(map[string]interface{}),
	}
}

// Utility functions
func parseInt(s string) int {
	// Simple integer parsing with error handling
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

func getNumCPU() int {
	// Would use runtime.NumCPU() in real implementation
	return 4
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
