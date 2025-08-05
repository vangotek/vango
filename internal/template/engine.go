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
)

// Engine handles template rendering
type Engine struct {
	config    *config.Config
	templates map[string]*template.Template
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
func NewEngine(cfg *config.Config) *Engine {
	engine := &Engine{
		config:    cfg,
		templates: make(map[string]*template.Template),
		funcMap:   createFuncMap(),
	}
	
	return engine
}

// LoadTemplates loads all templates from the layouts directory
func (e *Engine) LoadTemplates() error {
	layoutDir := e.config.LayoutDir
	
	return filepath.Walk(layoutDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		
		// Get template name relative to layouts directory
		relPath, err := filepath.Rel(layoutDir, path)
		if err != nil {
			return err
		}
		
		templateName := strings.ReplaceAll(relPath, "\\", "/")
		templateName = strings.TrimSuffix(templateName, ".html")
		
		// Parse template with function map
		tmpl := template.New(filepath.Base(path)).Funcs(e.funcMap)
		_, err = tmpl.ParseFiles(path)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}
		
		e.templates[templateName] = tmpl
		return nil
	})
}

// Render renders a page using the appropriate template
func (e *Engine) Render(page *content.Page, pages []*content.Page) (string, error) {
	// Determine which template to use
	templateName := e.getTemplateName(page)
	
	tmpl, exists := e.templates[templateName]
	if !exists {
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
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}
	
	return buf.String(), nil
}

// getTemplateName determines which template to use for a page
func (e *Engine) getTemplateName(page *content.Page) string {
	// Check for page-specific template
	if tmplName, ok := page.Params["layout"].(string); ok {
		if _, exists := e.templates[tmplName]; exists {
			return tmplName
		}
	}
	
	// Check for section-specific template
	if strings.Contains(page.Slug, "/") {
		section := strings.Split(page.Slug, "/")[0]
		sectionTemplate := section + "/single"
		if _, exists := e.templates[sectionTemplate]; exists {
			return sectionTemplate
		}
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
	tmpl, exists := e.templates[name]
	return tmpl, exists
}

// ListTemplates returns all available template names
func (e *Engine) ListTemplates() []string {
	names := make([]string, 0, len(e.templates))
	for name := range e.templates {
		names = append(names, name)
	}
	return names
}
