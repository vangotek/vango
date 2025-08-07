package content

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"
)

// Enhanced Page structure with additional features
type Page struct {
	// Front matter fields
	Title       string                 `toml:"title" yaml:"title"`
	Date        string                 `toml:"date" yaml:"date"`
	ParsedDate  time.Time
	Draft       bool                   `toml:"draft" yaml:"draft"`
	Description string                 `toml:"description" yaml:"description"`
	Tags        []string               `toml:"tags" yaml:"tags"`
	Categories  []string               `toml:"categories" yaml:"categories"`
	Author      string                 `toml:"author" yaml:"author"`
	Weight      int                    `toml:"weight" yaml:"weight"`
	Params      map[string]interface{} `toml:"params" yaml:"params"`
	
	// Enhanced metadata
	Language    string   `toml:"language" yaml:"language"`
	Translationkey string `toml:"translationkey" yaml:"translationkey"`
	Aliases     []string `toml:"aliases" yaml:"aliases"`
	Keywords    []string `toml:"keywords" yaml:"keywords"`
	
	// SEO and social
	MetaDescription string            `toml:"meta_description" yaml:"meta_description"`
	OpenGraph       map[string]string `toml:"opengraph" yaml:"opengraph"`
	TwitterCard     map[string]string `toml:"twitter_card" yaml:"twitter_card"`
	CanonicalURL    string            `toml:"canonical_url" yaml:"canonical_url"`
	
	// Publishing control
	PublishDate time.Time `toml:"publish_date" yaml:"publish_date"`
	ExpiryDate  time.Time `toml:"expiry_date" yaml:"expiry_date"`
	LastMod     time.Time `toml:"lastmod" yaml:"lastmod"`
	
	// Content organization
	Section     string `toml:"section" yaml:"section"`
	Type        string `toml:"type" yaml:"type"`
	Layout      string `toml:"layout" yaml:"layout"`
	
	// Computed fields
	Content     template.HTML
	Summary     template.HTML
	TableOfContents template.HTML
	WordCount   int
	ReadingTime int
	Slug        string
	URL         string
	Permalink   string
	FilePath    string
	OutputPath  string
	RelPermalink string
	
	// Enhanced features
	Hash        string            // Content hash for change detection
	Headings    []Heading         // Extracted headings for TOC
	Images      []Image           // Extracted images
	Links       []Link            // Extracted links
	CodeBlocks  []CodeBlock       // Extracted code blocks
	Related     []*Page           // Related pages
	Translations []*Page          // Page translations
	PrevInSection *Page           // Previous page in section
	NextInSection *Page           // Next page in section
	
	// Performance tracking
	ParseTime   time.Duration
	RenderTime  time.Duration
	LastBuilt   time.Time
}

// Heading represents a heading in the content
type Heading struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
	ID    string `json:"id"`
	Anchor string `json:"anchor"`
}

