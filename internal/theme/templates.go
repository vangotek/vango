// internal/theme/templates.go
package theme

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"
)

// GetThemeFunctions returns enhanced template functions for themes
func (tm *ThemeManager) GetThemeFunctions() template.FuncMap {
	return template.FuncMap{
		// Theme-specific functions
		"themeAsset": tm.getThemeAssetURL,
		"themeConfig": tm.getThemeConfigValue,
		"hasFeature": tm.hasFeature,
		"themeColor": tm.getThemeColor,
		
		// Enhanced content functions
		"excerpt":        tm.createExcerpt,
		"readingTime":    tm.calculateReadingTime,
		"wordCount":      tm.countWords,
		"tableOfContents": tm.generateTOC,
		"relatedPosts":   tm.getRelatedPosts,
		
		// SEO and social functions
		"metaDescription": tm.generateMetaDescription,
		"jsonLD":         tm.generateJSONLD,
		"openGraph":      tm.generateOpenGraph,
		"twitterCard":    tm.generateTwitterCard,
		
		// Media and asset functions
		"imageOptimize":  tm.optimizeImage,
		"responsiveImg":  tm.responsiveImage,
		"assetFingerprint": tm.assetFingerprint,
		
		// Date and time enhancements
		"isRecent":       tm.isRecent,
		"formatDate":     tm.formatDate,
		"timeFromNow":    tm.timeFromNow,
		"isoDate":        tm.isoDate,
		
		// Content transformation
		"markdownify":    tm.markdownify,
		"highlight":      tm.syntaxHighlight,
		"sanitizeHTML":   tm.sanitizeHTML,
		"truncateWords":  tm.truncateWords,
		"slugify":        tm.slugify,
		
		// Math and utilities
		"percentage":     tm.percentage,
		"round":          tm.round,
		"random":         tm.random,
		"uuid":           tm.generateUUID,
		
		// Collections and data
		"groupBy":        tm.groupBy,
		"sortBy":         tm.sortBy,
		"filterBy":       tm.filterBy,
		"unique":         tm.unique,
		"paginate":       tm.paginate,
		
		// Conditional helpers
		"ifNotEmpty":     tm.ifNotEmpty,
		"ifAny":          tm.ifAny,
		"ifAll":          tm.ifAll,
		"switch":         tm.switchCase,
	}
}

// Theme-specific functions
func (tm *ThemeManager) getThemeAssetURL(path string) string {
	if tm.activeTheme == nil {
		return "/static/" + path
	}
	return "/theme/" + path
}

func (tm *ThemeManager) getThemeConfigValue(key string) interface{} {
	config, err := tm.GetThemeConfig()
	if err != nil {
		return nil
	}
	
	// Navigate nested keys using dot notation
	parts := strings.Split(key, ".")
	var current interface{} = config
	
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[part]
		default:
			return nil
		}
	}
	
	return current
}

func (tm *ThemeManager) hasFeature(feature string) bool {
	config, err := tm.GetThemeConfig()
	if err != nil {
		return false
	}
	
	switch feature {
	case "dark_mode":
		return config.Features.DarkMode
	case "syntax":
		return config.Features.Syntax
	case "mathjax":
		return config.Features.MathJax
	case "toc":
		return config.Features.TableOfContents
	case "reading_time":
		return config.Features.ReadingTime
	case "share_buttons":
		return config.Features.ShareButtons
	case "related_posts":
		return config.Features.RelatedPosts
	case "analytics":
		return config.Features.Analytics
	default:
		return false
	}
}

func (tm *ThemeManager) getThemeColor(name string) string {
	config, err := tm.GetThemeConfig()
	if err != nil {
		return "#000000"
	}
	
	switch name {
	case "primary":
		return config.Colors.Primary
	case "secondary":
		return config.Colors.Secondary
	case "accent":
		return config.Colors.Accent
	case "background":
		return config.Colors.Background
	case "text":
		return config.Colors.Text
	default:
		return "#000000"
	}
}

// Content functions
func (tm *ThemeManager) createExcerpt(content string, maxWords int) string {
	words := strings.Fields(stripHTML(content))
	if len(words) <= maxWords {
		return strings.Join(words, " ")
	}
	return strings.Join(words[:maxWords], " ") + "..."
}

