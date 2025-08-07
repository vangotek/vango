package theme

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"vango/internal/config"
)

// Theme represents a VanGo theme
type Theme struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Author      string                 `json:"author"`
	Homepage    string                 `json:"homepage"`
	License     string                 `json:"license"`
	MinVersion  string                 `json:"min_vango_version"`
	Tags        []string               `json:"tags"`
	Features    []string               `json:"features"`
	Config      map[string]interface{} `json:"config"`
	Path        string                 `json:"-"`
	// Template and asset paths
	LayoutsDir string `json:"layouts_dir"`
	StaticDir  string `json:"static_dir"`
	AssetsDir  string `json:"assets_dir"`
	Templates map[string]string `json:"templates"`
	
	CSS       string `json:"css"`
}

// ThemeManager handles theme operations
type ThemeManager struct {
	config      *config.Config
	activeTheme *Theme
	themes      map[string]*Theme
	themesDir   string
	defaultTheme string
}

// ThemeConfig represents theme-specific configuration
type ThemeConfig struct {
	Colors      ColorScheme            `json:"colors"`
	Typography  Typography             `json:"typography"`
	Layout      LayoutConfig           `json:"layout"`
	Features    FeatureConfig          `json:"features"`
	CustomCSS   string                 `json:"custom_css"`
	CustomJS    string                 `json:"custom_js"`
	Params      map[string]interface{} `json:"params"`
}

// ColorScheme defines theme colors
type ColorScheme struct {
	Primary     string `json:"primary"`
	Secondary   string `json:"secondary"`
	Accent      string `json:"accent"`
	Background  string `json:"background"`
	Surface     string `json:"surface"`
	Text        string `json:"text"`
	TextMuted   string `json:"text_muted"`
	Border      string `json:"border"`
	Success     string `json:"success"`
	Warning     string `json:"warning"`
	Error       string `json:"error"`
	Info        string `json:"info"`
}

// Typography defines font settings
type Typography struct {
	FontFamily    string  `json:"font_family"`
	FontSize      string  `json:"font_size"`
	LineHeight    float64 `json:"line_height"`
	HeadingFont   string  `json:"heading_font"`
	MonospaceFont string  `json:"monospace_font"`
}

// LayoutConfig defines layout settings
type LayoutConfig struct {
	MaxWidth    string `json:"max_width"`
	Sidebar     bool   `json:"sidebar"`
	Navigation  string `json:"navigation"` // top, side, both
	Footer      bool   `json:"footer"`
	Comments    bool   `json:"comments"`
	SearchBox   bool   `json:"search_box"`
}

// FeatureConfig defines enabled features
type FeatureConfig struct {
	DarkMode        bool `json:"dark_mode"`
	Syntax          bool `json:"syntax_highlighting"`
	MathJax         bool `json:"mathjax"`
	TableOfContents bool `json:"table_of_contents"`
	ShareButtons    bool `json:"share_buttons"`
	ReadingTime     bool `json:"reading_time"`
	RelatedPosts    bool `json:"related_posts"`
	Analytics       bool `json:"analytics"`
}

// NewThemeManager creates a new theme manager
func NewThemeManager(cfg *config.Config) *ThemeManager {
	themesDir := "themes"
	if cfg.GetParam("themes_dir") != nil {
		if dir, ok := cfg.GetParam("themes_dir").(string); ok {
			themesDir = dir
		}
	}
	return &ThemeManager{
		config:    cfg,
		themes:    make(map[string]*Theme),
		themesDir: themesDir,
	}
}

// LoadThemes discovers and loads all available themes
func (tm *ThemeManager) LoadThemes() error {
	// Ensure themes directory exists
	if _, err := os.Stat(tm.themesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tm.themesDir, 0755); err != nil {
			return fmt.Errorf("failed to create themes directory: %w", err)
		}
	}
	// Walk through themes directory
	return filepath.Walk(tm.themesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() || path == tm.themesDir {
			return nil
		}
		// Check if this directory contains a theme
		themeFile := filepath.Join(path, "theme.json")
		if _, err := os.Stat(themeFile); os.IsNotExist(err) {
			return nil // Not a theme directory
		}
		// Load the theme
		theme, err := tm.loadTheme(path)
		if err != nil {
			fmt.Printf("Warning: failed to load theme from %s: %v\n", path, err)
			return nil // Continue loading other themes
		}
		tm.themes[theme.Name] = theme
		return nil
	})
}

