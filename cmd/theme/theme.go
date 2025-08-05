// cmd/theme/theme.go
package theme

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"vango/internal/config"
	"vango/internal/theme"
)

// Command represents the theme command
type Command struct {
	themeManager *theme.ThemeManager
}

// NewCommand creates a new theme command
func NewCommand() *Command {
	return &Command{}
}

// Execute runs the theme command
func (c *Command) Execute(args []string) error {
	if len(args) < 1 {
		return c.showHelp()
	}

	// Load configuration
	cfg, err := config.Load("config.toml")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	c.themeManager = theme.NewThemeManager(cfg)
	if err := c.themeManager.LoadThemes(); err != nil {
		return fmt.Errorf("failed to load themes: %w", err)
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "list":
		return c.listThemes()
	case "install":
		if len(subArgs) < 1 {
			return fmt.Errorf("theme name required")
		}
		return c.installTheme(subArgs[0])
	case "use":
		if len(subArgs) < 1 {
			return fmt.Errorf("theme name required")
		}
		return c.useTheme(subArgs[0])
	case "create":
		if len(subArgs) < 1 {
			return fmt.Errorf("theme name required")
		}
		template := "basic"
		if len(subArgs) > 1 {
			template = subArgs[1]
		}
		return c.createTheme(subArgs[0], template)
	case "info":
		if len(subArgs) < 1 {
			return fmt.Errorf("theme name required")
		}
		return c.showThemeInfo(subArgs[0])
	case "config":
		return c.showThemeConfig()
	case "validate":
		if len(subArgs) < 1 {
			return fmt.Errorf("theme name required")
		}
		return c.validateTheme(subArgs[0])
	default:
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}
}

// showHelp displays help information
func (c *Command) showHelp() error {
	fmt.Println("VanGo Theme Manager")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  vango theme <command> [arguments]")
	fmt.Println("")
	fmt.Println("Available Commands:")
	fmt.Println("  list                 List all available themes")
	fmt.Println("  install <name>       Install a theme")
	fmt.Println("  use <name>           Set active theme")
	fmt.Println("  create <name> [type] Create a new theme")
	fmt.Println("  info <name>          Show theme information")
	fmt.Println("  config               Show current theme configuration")
	fmt.Println("  validate <name>      Validate theme structure")
	fmt.Println("")
	fmt.Println("Theme Types:")
	fmt.Println("  basic      Simple, clean theme")
	fmt.Println("  blog       Blog-focused theme")
	fmt.Println("  portfolio  Portfolio/showcase theme")
	fmt.Println("  docs       Documentation theme")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  vango theme create mytheme basic")
	fmt.Println("  vango theme use mytheme")
	fmt.Println("  vango theme list")
	
	return nil
}

// listThemes displays all available themes
func (c *Command) listThemes() error {
	themes := c.themeManager.ListThemes()
	
	if len(themes) == 0 {
		fmt.Println("No themes found. Create a theme with 'vango theme create <name>'")
		return nil
	}

	fmt.Println("Available Themes:")
	fmt.Println("")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVERSION\tAUTHOR\tDESCRIPTION")
	fmt.Fprintln(w, "----\t-------\t------\t-----------")

	for _, theme := range themes {
		active := ""
		if c.themeManager.GetActiveTheme() != nil && c.themeManager.GetActiveTheme().Name == theme.Name {
			active = " (active)"
		}
		
		description := theme.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}
		
		fmt.Fprintf(w, "%s%s\t%s\t%s\t%s\n", 
			theme.Name, active, theme.Version, theme.Author, description)
	}

	return w.Flush()
}

// installTheme installs a theme (placeholder for future implementation)
func (c *Command) installTheme(name string) error {
	fmt.Printf("Installing theme: %s\n", name)
	fmt.Println("Theme installation from remote repositories is not yet implemented.")
	fmt.Println("You can create a theme with 'vango theme create <name>' or manually add themes to the themes/ directory.")
	return nil
}

// useTheme sets the active theme
func (c *Command) useTheme(name string) error {
	if err := c.themeManager.SetActiveTheme(name); err != nil {
		return err
	}
	
	fmt.Printf("Successfully set active theme to: %s\n", name)
	
	// Update config file
	cfg, err := config.Load("config.toml")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	cfg.Theme = name
	
	// Save updated config (this would need to be implemented in the config package)
	fmt.Println("Don't forget to update your config.toml file with:")
	fmt.Printf("theme = \"%s\"\n", name)
	
	return nil
}

// createTheme creates a new theme
func (c *Command) createTheme(name, template string) error {
	fmt.Printf("Creating theme '%s' with template '%s'...\n", name, template)
	
	if err := c.themeManager.CreateTheme(name, template); err != nil {
		return err
	}
	
	fmt.Printf("Theme '%s' created successfully!\n", name)
	fmt.Printf("Theme files are located in: themes/%s/\n", name)
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Printf("1. Edit themes/%s/theme.json to customize theme metadata\n", name)
	fmt.Printf("2. Modify templates in themes/%s/layouts/\n", name)
	fmt.Printf("3. Add styles to themes/%s/static/css/style.css\n", name)
	fmt.Printf("4. Use the theme with: vango theme use %s\n", name)
	
	return nil
}

