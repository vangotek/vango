package template

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vango/internal/config"
	"vango/internal/content"
	"vango/internal/theme"
)

// Engine handles template rendering
type Engine struct {
	config    *config.Config
	templates *template.Template // Use a single template set
	funcMap   template.FuncMap
}

// TemplateData represents data passed to templates
type TemplateData struct {
	Site   *config.Config
	Page   *content.Page
	Pages  []*content.Page
	Params map[string]interface{}
}

// NewEngine creates a new template engine
func NewEngine(cfg *config.Config, tm *theme.ThemeManager) *Engine {
	engine := &Engine{
		config:    cfg,
		templates: template.New("vango"), // Initialize a single root template set
		funcMap:   createFuncMap(),
	}

	// Add theme functions
	for name, fn := range tm.GetThemeFunctions() {
		engine.funcMap[name] = fn
	}

	engine.templates.Funcs(engine.funcMap) // Apply funcMap to the root template set

	return engine
}

// LoadTemplates loads all templates from the given directory and the default layout directory
func (e *Engine) LoadTemplates(themeLayoutDir string) error {
	// Load theme templates first (higher priority)
	if themeLayoutDir != "" && themeLayoutDir != e.config.LayoutDir {
		if err := e.parseAndAddTemplates(themeLayoutDir); err != nil {
			return fmt.Errorf("failed to parse theme templates: %w", err)
		}
	}

	// Then load default templates (lower priority - won't override existing)
	if err := e.parseAndAddTemplatesWithOverride(e.config.LayoutDir, false); err != nil {
		return fmt.Errorf("failed to parse default templates: %w", err)
	}

	return nil
}

// parseAndAddTemplatesWithOverride walks a directory, parses HTML files, and adds them to the template set with override control
func (e *Engine) parseAndAddTemplatesWithOverride(layoutDir string, allowOverride bool) error {
	return filepath.Walk(layoutDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		// Get template name relative to layouts directory, without .html extension
		relPath, err := filepath.Rel(layoutDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for template %s: %w", path, err)
		}
		templateName := strings.TrimSuffix(relPath, ".html")
		templateName = filepath.ToSlash(templateName) // Ensure forward slashes on all platforms

		// Skip if template already exists and override is not allowed
		if !allowOverride && e.templates.Lookup(templateName) != nil {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", path, err)
		}

		// Parse the template and add it to the main template set
		_, err = e.templates.New(templateName).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		return nil
	})
}

// parseAndAddTemplates walks a directory, parses HTML files, and adds them to the template set
func (e *Engine) parseAndAddTemplates(layoutDir string) error {
	return e.parseAndAddTemplatesWithOverride(layoutDir, true)
}

// Render renders a page using the appropriate template
func (e *Engine) Render(page *content.Page, pages []*content.Page) (string, error) {
	// Determine which template to use
	templateName := e.getTemplateName(page)
	
	tmpl := e.templates.Lookup(templateName) // Use Lookup on the single template set
	if tmpl == nil {
		return "", fmt.Errorf("template not found: %s", templateName)
	}
	
	// Prepare template data
	data := &TemplateData{
		Site:   e.config,
		Page:   page,
		Pages:  pages,
		Params: make(map[string]interface{}),
	}
	
	// Execute template
	var buf strings.Builder
	
	// Handle template inheritance for base templates
	if templateName == "_default/baseof" {
		// For base templates, we need to execute with proper context
		// The base template will call the appropriate content template
		err := e.templates.ExecuteTemplate(&buf, "_default/baseof", data)
		if err != nil {
			return "", fmt.Errorf("failed to execute base template: %w", err)
		}
	} else {
		// For non-base templates, execute directly
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
		}
	}
	
	return buf.String(), nil
}

// getTemplateName determines which template to use for a page
func (e *Engine) getTemplateName(page *content.Page) string {
	// Check for page-specific template
	if tmplName, ok := page.Params["layout"].(string); ok {
		if e.templates.Lookup(tmplName) != nil { // Use Lookup
			return tmplName
		}
	}
	
	// Check for section-specific template
	if strings.Contains(page.Slug, "/") {
		section := strings.Split(page.Slug, "/")[0]
		sectionTemplate := section + "/single"
		if e.templates.Lookup(sectionTemplate) != nil { // Use Lookup
			return sectionTemplate
		}
	}
	
	// For themes with base templates, try to use baseof as the main template
	if e.templates.Lookup("_default/baseof") != nil {
		fmt.Printf("🎨 Using base template: _default/baseof\n")
		return "_default/baseof"
	}
	
	// Default to single template
	return "_default/single"
}

// createFuncMap creates template functions
func createFuncMap() template.FuncMap {
	return template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": strings.Title,
		"trim":  strings.TrimSpace,
		"replace": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"split": strings.Split,
		"join": func(sep string, elems []string) string {
			return strings.Join(elems, sep)
		},
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"contains": strings.Contains,
		"now": time.Now,
		"dateFormat": func(layout string, date time.Time) string {
			return date.Format(layout)
		},
		"humanizeDate": func(date time.Time) string {
			return date.Format("January 2, 2006")
		},
		"timeAgo": func(date time.Time) string {
			duration := time.Since(date)
			switch {
			case duration.Hours() < 24:
				return fmt.Sprintf("%.0f hours ago", duration.Hours())
			case duration.Hours() < 24*7:
				return fmt.Sprintf("%.0f days ago", duration.Hours()/24)
			case duration.Hours() < 24*30:
				return fmt.Sprintf("%.0f weeks ago", duration.Hours()/(24*7))
			case duration.Hours() < 24*365:
				return fmt.Sprintf("%.0f months ago", duration.Hours()/(24*30))
			default:
				return fmt.Sprintf("%.0f years ago", duration.Hours()/(24*365))
			}
		},
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { 
			if b == 0 { return 0 }
			return a / b 
		},
		"seq": func(n int) []int {
			seq := make([]int, n)
			for i := range seq {
				seq[i] = i + 1
			}
			return seq
		},
		"dict": func(values ...interface{}) map[string]interface{} {
			if len(values)%2 != 0 {
				return nil
			}
			dict := make(map[string]interface{})
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					continue
				}
				dict[key] = values[i+1]
			}
			return dict
		},
		"default": func(defaultValue, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeCSS": func(s string) template.CSS {
			return template.CSS(s)
		},
		"safeJS": func(s string) template.JS {
			return template.JS(s)
		},
	}
}

// GetTemplate returns a template by name
func (e *Engine) GetTemplate(name string) (*template.Template, bool) {
	tmpl := e.templates.Lookup(name)
	return tmpl, tmpl != nil
}

// ListTemplates returns all available template names
func (e *Engine) ListTemplates() []string {
	names := make([]string, 0)
	for _, tmpl := range e.templates.Templates() {
		names = append(names, tmpl.Name())
	}
	return names
}