func (tm *ThemeManager) calculateReadingTime(content string) int {
	words := len(strings.Fields(stripHTML(content)))
	// Average reading speed: 200 words per minute
	minutes := (words + 199) / 200
	if minutes < 1 {
		return 1
	}
	return minutes
}

func (tm *ThemeManager) countWords(content string) int {
	return len(strings.Fields(stripHTML(content)))
}

func (tm *ThemeManager) generateTOC(content string) template.HTML {
	// Extract headings and generate table of contents
	re := regexp.MustCompile(`<h([1-6])(?:\s+id="([^"]*)")?[^>]*>([^<]+)</h[1-6]>`)
	matches := re.FindAllStringSubmatch(content, -1)
	
	if len(matches) == 0 {
		return ""
	}
	
	var toc strings.Builder
	toc.WriteString(`<nav class="table-of-contents"><ul>`)
	
	for i, match := range matches {
		level := match[1]
		id := match[2]
		text := match[3]
		
		if id == "" {
			id = tm.slugify(text)
		}
		
		toc.WriteString(fmt.Sprintf(
			`<li class="toc-level-%s"><a href="#%s">%s</a></li>`,
			level, id, text,
		)) 
        i++
	}
	
	toc.WriteString(`</ul></nav>`)
	return template.HTML(toc.String())
}
// SEO functions
func (tm *ThemeManager) generateMetaDescription(page interface{}) string {
	// Implementation to generate meta description from page content
	return ""
}

func (tm *ThemeManager) generateJSONLD(page interface{}) template.HTML {
	// Generate JSON-LD structured data
	return template.HTML(`<script type="application/ld+json">{}</script>`)
}

func (tm *ThemeManager) generateOpenGraph(page interface{}) template.HTML {
	// Generate Open Graph meta tags
	return ""
}

func (tm *ThemeManager) generateTwitterCard(page interface{}) template.HTML {
	// Generate Twitter Card meta tags
	return ""
}

// Media functions
func (tm *ThemeManager) optimizeImage(src string, width, height int) string {
	// Return optimized image URL (would integrate with image processing)
	return src
}

func (tm *ThemeManager) responsiveImage(src string, sizes []int) template.HTML {
	// Generate responsive image with srcset
	return template.HTML(fmt.Sprintf(`<img src="%s" alt="">`, src))
}

func (tm *ThemeManager) assetFingerprint(path string) string {
	// Generate fingerprinted asset URL for cache busting
	hash := md5.Sum([]byte(path + time.Now().String()))
	return path + "?v=" + hex.EncodeToString(hash[:4])
}

// Date functions
func (tm *ThemeManager) isRecent(date time.Time, days int) bool {
	return time.Since(date).Hours() < float64(days*24)
}

func (tm *ThemeManager) formatDate(format string, date time.Time) string {
	return date.Format(format)
}