// showThemeInfo displays detailed information about a theme
func (c *Command) showThemeInfo(name string) error {
	theme, exists := c.themeManager.GetTheme(name)
	if !exists {
		return fmt.Errorf("theme not found: %s", name)
	}
	
	fmt.Printf("Theme: %s\n", theme.Name)
	fmt.Printf("Version: %s\n", theme.Version)
	fmt.Printf("Author: %s\n", theme.Author)
	fmt.Printf("Description: %s\n", theme.Description)
	
	if theme.Homepage != "" {
		fmt.Printf("Homepage: %s\n", theme.Homepage)
	}
	
	if theme.License != "" {
		fmt.Printf("License: %s\n", theme.License)
	}
	
	if theme.MinVersion != "" {
		fmt.Printf("Min VanGo Version: %s\n", theme.MinVersion)
	}
	
	if len(theme.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(theme.Tags, ", "))
	}
	
	if len(theme.Features) > 0 {
		fmt.Printf("Features: %s\n", strings.Join(theme.Features, ", "))
	}
	
	fmt.Printf("Path: %s\n", theme.Path)
	
	// Show theme structure
	fmt.Println("\nTheme Structure:")
	fmt.Printf("  Layouts: %s\n", theme.LayoutsDir)
	fmt.Printf("  Static: %s\n", theme.StaticDir)
	fmt.Printf("  Assets: %s\n", theme.AssetsDir)
	
	return nil
}

// showThemeConfig displays the current theme configuration
func (c *Command) showThemeConfig() error {
	activeTheme := c.themeManager.GetActiveTheme()
	if activeTheme == nil {
		fmt.Println("No active theme set.")
		return nil
	}
	
	fmt.Printf("Active Theme: %s\n", activeTheme.Name)
	fmt.Println("")
	
	config, err := c.themeManager.GetThemeConfig()
	if err != nil {
		return fmt.Errorf("failed to get theme config: %w", err)
	}
	
	// Display color scheme
	fmt.Println("Color Scheme:")
	fmt.Printf("  Primary: %s\n", config.Colors.Primary)
	fmt.Printf("  Secondary: %s\n", config.Colors.Secondary)
	fmt.Printf("  Accent: %s\n", config.Colors.Accent)
	fmt.Printf("  Background: %s\n", config.Colors.Background)
	fmt.Printf("  Text: %s\n", config.Colors.Text)
	fmt.Println("")
	
	// Display typography
	fmt.Println("Typography:")
	fmt.Printf("  Font Family: %s\n", config.Typography.FontFamily)
	fmt.Printf("  Font Size: %s\n", config.Typography.FontSize)
	fmt.Printf("  Line Height: %.1f\n", config.Typography.LineHeight)
	fmt.Println("")
	
	// Display layout
	fmt.Println("Layout:")
	fmt.Printf("  Max Width: %s\n", config.Layout.MaxWidth)
	fmt.Printf("  Sidebar: %t\n", config.Layout.Sidebar)
	fmt.Printf("  Navigation: %s\n", config.Layout.Navigation)
	fmt.Printf("  Footer: %t\n", config.Layout.Footer)
	fmt.Println("")
	
	// Display features
	fmt.Println("Features:")
	fmt.Printf("  Dark Mode: %t\n", config.Features.DarkMode)
	fmt.Printf("  Syntax Highlighting: %t\n", config.Features.Syntax)
	fmt.Printf("  MathJax: %t\n", config.Features.MathJax)
	fmt.Printf("  Table of Contents: %t\n", config.Features.TableOfContents)
	fmt.Printf("  Reading Time: %t\n", config.Features.ReadingTime)
	
	return nil
}

// validateTheme validates a theme's structure and configuration
func (c *Command) validateTheme(name string) error {
	theme, exists := c.themeManager.GetTheme(name)
	if !exists {
		return fmt.Errorf("theme not found: %s", name)
	}
	
	fmt.Printf("Validating theme: %s\n", name)
	fmt.Println("")
	
	issues := []string{}
	
	// Check required files
	requiredFiles := map[string]string{
		"theme.json":                    "Theme metadata file",
		"layouts/_default/single.html": "Single page template",
		"layouts/_default/list.html":   "List page template",
	}
	
	for file, description := range requiredFiles {
		path := fmt.Sprintf("%s/%s", theme.Path, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			issues = append(issues, fmt.Sprintf("Missing %s (%s)", file, description))
		} else {
			fmt.Printf("✓ %s exists\n", description)
		}
	}
	
	// Check optional but recommended files
	recommendedFiles := map[string]string{
		"static/css/style.css": "Main stylesheet",
		"layouts/partials":     "Partials directory",
		"config.json":          "Theme configuration",
		"README.md":            "Theme documentation",
	}
	
	for file, description := range recommendedFiles {
		path := fmt.Sprintf("%s/%s", theme.Path, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("⚠ Missing %s (%s) - recommended\n", file, description)
		} else {
			fmt.Printf("✓ %s exists\n", description)
		}
	}
	
	// Validate theme.json structure
	if theme.Name == "" {
		issues = append(issues, "Theme name is empty")
	}
	if theme.Version == "" {
		issues = append(issues, "Theme version is empty")
	}
	if theme.Author == "" {
		fmt.Println("⚠ Theme author is empty - recommended to set")
	}
	if theme.Description == "" {
		fmt.Println("⚠ Theme description is empty - recommended to set")
	}
	
	fmt.Println("")
	
	if len(issues) == 0 {
		fmt.Printf("✅ Theme '%s' validation passed!\n", name)
	} else {
		fmt.Printf("❌ Theme '%s' has %d issues:\n", name, len(issues))
		for _, issue := range issues {
			fmt.Printf("  • %s\n", issue)
		}
		return fmt.Errorf("theme validation failed")
	}
	
	return nil
}