// loadTheme loads a single theme from a directory
func (tm *ThemeManager) loadTheme(themePath string) (*Theme, error) {
	themeFile := filepath.Join(themePath, "theme.json")
	data, err := os.ReadFile(themeFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme.json: %w", err)
	}
	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("failed to parse theme.json: %w", err)
	}
	theme.Path = themePath
	// Set default paths if not specified
	if theme.LayoutsDir == "" {
		theme.LayoutsDir = "layouts"
	}
	if theme.StaticDir == "" {
		theme.StaticDir = "static"
	}
	if theme.AssetsDir == "" {
		theme.AssetsDir = "assets"
	}
	// Validate theme structure
	if err := tm.validateTheme(&theme); err != nil {
		return nil, fmt.Errorf("invalid theme structure: %w", err)
	}
	return &theme, nil
}

// InstallTheme installs a theme from a remote source or local path
func (tm *ThemeManager) InstallTheme(source string) error {
    if source == "" {
        return fmt.Errorf("theme source cannot be empty")
    }
    
    // For now, return an error indicating the feature is not implemented
    // You can expand this later to handle actual theme installation
    return fmt.Errorf("theme installation from remote sources is not yet implemented. Use 'vango theme create <name>' to create a new theme or manually copy themes to the themes/ directory")
}

// Alternative: If you want a basic implementation that copies from a local directory
func (tm *ThemeManager) InstallThemeFromPath(sourcePath, themeName string) error {
    if sourcePath == "" || themeName == "" {
        return fmt.Errorf("source path and theme name cannot be empty")
    }
    
    // Check if source exists
    if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
        return fmt.Errorf("source path does not exist: %s", sourcePath)
    }
    
    // Check if theme already exists
    themePath := filepath.Join(tm.themesDir, themeName)
    if _, err := os.Stat(themePath); !os.IsNotExist(err) {
        return fmt.Errorf("theme already exists: %s", themeName)
    }
    
    // Copy the theme
    return tm.copyDir(sourcePath, themePath)
}

// validateTheme checks if a theme has the required structure
func (tm *ThemeManager) validateTheme(theme *Theme) error {
    // Check required templates
    requiredTemplates := []string{
        "layouts/_default/single.html",
        "layouts/_default/list.html",
    }
    
    for _, template := range requiredTemplates {
        templatePath := filepath.Join(theme.Path, template)
        if _, err := os.Stat(templatePath); os.IsNotExist(err) {
            return fmt.Errorf("required template missing: %s", template)
        }
    }
    
    // Validate theme.json structure
    if theme.Name == "" {
        return fmt.Errorf("theme name cannot be empty")
    }
    
    return nil
}

// SetActiveTheme sets the currently active theme
func (tm *ThemeManager) SetActiveTheme(themeName string) error {
	theme, exists := tm.themes[themeName]
	if !exists {
		return fmt.Errorf("theme not found: %s", themeName)
	}
	tm.activeTheme = theme
	tm.config.Theme = themeName
	return nil
}

func (tm *ThemeManager) SetDefaultTheme(name string) {
	tm.defaultTheme = name
}


// GetActiveTheme returns the currently active theme
func (tm *ThemeManager) GetActiveTheme() *Theme {
	return tm.activeTheme
}

// GetTheme returns a theme by name
func (tm *ThemeManager) GetTheme(name string) (*Theme, bool) {
	theme, exists := tm.themes[name]
	return theme, exists
}

// ListThemes returns all available themes
func (tm *ThemeManager) ListThemes() map[string]*Theme {
	return tm.themes
}

// GetThemeTemplatesPath returns the templates path for the active theme
func (tm *ThemeManager) GetThemeTemplatesPath() string {
	if tm.activeTheme == nil {
		return tm.config.LayoutDir
	}
	return filepath.Join(tm.activeTheme.Path, tm.activeTheme.LayoutsDir)
}