func (tm *ThemeManager) timeFromNow(date time.Time) string {
	duration := time.Since(date)
	
	switch {
	case duration.Minutes() < 1:
		return "just now"
	case duration.Hours() < 1:
		return fmt.Sprintf("%.0f minutes ago", duration.Minutes())
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
}

func (tm *ThemeManager) isoDate(date time.Time) string {
	return date.Format(time.RFC3339)
}

// Content transformation
func (tm *ThemeManager) markdownify(content string) template.HTML {
	// Convert markdown to HTML (would use goldmark)
	return template.HTML(content)
}

func (tm *ThemeManager) syntaxHighlight(code, language string) template.HTML {
	// Apply syntax highlighting
	return template.HTML(fmt.Sprintf(
		`<pre><code class="language-%s">%s</code></pre>`,
		language, template.HTMLEscapeString(code),
	))
}

func (tm *ThemeManager) sanitizeHTML(content string) template.HTML {
	// Sanitize HTML content (would use bluemonday)
	return template.HTML(content)
}

func (tm *ThemeManager) truncateWords(content string, maxWords int) string {
	words := strings.Fields(content)
	if len(words) <= maxWords {
		return content
	}
	return strings.Join(words[:maxWords], " ") + "..."
}

func (tm *ThemeManager) slugify(text string) string {
	// Convert text to URL-friendly slug
	re := regexp.MustCompile(`[^a-zA-Z0-9\s-]`)
	text = re.ReplaceAllString(text, "")
	text = strings.ToLower(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, "-")
	text = regexp.MustCompile(`-+`).ReplaceAllString(text, "-")
	return strings.Trim(text, "-")
}

// Math functions
func (tm *ThemeManager) percentage(current, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(current) / float64(total) * 100
}

func (tm *ThemeManager) round(num float64, precision int) float64 {
	multiplier := float64(1)
	for i := 0; i < precision; i++ {
		multiplier *= 10
	}
	return float64(int(num*multiplier+0.5)) / multiplier
}

func (tm *ThemeManager) random(min, max int) int {
	// Would use crypto/rand for better randomness
	return min
}

func (tm *ThemeManager) generateUUID() string {
	// Generate UUID (simplified)
	return "uuid-placeholder"
}

// Collection functions
func (tm *ThemeManager) groupBy(items []interface{}, key string) map[string][]interface{} {
	groups := make(map[string][]interface{})
	// Implementation for grouping items by a field
	return groups
}

func (tm *ThemeManager) sortBy(items []interface{}, key string) []interface{} {
	// Implementation for sorting items by a field
	return items
}

func (tm *ThemeManager) filterBy(items []interface{}, key string, value interface{}) []interface{} {
	// Implementation for filtering items
	return items
}

func (tm *ThemeManager) unique(items []interface{}) []interface{} {
	seen := make(map[interface{}]bool)
	var result []interface{}
	
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

func (tm *ThemeManager) paginate(items []interface{}, page, perPage int) []interface{} {
	start := (page - 1) * perPage
	end := start + perPage
	
	if start >= len(items) {
		return []interface{}{}
	}
	
	if end > len(items) {
		end = len(items)
	}
	
	return items[start:end]
}

// Conditional helpers
func (tm *ThemeManager) ifNotEmpty(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return v != ""
	case []interface{}:
		return len(v) > 0
	case map[string]interface{}:
		return len(v) > 0
	case nil:
		return false
	default:
		return true
	}
}

func (tm *ThemeManager) ifAny(values ...interface{}) bool {
	for _, value := range values {
		if tm.ifNotEmpty(value) {
			return true
		}
	}
	return false
}

func (tm *ThemeManager) ifAll(values ...interface{}) bool {
	for _, value := range values {
		if !tm.ifNotEmpty(value) {
			return false
		}
	}
	return true
}

func (tm *ThemeManager) switchCase(value interface{}, cases ...interface{}) interface{} {
	// Implementation for switch-case logic in templates
	if len(cases)%2 != 0 {
		// Odd number of cases means there's a default value
		defaultValue := cases[len(cases)-1]
		cases = cases[:len(cases)-1]
		
		for i := 0; i < len(cases); i += 2 {
			if cases[i] == value {
				return cases[i+1]
			}
		}
		
		return defaultValue
	}
	
	for i := 0; i < len(cases); i += 2 {
		if cases[i] == value {
			return cases[i+1]
		}
	}
	
	return nil
}

func (tm *ThemeManager) getRelatedPosts(currentPage interface{}, allPages []interface{}, limit int) []interface{} {
	// Implementation for finding related posts based on tags/categories
	return []interface{}{}
}

// Helper function to strip HTML tags
func stripHTML(content string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(content, "")
}
// getBasicTemplates returns basic theme templates
func (tm *ThemeManager) getBasicTemplates() map[string]string {
	return map[string]string{
		"layouts/_default/single.html": `<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Page.Title }} | {{ .Site.Title }}</title>
    <meta name="description" content="{{ default .Site.Description .Page.Description }}">
    <link rel="stylesheet" href="{{ themeAsset "css/style.css" }}">
</head>
<body>
    <header class="site-header">
        <nav class="nav-container">
            <a href="/" class="site-title">{{ .Site.Title }}</a>
            <ul class="nav-links">
                <li><a href="/">Home</a></li>
                <li><a href="/about/">About</a></li>
            </ul>
        </nav>
    </header>
    <main class="main-content">
        <article class="post">
            <header class="post-header">
                <h1 class="post-title">{{ .Page.Title }}</h1>
                {{ if hasFeature "reading_time" }}
                <div class="post-meta">
                    {{ if .Page.ReadingTime }}{{ .Page.ReadingTime }} min read{{ end }}
                </div>
                {{ end }}
            </header>
            <div class="post-content">
                {{ .Page.Content }}
            </div>
        </article>
    </main>
    <footer class="site-footer">
        <div class="footer-container">
            <p>&copy; {{ dateFormat "2006" now }} {{ .Site.Author }}</p>
        </div>
    </footer>
</body>
</html>`,
		"layouts/_default/list.html": `<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Site.Title }}</title>
    <meta name="description" content="{{ .Site.Description }}">
    <link rel="stylesheet" href="{{ themeAsset "css/style.css" }}">
</head>
<body>
    <header class="site-header">
        <nav class="nav-container">
            <a href="/" class="site-title">{{ .Site.Title }}</a>
            <ul class="nav-links">
                <li><a href="/">Home</a></li>
                <li><a href="/about/">About</a></li>
            </ul>
        </nav>
    </header>
    <main class="main-content">
        <div class="home-hero">
            <h1 class="hero-title">{{ .Site.Title }}</h1>
            <p class="hero-description">{{ .Site.Description }}</p>
        </div>
        <section class="posts-list">
            <h2>Recent Posts</h2>
            {{ range .Pages }}
            <article class="post-summary">
                <h3><a href="{{ .URL }}">{{ .Title }}</a></h3>
                <div class="post-meta">
                    <time datetime="{{ dateFormat "2006-01-02" .ParsedDate }}">
                        {{ humanizeDate .ParsedDate }}
                    </time>
                </div>
                <p class="post-excerpt">{{ .Summary }}</p>
            </article>
            {{ end }}
        </section>
    </main>
    <footer class="site-footer">
        <div class="footer-container">
            <p>&copy; {{ dateFormat "2006" now }} {{ .Site.Author }}</p>
        </div>
    </footer>
</body>
</html>`,
	}
}

// getBlogTemplates returns blog-focused theme templates
func (tm *ThemeManager) getBlogTemplates() map[string]string {
	return map[string]string{
		"layouts/_default/single.html": `<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Page.Title }} | {{ .Site.Title }}</title>
    <meta name="description" content="{{ default .Site.Description .Page.Description }}">
    <link rel="stylesheet" href="{{ themeAsset "css/style.css" }}">
    {{ if hasFeature "syntax" }}
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
    <script>hljs.highlightAll();</script>
    {{ end }}
</head>
<body>
    <header class="site-header">
        <nav class="nav-container">
            <a href="/" class="site-title">{{ .Site.Title }}</a>
            <ul class="nav-links">
                <li><a href="/">Home</a></li>
                <li><a href="/about/">About</a></li>
                <li><a href="/posts/">Posts</a></li>
            </ul>
        </nav>
    </header>
    <main class="main-content">
        <article class="post">
            <header class="post-header">
                <h1 class="post-title">{{ .Page.Title }}</h1>
                <div class="post-meta">
                    <time datetime="{{ dateFormat "2006-01-02" .Page.ParsedDate }}">
                        {{ humanizeDate .Page.ParsedDate }}
                    </time>
                    {{ if .Page.Author }}
                        by <span class="author">{{ .Page.Author }}</span>
                    {{ end }}
                    {{ if hasFeature "reading_time" }}
                        â€¢ {{ .Page.ReadingTime }} min read
                    {{ end }}
                </div>
                {{ if .Page.Tags }}
                <div class="post-tags">
                    {{ range .Page.Tags }}
                        <span class="tag">#{{ . }}</span>
                    {{ end }}
                </div>
                {{ end }}
            </header>
            <div class="post-content">
                {{ .Page.Content }}
            </div>
            {{ if hasFeature "share" }}
            <div class="post-share">
                <h4>Share this post</h4>
                <a href="https://twitter.com/intent/tweet?text={{ .Page.Title }}&url={{ .Site.BaseURL }}{{ .Page.URL }}" target="_blank">Twitter</a>
                <a href="https://www.facebook.com/sharer/sharer.php?u={{ .Site.BaseURL }}{{ .Page.URL }}" target="_blank">Facebook</a>
            </div>
            {{ end }}
        </article>
    </main>
    <footer class="site-footer">
        <div class="footer-container">
            <p>&copy; {{ dateFormat "2006" now }} {{ .Site.Author }}</p>
        </div>
    </footer>
</body>
</html>`,
		"layouts/_default/list.html": `<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Site.Title }}</title>
    <meta name="description" content="{{ .Site.Description }}">
    <link rel="stylesheet" href="{{ themeAsset "css/style.css" }}">
</head>
<body>
    <header class="site-header">
        <nav class="nav-container">
            <a href="/" class="site-title">{{ .Site.Title }}</a>
            <ul class="nav-links">
                <li><a href="/">Home</a></li>
                <li><a href="/about/">About</a></li>
                <li><a href="/posts/">Posts</a></li>
            </ul>
        </nav>
    </header>
    <main class="main-content">
        <div class="blog-hero">
            <h1 class="hero-title">{{ .Site.Title }}</h1>
            <p class="hero-description">{{ .Site.Description }}</p>
        </div>
        <section class="posts-grid">
            {{ range .Pages }}
            <article class="post-card">
                <h2><a href="{{ .URL }}">{{ .Title }}</a></h2>
                <div class="post-meta">
                    <time datetime="{{ dateFormat "2006-01-02" .ParsedDate }}">
                        {{ humanizeDate .ParsedDate }}
                    </time>
                    {{ if hasFeature "reading_time" }}
                        â€¢ {{ .ReadingTime }} min read
                    {{ end }}
                </div>
                <p class="post-excerpt">{{ .Summary }}</p>
                {{ if .Tags }}
                <div class="post-tags">
                    {{ range .Tags }}
                        <span class="tag">#{{ . }}</span>
                    {{ end }}
                </div>
                {{ end }}
            </article>
            {{ end }}
        </section>
    </main>
    <footer class="site-footer">
        <div class="footer-container">
            <p>&copy; {{ dateFormat "2006" now }} {{ .Site.Author }}</p>
        </div>
    </footer>
</body>
</html>`,
	}
}

// getPortfolioTemplates returns portfolio theme templates
func (tm *ThemeManager) getPortfolioTemplates() map[string]string {
	return map[string]string{
		"layouts/_default/single.html": `<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Page.Title }} | {{ .Site.Title }}</title>
    <meta name="description" content="{{ default .Site.Description .Page.Description }}">
    <link rel="stylesheet" href="{{ themeAsset "css/style.css" }}">
</head>
<body>
    <nav class="portfolio-nav">
        <a href="/" class="nav-logo">{{ .Site.Title }}</a>
        <ul class="nav-menu">
            <li><a href="/">Home</a></li>
            <li><a href="/projects/">Projects</a></li>
            <li><a href="/about/">About</a></li>
            <li><a href="/contact/">Contact</a></li>
        </ul>
    </nav>
    <main class="portfolio-main">
        <article class="project-detail">
            <header class="project-header">
                <h1 class="project-title">{{ .Page.Title }}</h1>
                {{ if .Page.Params.technologies }}
                <div class="project-tech">
                    {{ range .Page.Params.technologies }}
                        <span class="tech-tag">{{ . }}</span>
                    {{ end }}
                </div>
                {{ end }}
            </header>
            <div class="project-content">
                {{ .Page.Content }}
            </div>
            {{ if .Page.Params.demo_url }}
            <div class="project-links">
                <a href="{{ .Page.Params.demo_url }}" target="_blank" class="btn btn-primary">View Demo</a>
                {{ if .Page.Params.github_url }}
                <a href="{{ .Page.Params.github_url }}" target="_blank" class="btn btn-secondary">View Code</a>
                {{ end }}
            </div>
            {{ end }}
        </article>
    </main>
</body>
</html>`,
		"layouts/_default/list.html": `<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Site.Title }}</title>
    <meta name="description" content="{{ .Site.Description }}">
    <link rel="stylesheet" href="{{ themeAsset "css/style.css" }}">
</head>
<body>
    <nav class="portfolio-nav">
        <a href="/" class="nav-logo">{{ .Site.Title }}</a>
        <ul class="nav-menu">
            <li><a href="/">Home</a></li>
            <li><a href="/projects/">Projects</a></li>
            <li><a href="/about/">About</a></li>
            <li><a href="/contact/">Contact</a></li>
        </ul>
    </nav>
    <main class="portfolio-main">
        <section class="hero-section">
            <div class="hero-content">
                <h1 class="hero-title">{{ .Site.Title }}</h1>
                <p class="hero-subtitle">{{ .Site.Description }}</p>
            </div>
        </section>
        <section class="projects-section" id="projects">
            <h2>My Projects</h2>
            <div class="projects-grid">
                {{ range .Pages }}
                <article class="project-card">
                    <div class="project-image">
                        {{ if .Params.image }}
                        <img src="{{ .Params.image }}" alt="{{ .Title }}">
                        {{ else }}
                        <div class="project-placeholder">{{ substr .Title 0 1 }}</div>
                        {{ end }}
                    </div>
                    <div class="project-info">
                        <h3><a href="{{ .URL }}">{{ .Title }}</a></h3>
                        <p class="project-description">{{ .Summary }}</p>
                        {{ if .Params.technologies }}
                        <div class="project-tech">
                            {{ range .Params.technologies }}
                                <span class="tech-tag">{{ . }}</span>
                            {{ end }}
                        </div>
                        {{ end }}
                    </div>
                </article>
                {{ end }}
            </div>
        </section>
    </main>
</body>
</html>`,
	}
}

// getDocsTemplates returns documentation theme templates
func (tm *ThemeManager) getDocsTemplates() map[string]string {
	return map[string]string{
		"layouts/_default/single.html": `<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Page.Title }} | {{ .Site.Title }}</title>
    <meta name="description" content="{{ default .Site.Description .Page.Description }}">
    <link rel="stylesheet" href="{{ themeAsset "css/style.css" }}">
    {{ if hasFeature "syntax" }}
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
    <script>hljs.highlightAll();</script>
    {{ end }}
</head>
<body class="docs-layout">
    <nav class="docs-nav">
        <div class="nav-brand">
            <a href="/">{{ .Site.Title }}</a>
        </div>
        <div class="nav-search">
            <input type="search" placeholder="Search docs...">
        </div>
    </nav>
    <div class="docs-container">
        <aside class="docs-sidebar">
            <nav class="sidebar-nav">
                <h3>Navigation</h3>
                <ul>
                    <li><a href="/">Home</a></li>
                    <li><a href="/getting-started/">Getting Started</a></li>
                    <li><a href="/api/">API Reference</a></li>
                    <li><a href="/examples/">Examples</a></li>
                </ul>
            </nav>
        </aside>
        <main class="docs-main">
            <article class="docs-article">
                <header class="docs-header">
                    <h1>{{ .Page.Title }}</h1>
                    {{ if .Page.Description }}
                    <p class="docs-description">{{ .Page.Description }}</p>
                    {{ end }}
                </header>
                {{ if hasFeature "toc" }}
                <nav class="docs-toc">
                    <h4>Table of Contents</h4>
                    <!-- TOC would be generated here -->
                </nav>
                {{ end }}
                <div class="docs-content">
                    {{ .Page.Content }}
                </div>
                <footer class="docs-footer">
                    <div class="docs-navigation">
                        <!-- Previous/Next navigation would go here -->
                    </div>
                </footer>
            </article>
        </main>
    </div>
</body>
</html>`,
		"layouts/_default/list.html": `<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Site.Title }}</title>
    <meta name="description" content="{{ .Site.Description }}">
    <link rel="stylesheet" href="{{ themeAsset "css/style.css" }}">
</head>
<body class="docs-layout">
    <nav class="docs-nav">
        <div class="nav-brand">
            <a href="/">{{ .Site.Title }}</a>
        </div>
        <div class="nav-search">
            <input type="search" placeholder="Search docs...">
        </div>
    </nav>
    <div class="docs-container">
        <aside class="docs-sidebar">
            <nav class="sidebar-nav">
                <h3>Documentation</h3>
                <ul>
                    {{ range .Pages }}
                    <li><a href="{{ .URL }}">{{ .Title }}</a></li>
                    {{ end }}
                </ul>
            </nav>
        </aside>
        <main class="docs-main">
            <div class="docs-home">
                <header class="docs-hero">
                    <h1>{{ .Site.Title }}</h1>
                    <p>{{ .Site.Description }}</p>
                </header>
                <section class="docs-sections">
                    {{ range .Pages }}
                    <article class="docs-card">
                        <h2><a href="{{ .URL }}">{{ .Title }}</a></h2>
                        <p>{{ .Summary }}</p>
                    </article>
                    {{ end }}
                </section>
            </div>
        </main>
    </div>
</body>
</html>`,
	}

}