package content

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Page represents a content page
type Page struct {
	Title       string                 `toml:"title"`
	Date        string              `toml:"date"`
	ParsedDate  time.Time
	Draft       bool                   `toml:"draft"`
	Description string                 `toml:"description"`
	Tags        []string               `toml:"tags"`
	Categories  []string               `toml:"categories"`
	Author      string                 `toml:"author"`
	Weight      int                    `toml:"weight"`
	Params      map[string]interface{} `toml:"params"`
	
	// Computed fields
	Content     template.HTML
	Summary     template.HTML
	WordCount   int
	ReadingTime int
	Slug        string
	URL         string
	FilePath    string
	OutputPath  string
}

// Parser handles content parsing
type Parser struct {
	markdown goldmark.Markdown
}

// NewParser creates a new content parser
func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.Linkify,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	return &Parser{
		markdown: md,
	}
}

// ParseFile parses a content file and returns a Page
func (p *Parser) ParseFile(filePath string, contentDir string) (*Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	
	// Read front matter
	var frontMatter strings.Builder
	var body strings.Builder
	var inFrontMatter bool
	var frontMatterDelim string

	if scanner.Scan() {
		firstLine := scanner.Text()
		if firstLine == "+++" || firstLine == "---" {
			inFrontMatter = true
			frontMatterDelim = firstLine
		} else {
			body.WriteString(firstLine + "\n")
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		
		if inFrontMatter {
			if line == frontMatterDelim {
				inFrontMatter = false
				continue
			}
			frontMatter.WriteString(line + "\n")
		} else {
			body.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	// Parse front matter
	page := &Page{
		ParsedDate:   time.Now(),
		Params: make(map[string]interface{}),
		FilePath: filePath,
	}

	if frontMatter.Len() > 0 {
		var err error
		switch frontMatterDelim {
		case "+++":
			err = toml.Unmarshal([]byte(frontMatter.String()), page)
		case "---":
			// YAML support could be added here
			return nil, fmt.Errorf("YAML front matter not yet supported")
		}
		
		if err != nil {
			return nil, fmt.Errorf("failed to parse front matter in %s: %w", filePath, err)
		}

		// Parse date
		if page.Date != "" {
			layouts := []string{time.RFC3339, "2006-01-02"}
			for _, layout := range layouts {
				t, err := time.Parse(layout, page.Date)
				if err == nil {
					page.ParsedDate = t
					break
				}
			}
			if page.ParsedDate.IsZero() {
				return nil, fmt.Errorf("failed to parse date in %s: %s", filePath, page.Date)
			}
		} else {
			page.ParsedDate = time.Now()
		}
	}

	// Convert markdown to HTML
	var htmlBuf strings.Builder
	if err := p.markdown.Convert([]byte(body.String()), &htmlBuf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown in %s: %w", filePath, err)
	}

	page.Content = template.HTML(htmlBuf.String())

	// Calculate reading time (average 200 words per minute)
	words := strings.Fields(body.String())
	page.WordCount = len(words)
	page.ReadingTime = (page.WordCount + 199) / 200 // Round up

	// Generate slug and URL
	relPath, err := filepath.Rel(contentDir, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path: %w", err)
	}

	page.Slug = strings.TrimSuffix(relPath, filepath.Ext(relPath))
	page.URL = "/" + strings.ReplaceAll(page.Slug, "\\", "/") + "/"

	// Set default title if not provided
	if page.Title == "" {
		page.Title = strings.ReplaceAll(filepath.Base(page.Slug), "-", " ")
		page.Title = strings.Title(page.Title)
	}

	return page, nil
}

// ShouldBuild determines if a page should be built based on its properties
func (p *Page) ShouldBuild(buildDrafts, buildFuture bool) bool {
	if p.Draft && !buildDrafts {
		return false
	}
	
	if p.Date != "" && p.ParsedDate.After(time.Now()) && !buildFuture {
		return false
	}
	
	return true
}

// GetParam returns a parameter value by key
func (p *Page) GetParam(key string) interface{} {
	return p.Params[key]
}

// SetParam sets a parameter value
func (p *Page) SetParam(key string, value interface{}) {
	p.Params[key] = value
}