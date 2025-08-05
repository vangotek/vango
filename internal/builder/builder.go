package builder

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vango/internal/config"
	"vango/internal/content"
	"vango/internal/template"
)

// Builder handles site building
type Builder struct {
	config   *config.Config
	parser   *content.Parser
	engine   *template.Engine
	pages    []*content.Page
}

// New creates a new builder
func New(cfg *config.Config) *Builder {
	return &Builder{
		config: cfg,
		parser: content.NewParser(),
		engine: template.NewEngine(cfg),
		pages:  make([]*content.Page, 0),
	}
}

// Build builds the entire site
func (b *Builder) Build() error {
	fmt.Println("Building site...")

	// Clean public directory if configured
	if b.config.CleanBuild {
		if err := b.cleanPublicDir(); err != nil {
			return fmt.Errorf("failed to clean public directory: %w", err)
		}
	}

	// Ensure public directory exists
	if err := os.MkdirAll(b.config.PublicDir, 0755); err != nil {
		return fmt.Errorf("failed to create public directory: %w", err)
	}

	// Load templates
	if err := b.engine.LoadTemplates(); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Parse content files
	if err := b.parseContent(); err != nil {
		return fmt.Errorf("failed to parse content: %w", err)
	}

	// Generate pages
	if err := b.generatePages(); err != nil {
		return fmt.Errorf("failed to generate pages: %w", err)
	}

	// Copy static assets
	if err := b.copyStaticFiles(); err != nil {
		return fmt.Errorf("failed to copy static files: %w", err)
	}

	fmt.Printf("Generated %d pages\n", len(b.pages))
	return nil
}

// cleanPublicDir removes and recreates the public directory
func (b *Builder) cleanPublicDir() error {
	if _, err := os.Stat(b.config.PublicDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(b.config.PublicDir); err != nil {
			return err
		}
	}
	return nil
}

// parseContent walks the content directory and parses all markdown files
func (b *Builder) parseContent() error {
	b.pages = make([]*content.Page, 0)

	return filepath.Walk(b.config.ContentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		// Parse the content file
		page, err := b.parser.ParseFile(path, b.config.ContentDir)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		// Check if page should be built
		if !page.ShouldBuild(b.config.BuildDrafts, b.config.BuildFuture) {
			fmt.Printf("Skipping %s (draft: %v, future: %v)\n", path, page.Draft, page.ParsedDate.After(time.Now()))
			return nil
		}

		b.pages = append(b.pages, page)
		return nil
	})
}

// generatePages renders and writes all pages
func (b *Builder) generatePages() error {
	for _, page := range b.pages {
		if err := b.generatePage(page); err != nil {
			return fmt.Errorf("failed to generate page %s: %w", page.FilePath, err)
		}
	}
	return nil
}

// generatePage renders and writes a single page
func (b *Builder) generatePage(page *content.Page) error {
	// Render the page
	html, err := b.engine.Render(page, b.pages)
	if err != nil {
		return err
	}

	// Determine output path
	outputPath := filepath.Join(b.config.PublicDir, page.Slug, "index.html")
	
	// Handle root-level content files
	if !strings.Contains(page.Slug, "/") {
		outputPath = filepath.Join(b.config.PublicDir, page.Slug, "index.html")
	}

	// Create output directory
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
	}

	// Write HTML file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(html); err != nil {
		return fmt.Errorf("failed to write output file %s: %w", outputPath, err)
	}

	page.OutputPath = outputPath
	fmt.Printf("Generated: %s\n", outputPath)
	return nil
}

// copyStaticFiles copies static assets to the public directory
func (b *Builder) copyStaticFiles() error {
	staticDir := b.config.StaticDir
	
	// Check if static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		fmt.Printf("Static directory %s does not exist, skipping\n", staticDir)
		return nil
	}

	staticOutputDir := filepath.Join(b.config.PublicDir, "static")
	
	return filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Get relative path from static directory
		relPath, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}

		// Create output path
		outputPath := filepath.Join(staticOutputDir, relPath)
		
		// Create output directory
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
		}

		// Copy file
		return b.copyFile(path, outputPath)
	})
}

// copyFile copies a file from src to dst
func (b *Builder) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Copy file permissions
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	return destFile.Chmod(sourceInfo.Mode())
}

// GetPages returns all parsed pages
func (b *Builder) GetPages() []*content.Page {
	return b.pages
}

// GetPageBySlug returns a page by its slug
func (b *Builder) GetPageBySlug(slug string) *content.Page {
	for _, page := range b.pages {
		if page.Slug == slug {
			return page
		}
	}
	return nil
}