// GetThemeStaticPath returns the static assets path for the active theme
func (tm *ThemeManager) GetThemeStaticPath() string {
	if tm.activeTheme == nil {
		return tm.config.StaticDir
	}
	return filepath.Join(tm.activeTheme.Path, tm.activeTheme.StaticDir)
}

// GetThemeAssetsPath returns the assets path for the active theme
func (tm *ThemeManager) GetThemeAssetsPath() string {
	if tm.activeTheme == nil {
		return ""
	}
	return filepath.Join(tm.activeTheme.Path, tm.activeTheme.AssetsDir)
}

// GetThemeConfig returns the theme configuration
func (tm *ThemeManager) GetThemeConfig() (*ThemeConfig, error) {
	if tm.activeTheme == nil {
		return tm.getDefaultThemeConfig(), nil
	}
	configPath := filepath.Join(tm.activeTheme.Path, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return tm.getDefaultThemeConfig(), nil
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme config: %w", err)
	}
	var config ThemeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse theme config: %w", err)
	}
	return &config, nil
}

// getDefaultThemeConfig returns the default theme configuration
func (tm *ThemeManager) getDefaultThemeConfig() *ThemeConfig {
	return &ThemeConfig{
		Colors: ColorScheme{
			Primary:    "#007bff",
			Secondary:  "#6c757d",
			Accent:     "#28a745",
			Background: "#ffffff",
			Surface:    "#f8f9fa",
			Text:       "#333333",
			TextMuted:  "#6c757d",
			Border:     "#e9ecef",
			Success:    "#28a745",
			Warning:    "#ffc107",
			Error:      "#dc3545",
			Info:       "#17a2b8",
		},
		Typography: Typography{
			FontFamily:    "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
			FontSize:      "16px",
			LineHeight:    1.6,
			HeadingFont:   "inherit",
			MonospaceFont: "'Courier New', monospace",
		},
		Layout: LayoutConfig{
			MaxWidth:   "800px",
			Sidebar:    false,
			Navigation: "top",
			Footer:     true,
			Comments:   false,
			SearchBox:  false,
		},
		Features: FeatureConfig{
			DarkMode:        true,
			Syntax:          true,
			MathJax:         false,
			TableOfContents: true,
			ShareButtons:    false,
			ReadingTime:     true,
			RelatedPosts:    false,
			Analytics:       false,
		},
		Params: make(map[string]interface{}),
	}
}

// CopyThemeAssets copies theme assets to the public directory
func (tm *ThemeManager) CopyThemeAssets(publicDir string) error {
	if tm.activeTheme == nil {
		return nil
	}
	staticPath := tm.GetThemeStaticPath()
	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		return nil // No static assets to copy
	}
	destPath := filepath.Join(publicDir, "theme")
	return tm.copyDir(staticPath, destPath)
}

// copyDir recursively copies a directory
func (tm *ThemeManager) copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		// Create destination directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()
		_, err = io.Copy(destFile, srcFile)
		return err
	})
}