// Image represents an image in the content
type Image struct {
	Src    string `json:"src"`
	Alt    string `json:"alt"`
	Title  string `json:"title"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Link represents a link in the content
type Link struct {
	URL    string `json:"url"`
	Text   string `json:"text"`
	Title  string `json:"title"`
	External bool `json:"external"`
}

// CodeBlock represents a code block in the content
type CodeBlock struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Lines    int    `json:"lines"`
}

// Enhanced Parser with additional features
type Parser struct {
	markdown goldmark.Markdown
	options  ParserOptions
}

// ParserOptions configures the parser behavior
type ParserOptions struct {
	ExtractImages     bool
	ExtractLinks      bool
	ExtractHeadings   bool
	ExtractCodeBlocks bool
	GenerateTOC       bool
	EnableSummary     bool
	SummaryLength     int
	EnableAnchors     bool
	SafeMode          bool
}

// NewParser creates a parser with sensible default options.
func NewParser() *Parser {
    // Define your default options here
    defaultOptions := ParserOptions{
        ExtractHeadings:   true,
        ExtractLinks:      true,
        GenerateTOC:       true,
        EnableSummary:     true,
        SummaryLength:     300,
        EnableAnchors:     true,
        SafeMode:          false,
    }
    return NewParserWithOptions(defaultOptions)
}

// NewParserWithOptions creates a parser with custom options
func NewParserWithOptions(options ParserOptions) *Parser {
	extensions := []goldmark.Extender{
		extension.GFM,
		extension.Table,
		extension.Strikethrough,
		extension.Linkify,
		extension.TaskList,
		extension.Footnote,
		extension.DefinitionList,
	}

	// Add typographer extension for smart quotes
	extensions = append(extensions, extension.Typographer)

	parserOptions := []parser.Option{
		parser.WithAutoHeadingID(),
		parser.WithAttribute(),
	}

	var rendererOptions []renderer.Option
	if !options.SafeMode {
		rendererOptions = []renderer.Option{
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		}
	} else {
		rendererOptions = []renderer.Option{
			html.WithHardWraps(),
			html.WithXHTML(),
		}
	}

	md := goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(parserOptions...),
		goldmark.WithRendererOptions(rendererOptions...),
	)

	return &Parser{
		markdown: md,
		options:  options,
	}
}

// ParseFile parses a content file with enhanced features
func (p *Parser) ParseFile(filePath string, contentDir string) (*Page, error) {
	startTime := time.Now()
	
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	
	// Read front matter and body
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

	// Initialize page with enhanced defaults
	page := &Page{
		ParsedDate:    time.Now(),
		Params:       make(map[string]interface{}),
		FilePath:     filePath,
		OpenGraph:    make(map[string]string),
		TwitterCard:  make(map[string]string),
		Headings:     make([]Heading, 0),
		Images:       make([]Image, 0),
		Links:        make([]Link, 0),
		CodeBlocks:   make([]CodeBlock, 0),
		Related:      make([]*Page, 0),
		Translations: make([]*Page, 0),
		LastBuilt:    time.Now(),
	}

	// Parse front matter
	if frontMatter.Len() > 0 {
		if err := p.parseFrontMatter(frontMatter.String(), frontMatterDelim, page); err != nil {
			return nil, fmt.Errorf("failed to parse front matter in %s: %w", filePath, err)
		}
	}

	// Generate content hash for change detection
	bodyContent := body.String()
	page.Hash = p.generateContentHash(bodyContent)

	// Process content with enhanced features
	if err := p.processContent(bodyContent, page); err != nil {
		return nil, fmt.Errorf("failed to process content in %s: %w", filePath, err)
	}

	// Generate URL and slug
	if err := p.generateURLs(page, contentDir); err != nil {
		return nil, fmt.Errorf("failed to generate URLs for %s: %w", filePath, err)
	}

	// Set defaults
	p.setDefaults(page)

	page.ParseTime = time.Since(startTime)
	return page, nil
}

// parseFrontMatter parses TOML or YAML front matter
func (p *Parser) parseFrontMatter(content, delimiter string, page *Page) error {
	var err error
	
	switch delimiter {
	case "+++":
		err = toml.Unmarshal([]byte(content), page)
	case "---":
		err = yaml.Unmarshal([]byte(content), page)
	default:
		// Auto-detect format
		if strings.Contains(content, ":") && !strings.Contains(content, "=") {
			err = yaml.Unmarshal([]byte(content), page)
		} else {
			err = toml.Unmarshal([]byte(content), page)
		}
	}
	
	if err != nil {
		return err
	}

	// Parse dates
	if err := p.parseDates(page); err != nil {
		return err
	}

	return nil
}

// parseDates parses various date fields
func (p *Parser) parseDates(page *Page) error {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"01/02/2006",
		"2006/01/02",
	}

	// Parse main date
	if page.Date != "" {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, page.Date); err == nil {
				page.ParsedDate = t
				break
			}
		}
	}

	// Parse publish date
	if publishDate, ok := page.Params["publish_date"].(string); ok {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, publishDate); err == nil {
				page.PublishDate = t
				break
			}
		}
	}

	// Parse expiry date
	if expiryDate, ok := page.Params["expiry_date"].(string); ok {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, expiryDate); err == nil {
				page.ExpiryDate = t
				break
			}
		}
	}

	// Parse last modified date
	if lastMod, ok := page.Params["lastmod"].(string); ok {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, lastMod); err == nil {
				page.LastMod = t
				break
			}
		}
	}

	return nil
}

// processContent converts markdown and extracts features
func (p *Parser) processContent(content string, page *Page) error {
	// Convert markdown to HTML
	var htmlBuf strings.Builder
	if err := p.markdown.Convert([]byte(content), &htmlBuf); err != nil {
		return err
	}
	
	htmlContent := htmlBuf.String()
	page.Content = template.HTML(htmlContent)

	// Extract features if enabled
	if p.options.ExtractHeadings {
		page.Headings = p.extractHeadings(htmlContent)
	}

	if p.options.ExtractImages {
		page.Images = p.extractImages(htmlContent)
	}

	if p.options.ExtractLinks {
		page.Links = p.extractLinks(htmlContent)
	}

	if p.options.ExtractCodeBlocks {
		page.CodeBlocks = p.extractCodeBlocks(content)
	}

	if p.options.GenerateTOC && len(page.Headings) > 0 {
		page.TableOfContents = p.generateTableOfContents(page.Headings)
	}

	if p.options.EnableSummary {
		page.Summary = p.generateSummary(content, p.options.SummaryLength)
	}

	// Calculate reading metrics
	words := strings.Fields(p.stripHTML(content))
	page.WordCount = len(words)
	page.ReadingTime = p.calculateReadingTime(page.WordCount)

	return nil
}

// extractHeadings extracts headings from HTML content
func (p *Parser) extractHeadings(html string) []Heading {
	re := regexp.MustCompile(`<h([1-6])(?:\s+id="([^"]*)")?[^>]*>([^<]+)</h[1-6]>`)
	matches := re.FindAllStringSubmatch(html, -1)
	
	var headings []Heading
	for _, match := range matches {
		level := parseInt(match[1])
		id := match[2]
		text := strings.TrimSpace(match[3])
		
		if id == "" {
			id = p.slugify(text)
		}
		
		headings = append(headings, Heading{
			Level:  level,
			Text:   text,
			ID:     id,
			Anchor: "#" + id,
		})
	}
	
	return headings
}

