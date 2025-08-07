package builder

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"vango/internal/config"
	"vango/internal/content"
	"vango/internal/template"
	"vango/internal/theme"
)

// Builder handles site building
type Builder struct {
	config       *config.Config
	parser       *content.Parser
	engine       *template.Engine
	pages        []*content.Page
	themeManager *theme.ThemeManager
	
	// Performance enhancements
	workers      int
	cache        map[string]time.Time // File modification cache
	cacheMutex   sync.RWMutex
}

// New creates a new builder
func New(cfg *config.Config) *Builder {
	workers := runtime.NumCPU()
	if workers > 8 {
		workers = 8 // Cap at 8 for optimal performance
	}
	
	tm := theme.NewThemeManager(cfg)
	return &Builder{
		config:       cfg,
		parser:       content.NewParser(),
		engine:       template.NewEngine(cfg, tm),
		pages:        make([]*content.Page, 0),
		themeManager: tm,
		workers:      workers,
		cache:        make(map[string]time.Time),
	}
}

// Build builds the entire site
func (b *Builder) Build() error {
	start := time.Now()
	fmt.Printf("🏗️  Building site with %d workers...\n", b.workers)

	// Load themes and set active theme
	if err := b.themeManager.LoadThemes(); err != nil {
		return fmt.Errorf("failed to load themes: %w", err)
	}
	if b.config.Theme != "" {
		if err := b.themeManager.SetActiveTheme(b.config.Theme); err != nil {
			return fmt.Errorf("failed to set active theme: %w", err)
		}
		fmt.Printf("📦 Using theme: %s\n", b.config.Theme)
	}

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

	// Load templates with caching
	if err := b.engine.LoadTemplates(b.themeManager.GetThemeTemplatesPath()); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Parse content files in parallel
	if err := b.parseContentParallel(); err != nil {
		return fmt.Errorf("failed to parse content: %w", err)
	}

	// Generate pages in parallel
	if err := b.generatePagesParallel(); err != nil {
		return fmt.Errorf("failed to generate pages: %w", err)
	}

	// Copy static assets and theme assets in parallel
	errChan := make(chan error, 2)
	go func() {
		errChan <- b.copyStaticFiles()
	}()
	go func() {
		errChan <- b.themeManager.CopyThemeAssets(b.config.PublicDir)
	}()

	// Wait for both operations to complete
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return fmt.Errorf("failed to copy assets: %w", err)
		}
	}

	duration := time.Since(start)
	fmt.Printf("✅ Generated %d pages in %v\n", len(b.pages), duration)
	return nil
}

// parseContentParallel parses content files using worker goroutines
func (b *Builder) parseContentParallel() error {
	// Collect all markdown files
	var files []string
	err := filepath.Walk(b.config.ContentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".md") {
			// Check cache for file modification time
			if b.isFileModified(path, info.ModTime()) {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(files) == 0 {
		fmt.Println("📝 No content files to process")
		return nil
	}

	fmt.Printf("📝 Processing %d content files...\n", len(files))

	// Create worker pool
	fileChan := make(chan string, len(files))
	resultChan := make(chan *content.Page, len(files))
	errorChan := make(chan error, len(files))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < b.workers; i++ {
		wg.Add(1)
		go b.contentWorker(&wg, fileChan, resultChan, errorChan)
	}

	// Send files to workers
	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	// Wait for workers to complete
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Collect results
	var pages []*content.Page
	var errors []error

	for page := range resultChan {
		if page != nil {
			pages = append(pages, page)
		}
	}

	for err := range errorChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("content parsing errors: %v", errors[0])
	}

	b.pages = pages
	return nil
}

// contentWorker processes content files
func (b *Builder) contentWorker(wg *sync.WaitGroup, fileChan <-chan string, resultChan chan<- *content.Page, errorChan chan<- error) {
	defer wg.Done()
	
	for filePath := range fileChan {
		page, err := b.parser.ParseFile(filePath, b.config.ContentDir)
		if err != nil {
			errorChan <- fmt.Errorf("failed to parse %s: %w", filePath, err)
			continue
		}

		// Check if page should be built
		if !page.ShouldBuild(b.config.BuildDrafts, b.config.BuildFuture) {
			continue
		}

		resultChan <- page
	}
}

// generatePagesParallel renders pages using worker goroutines
func (b *Builder) generatePagesParallel() error {
	if len(b.pages) == 0 {
		return nil
	}

	fmt.Printf("🎨 Rendering %d pages...\n", len(b.pages))

	// Create worker pool for page generation
	pageChan := make(chan *content.Page, len(b.pages))
	errorChan := make(chan error, len(b.pages))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < b.workers; i++ {
		wg.Add(1)
		go b.pageWorker(&wg, pageChan, errorChan)
	}

	// Send pages to workers
	for _, page := range b.pages {
		pageChan <- page
	}
	close(pageChan)

	// Wait for workers to complete
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Collect errors
	var errors []error
	for err := range errorChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("page generation errors: %v", errors[0])
	}

	return nil
}

// pageWorker renders individual pages
func (b *Builder) pageWorker(wg *sync.WaitGroup, pageChan <-chan *content.Page, errorChan chan<- error) {
	defer wg.Done()
	
	for page := range pageChan {
		if err := b.generatePage(page); err != nil {
			errorChan <- fmt.Errorf("failed to generate page %s: %w", page.FilePath, err)
		}
	}
}

// isFileModified checks if a file has been modified since last build
func (b *Builder) isFileModified(path string, modTime time.Time) bool {
	b.cacheMutex.RLock()
	cached, exists := b.cache[path]
	b.cacheMutex.RUnlock()
	
	if !exists || modTime.After(cached) {
		b.cacheMutex.Lock()
		b.cache[path] = modTime
		b.cacheMutex.Unlock()
		return true
	}
	return false
}

// IncrementalBuild performs incremental build based on changed files
func (b *Builder) IncrementalBuild(changedFiles []string) error {
	start := time.Now()
	fmt.Printf("🔄 Incremental build for %d changed files...\n", len(changedFiles))

	var needsFullRebuild bool
	var contentFiles []string

	for _, file := range changedFiles {
		switch {
		case strings.HasSuffix(file, ".html"):
			// Template changed, need full rebuild
			needsFullRebuild = true
		case strings.HasSuffix(file, ".toml"):
			// Config changed, need full rebuild
			needsFullRebuild = true
		case strings.HasSuffix(file, ".md"):
			// Content file changed
			contentFiles = append(contentFiles, file)
		case strings.Contains(file, b.config.StaticDir):
			// Static file changed, just copy
			if err := b.copyStaticFiles(); err != nil { // Removed argument (file). Check for bugs in this line.
				return fmt.Errorf("failed to copy static file: %w", err)
			}
		}
	}

	if needsFullRebuild {
		return b.Build()
	}

	// Process only changed content files
	for _, file := range contentFiles {
		if err := b.rebuildContentFile(file); err != nil {
			return fmt.Errorf("failed to rebuild content file %s: %w", file, err)
		}
	}

	duration := time.Since(start)
	fmt.Printf("✅ Incremental build completed in %v\n", duration)
	return nil
}

// Additional helper methods for incremental builds...
func (b *Builder) rebuildContentFile(filePath string) error {
	page, err := b.parser.ParseFile(filePath, b.config.ContentDir)
	if err != nil {
		return err
	}

	if !page.ShouldBuild(b.config.BuildDrafts, b.config.BuildFuture) {
		return nil
	}

	return b.generatePage(page)
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