// GenerateThemeCSS generates CSS variables from theme configuration
func (tm *ThemeManager) GenerateThemeCSS() (string, error) {
	config, err := tm.GetThemeConfig()
	if err != nil {
		return "", err
	}
	var css strings.Builder
	css.WriteString(":root {\n")
	// Color variables
	css.WriteString(fmt.Sprintf("  --color-primary: %s;\n", config.Colors.Primary))
	css.WriteString(fmt.Sprintf("  --color-secondary: %s;\n", config.Colors.Secondary))
	css.WriteString(fmt.Sprintf("  --color-accent: %s;\n", config.Colors.Accent))
	css.WriteString(fmt.Sprintf("  --color-background: %s;\n", config.Colors.Background))
	css.WriteString(fmt.Sprintf("  --color-surface: %s;\n", config.Colors.Surface))
	css.WriteString(fmt.Sprintf("  --color-text: %s;\n", config.Colors.Text))
	css.WriteString(fmt.Sprintf("  --color-text-muted: %s;\n", config.Colors.TextMuted))
	css.WriteString(fmt.Sprintf("  --color-border: %s;\n", config.Colors.Border))
	css.WriteString(fmt.Sprintf("  --color-success: %s;\n", config.Colors.Success))
	css.WriteString(fmt.Sprintf("  --color-warning: %s;\n", config.Colors.Warning))
	css.WriteString(fmt.Sprintf("  --color-error: %s;\n", config.Colors.Error))
	css.WriteString(fmt.Sprintf("  --color-info: %s;\n", config.Colors.Info))
	// Typography variables
	css.WriteString(fmt.Sprintf("  --font-family: %s;\n", config.Typography.FontFamily))
	css.WriteString(fmt.Sprintf("  --font-size: %s;\n", config.Typography.FontSize))
	css.WriteString(fmt.Sprintf("  --line-height: %g;\n", config.Typography.LineHeight))
	css.WriteString(fmt.Sprintf("  --heading-font: %s;\n", config.Typography.HeadingFont))
	css.WriteString(fmt.Sprintf("  --monospace-font: %s;\n", config.Typography.MonospaceFont))
	// Layout variables
	css.WriteString(fmt.Sprintf("  --max-width: %s;\n", config.Layout.MaxWidth))
	css.WriteString("}\n")
	// Add custom CSS if provided
	if config.CustomCSS != "" {
		css.WriteString("\n")
		css.WriteString(config.CustomCSS)
	}
	return css.String(), nil
}

// CreateTheme creates a new theme from template
func (tm *ThemeManager) CreateTheme(name, template string) error {
	themePath := filepath.Join(tm.themesDir, name)
	// Check if theme already exists
	if _, err := os.Stat(themePath); !os.IsNotExist(err) {
		return fmt.Errorf("theme already exists: %s", name)
	}
	// Create theme directory
	if err := os.MkdirAll(themePath, 0755); err != nil {
		return fmt.Errorf("failed to create theme directory: %w", err)
	}
	// Create theme structure
	dirs := []string{
		filepath.Join(themePath, "layouts", "_default"),
		filepath.Join(themePath, "layouts", "partials"),
		filepath.Join(themePath, "static", "css"),
		filepath.Join(themePath, "static", "js"),
		filepath.Join(themePath, "assets"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	// Create theme.json
	theme := Theme{
		Name:        name,
		Version:     "1.0.0",
		Description: fmt.Sprintf("A VanGo theme called %s", name),
		Author:      "Unknown",
		License:     "MIT",
		MinVersion:  "1.0.0",
		Tags:        []string{"simple", "clean"},
		Features:    []string{"responsive", "dark-mode"},
		Config:      make(map[string]interface{}),
		LayoutsDir:  "layouts",
		StaticDir:   "static",
		AssetsDir:   "assets",
	}
	themeData, err := json.MarshalIndent(theme, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal theme data: %w", err)
	}
	if err := os.WriteFile(filepath.Join(themePath, "theme.json"), themeData, 0644); err != nil {
		return fmt.Errorf("failed to write theme.json: %w", err)
	}
	// Create basic templates based on template type
	if err := tm.createThemeTemplates(themePath, template); err != nil {
		return fmt.Errorf("failed to create theme templates: %w", err)
	}
	return nil
}

// createThemeTemplates creates basic templates for a new theme
func (tm *ThemeManager) createThemeTemplates(themePath, template string) error {
	var templates map[string]string
	var css string
	switch template {
	case "blog":
		templates = tm.getBlogTemplates()
		css = tm.getBlogCSS()
	case "portfolio":
		templates = tm.getPortfolioTemplates()
		css = tm.getPortfolioCSS()
	case "docs":
		templates = tm.getDocsTemplates()
		css = tm.getDocsCSS()
	default: // basic
		templates = tm.getBasicTemplates()
		css = tm.getBasicCSS()
	}
	// Write templates
	for path, content := range templates {
		fullPath := filepath.Join(themePath, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}
	// Write CSS
	cssPath := filepath.Join(themePath, "static", "css", "style.css")
	return os.WriteFile(cssPath, []byte(css), 0644)
}