// extractImages extracts images from HTML content
func (p *Parser) extractImages(html string) []Image {
	re := regexp.MustCompile(`<img[^>]+src="([^"]*)"[^>]*(?:alt="([^"]*)")?[^>]*(?:title="([^"]*)")?[^>]*>`)
	matches := re.FindAllStringSubmatch(html, -1)
	
	var images []Image
	for _, match := range matches {
		images = append(images, Image{
			Src:   match[1],
			Alt:   match[2],
			Title: match[3],
		})
	}
	
	return images
}

// extractLinks extracts links from HTML content
func (p *Parser) extractLinks(html string) []Link {
	re := regexp.MustCompile(`<a[^>]+href="([^"]*)"[^>]*(?:title="([^"]*)")?[^>]*>([^<]*)</a>`)
	matches := re.FindAllStringSubmatch(html, -1)
	
	var links []Link
	for _, match := range matches {
		url := match[1]
		isExternal := strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
		
		links = append(links, Link{
			URL:      url,
			Title:    match[2],
			Text:     match[3],
			External: isExternal,
		})
	}
	
	return links
}

// extractCodeBlocks extracts code blocks from markdown content
func (p *Parser) extractCodeBlocks(content string) []CodeBlock {
	re := regexp.MustCompile("```(\\w+)?\\n([\\s\\S]*?)```")
	matches := re.FindAllStringSubmatch(content, -1)
	
	var codeBlocks []CodeBlock
	for _, match := range matches {
		language := match[1]
		code := strings.TrimSpace(match[2])
		lines := len(strings.Split(code, "\n"))
		
		codeBlocks = append(codeBlocks, CodeBlock{
			Language: language,
			Code:     code,
			Lines:    lines,
		})
	}
	
	return codeBlocks
}

// generateTableOfContents creates a TOC from headings
func (p *Parser) generateTableOfContents(headings []Heading) template.HTML {
	if len(headings) == 0 {
		return ""
	}
	
	var toc strings.Builder
	toc.WriteString(`<nav class="table-of-contents" role="navigation" aria-label="Table of contents">`)
	toc.WriteString(`<h3>Table of Contents</h3><ul>`)
	
	for _, heading := range headings {
		toc.WriteString(fmt.Sprintf(
			`<li class="toc-level-%d"><a href="%s">%s</a></li>`,
			heading.Level, heading.Anchor, heading.Text,
		))
	}
	
	toc.WriteString(`</ul></nav>`)
	return template.HTML(toc.String())
}

// generateSummary creates a summary from content
func (p *Parser) generateSummary(content string, maxLength int) template.HTML {
	// Remove code blocks first
	re := regexp.MustCompile("```[\\s\\S]*?```")
	content = re.ReplaceAllString(content, "")
	
	// Remove markdown formatting
	content = p.stripMarkdown(content)
	
	// Truncate to maxLength
	if len(content) <= maxLength {
		return template.HTML(content)
	}
	
	// Find the last complete sentence within maxLength
	truncated := content[:maxLength]
	lastPeriod := strings.LastIndex(truncated, ".")
	lastExclamation := strings.LastIndex(truncated, "!")
	lastQuestion := strings.LastIndex(truncated, "?")
	
	lastSentence := max(max(lastPeriod, lastExclamation), lastQuestion)
	if lastSentence > maxLength-50 { // Only use if it's not too short
		truncated = content[:lastSentence+1]
	} else {
		truncated += "..."
	}
	
	return template.HTML(truncated)
}

// generateURLs creates URL and slug for the page
func (p *Parser) generateURLs(page *Page, contentDir string) error {
	relPath, err := filepath.Rel(contentDir, page.FilePath)
	if err != nil {
		return err
	}

	page.Slug = strings.TrimSuffix(relPath, filepath.Ext(relPath))
	page.Slug = strings.ReplaceAll(page.Slug, "\\", "/")
	
	// Generate section from file path
	pathParts := strings.Split(page.Slug, "/")
	if len(pathParts) > 1 {
		page.Section = pathParts[0]
	}
	
	// Generate URLs
	page.URL = "/" + page.Slug + "/"
	page.RelPermalink = page.URL
	page.Permalink = page.URL // Would be full URL with baseURL in production

	return nil
}

// setDefaults sets default values for the page
func (p *Parser) setDefaults(page *Page) {
	if page.Title == "" {
		page.Title = strings.ReplaceAll(filepath.Base(page.Slug), "-", " ")
		page.Title = strings.Title(page.Title)
	}
	
	if page.Type == "" {
		page.Type = page.Section
		if page.Type == "" {
			page.Type = "page"
		}
	}
	
	if page.Language == "" {
		page.Language = "en"
	}
	
	if page.MetaDescription == "" && len(page.Summary) > 0 {
		page.MetaDescription = string(page.Summary)
	}
}

// Helper functions
func (p *Parser) generateContentHash(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

func (p *Parser) stripHTML(content string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(content, "")
}

func (p *Parser) stripMarkdown(content string) string {
	// Remove various markdown elements
	patterns := []string{
		`#{1,6}\s*`,      // Headers
		`\*\*([^*]+)\*\*`, // Bold
		`\*([^*]+)\*`,     // Italic
		"`([^`]+)`",       // Inline code
		`\[([^\]]+)\]\([^\)]+\)`, // Links
		`!\[([^\]]*)\]\([^\)]+\)`, // Images
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		content = re.ReplaceAllString(content, "$1")
	}
	
	return content
}

func (p *Parser) slugify(text string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\s-]`)
	text = re.ReplaceAllString(text, "")
	text = strings.ToLower(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, "-")
	text = regexp.MustCompile(`-+`).ReplaceAllString(text, "-")
	return strings.Trim(text, "-")
}

func (p *Parser) calculateReadingTime(wordCount int) int {
	// Average reading speed: 200 words per minute
	minutes := (wordCount + 199) / 200
	if minutes < 1 {
		return 1
	}
	return minutes
}


// Enhanced page methods
func (page *Page) ShouldBuild(buildDrafts, buildFuture bool) bool {
	if page.Draft && !buildDrafts {
		return false
	}
	
	if !page.PublishDate.IsZero() && page.PublishDate.After(time.Now()) && !buildFuture {
		return false
	}
	
	if !page.ExpiryDate.IsZero() && page.ExpiryDate.Before(time.Now()) {
		return false
	}
	
	return true
}

func (page *Page) IsExpired() bool {
	return !page.ExpiryDate.IsZero() && page.ExpiryDate.Before(time.Now())
}

func (page *Page) IsFuture() bool {
	return !page.PublishDate.IsZero() && page.PublishDate.After(time.Now())
}

func (page *Page) HasChanged(hash string) bool {
	return page.Hash != hash
}

func (page *Page) GetRelatedByTags(allPages []*Page, limit int) []*Page {
	if len(page.Tags) == 0 {
		return []*Page{}
	}
	
	type pageScore struct {
		page  *Page
		score int
	}
	
	var scored []pageScore
	
	for _, other := range allPages {
		if other.FilePath == page.FilePath {
			continue
		}
		
		score := 0
		for _, tag := range page.Tags {
			for _, otherTag := range other.Tags {
				if tag == otherTag {
					score++
				}
			}
		}
		
		if score > 0 {
			scored = append(scored, pageScore{other, score})
		}
	}
	
	// Sort by score (highest first)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})
	
	// Return top results
	var related []*Page
	for i, ps := range scored {
		if i >= limit {
			break
		}
		related = append(related, ps.page)
	}
	
	return related
}

// Utility functions
func parseInt(s string) int {
	switch s {
	case "1": return 1
	case "2": return 2
	case "3": return 3
	case "4": return 4
	case "5": return 5
	case "6": return 6
	default: return 1
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GetParam returns a parameter value by key
func (p *Page) GetParam(key string) interface{} {
	return p.Params[key]
}

// SetParam sets a parameter value
func (p *Page) SetParam(key string, value interface{}) {
	p.Params[key] = value
}