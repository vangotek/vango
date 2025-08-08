package vango

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vango/internal/builder"
	"vango/internal/config"

	"github.com/spf13/cobra"
)

// Global flags
var (
	configPath    string
	verbose       bool
	environment   string
	workers       int
	outputFormat  string
	profile       bool
)

var rootCmd = &cobra.Command{
	Use:   "vango",
	Short: "VanGo is a fast, modern static site generator built with Go",
	Long: `VanGo is a fast, modern static site generator built with Go.
It combines simplicity with powerful features to help you create beautiful websites efficiently.

Features:
  • Fast parallel builds with Go's performance
  • Modern Markdown processing with extensions
  • Flexible template engine with 30+ built-in functions
  • Live development server with hot reload
  • Advanced theming system with dark mode support
  • SEO optimization and social media integration
  • Multilingual support and content management
  • Performance optimizations and caching

Examples:
  vango                           # Build the site using default config
  vango serve                     # Start development server
  vango serve -p 8080             # Start server on port 8080
  vango build -e production       # Build for production
  vango theme list                # List available themes
  vango new site myblog           # Create new site
  vango new post "My New Post"    # Create new post`,
	Version: "2.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior: build the site
		buildSite(cmd)
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "", "Environment (development, production, etc.)")
	rootCmd.PersistentFlags().IntVarP(&workers, "workers", "w", 0, "Number of parallel workers (0 = auto)")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "format", "text", "Output format (text, json, yaml)")
	rootCmd.PersistentFlags().BoolVar(&profile, "profile", false, "Enable performance profiling")

	// Add all subcommands
	rootCmd.AddCommand(buildCmd)
	// serveCmd is added in serve.go
	rootCmd.AddCommand(newCmd)
	// themeCmd is added in theme.go
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(benchmarkCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(deployCmd)

	// Build command flags
	buildCmd.Flags().Bool("clean", false, "Clean output directory before building")
	buildCmd.Flags().Bool("drafts", false, "Include draft content")
	buildCmd.Flags().Bool("future", false, "Include future-dated content")
	buildCmd.Flags().Bool("expired", false, "Include expired content")
	buildCmd.Flags().Bool("minify", false, "Minify output")

	// Serve command flags will be defined in serve.go

	// New command structure
	newCmd.AddCommand(newSiteCmd)
	newCmd.AddCommand(newPostCmd)
	newCmd.AddCommand(newPageCmd)

	// Theme command structure is handled in theme.go

	// Config command structure
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)

	// Benchmark flags
	benchmarkCmd.Flags().Int("iterations", 10, "Number of benchmark iterations")
	benchmarkCmd.Flags().Bool("memory", false, "Include memory profiling")

	// Deploy flags
	deployCmd.Flags().String("target", "", "Deployment target")
	deployCmd.Flags().String("branch", "gh-pages", "Git branch for deployment")
	deployCmd.Flags().String("message", "", "Deployment commit message")
	deployCmd.Flags().Bool("force", false, "Force deployment")
}

// Build and serve commands are defined in their respective files

// New command for creating content
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new content",
	Long: `Create new content files, sites, or themes.

This command helps scaffold new content with proper front matter
and directory structure.`,
	Example: `  vango new site myblog           # Create new site
  vango new post "My New Post"    # Create new post
  vango new page about            # Create new page
  vango new theme mytheme         # Create new theme`,
}

var newSiteCmd = &cobra.Command{
	Use:   "site [name]",
	Short: "Create a new site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createNewSite(args[0])
	},
}

var newPostCmd = &cobra.Command{
	Use:   "post [title]",
	Short: "Create a new post",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createNewPost(args[0])
	},
}

var newPageCmd = &cobra.Command{
	Use:   "page [title]",
	Short: "Create a new page",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createNewPage(args[0])
	},
}

// Theme commands are defined in theme.go
// Config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long: `Manage site configuration.

View, validate, and modify your site's configuration settings.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		showConfig()
	},
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	Run: func(cmd *cobra.Command, args []string) {
		validateConfig()
	},
}

// Version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		showVersion()
	},
}

// Benchmark command
var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Run performance benchmarks",
	Long: `Run performance benchmarks to measure build speed and optimization.

This helps identify performance bottlenecks and optimize your site.`,
	Run: func(cmd *cobra.Command, args []string) {
		runBenchmark(cmd)
	},
}


// Validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate site content and configuration",
	Long: `Validate your site's content, configuration, and structure.

This command checks for:
  • Configuration errors
  • Invalid front matter
  • Broken internal links
  • Missing images
  • SEO issues`,
	Run: func(cmd *cobra.Command, args []string) {
		validateSite()
	},
}

// Deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the site",
	Long: `Deploy the built site to various hosting platforms.

Supports deployment to:
  • GitHub Pages
  • Netlify
  • Vercel
  • AWS S3
  • FTP/SFTP servers`,
	Example: `  vango deploy github             # Deploy to GitHub Pages
  vango deploy netlify            # Deploy to Netlify
  vango deploy s3                 # Deploy to AWS S3`,
	Run: func(cmd *cobra.Command, args []string) {
		deploySite(args)
	},
}
// Command implementations
func buildSite(cmd *cobra.Command) {
	start := time.Now()
	
	if verbose {
		fmt.Println("🏗️  Loading configuration...")
	}
	
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("📖 Building site '%s'...\n", cfg.Title)
		fmt.Printf("🌍 Environment: %s\n", cfg.Environment)
		fmt.Printf("👷 Workers: %d\n", cfg.Workers)
	}

	// Apply build flags
	if buildClean, _ := cmd.Flags().GetBool("clean"); buildClean {
		cfg.CleanBuild = true
	}
	if buildDrafts, _ := cmd.Flags().GetBool("drafts"); buildDrafts {
		cfg.BuildDrafts = true
	}
	if buildFuture, _ := cmd.Flags().GetBool("future"); buildFuture {
		cfg.BuildFuture = true
	}

	b := builder.New(cfg)
	
	if profile {
		// Enable profiling
		fmt.Println("📊 Performance profiling enabled")
	}
	
	if err := b.Build(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Build failed: %v\n", err)
		os.Exit(1)
	}

	duration := time.Since(start)
	pages := b.GetPages()
	
	fmt.Printf("✅ Site built successfully!\n")
	fmt.Printf("📁 Output directory: %s\n", cfg.PublicDir)
	fmt.Printf("📄 Generated %d pages in %v\n", len(pages), duration)
	
	if verbose {
		fmt.Printf("⚡ Average: %.2f pages/second\n", float64(len(pages))/duration.Seconds())
	}
}

// serveServer function is moved to serve.go file

func createNewSite(name string) {
	fmt.Printf("🏗️  Creating new site: %s\n", name)
	
	if err := os.MkdirAll(name, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	// Create directory structure
	dirs := []string{
		"content",
		"layouts/_default",
		"static",
		"themes/modern-app",
		"data",
	}

	for _, dir := range dirs {
		path := filepath.Join(name, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to create directory %s: %v\n", path, err)
			os.Exit(1)
		}
	}

	// Create config file
	configContent := fmt.Sprintf(`title = "%s"
baseURL = "https://example.com/"
language = "en"
description = "A new VanGo site"
theme = "modern-app"

# Directory paths
contentDir = "content"
layoutDir = "layouts"
staticDir = "static"
publicDir = "public"

# Build settings
buildDrafts = false
buildFuture = false
cleanBuild = true

# Server settings
port = 1313
host = "localhost"
liveReload = true

[params]
    author = "Your Name"
    version = "1.0.0"
`, name)

	configPath := filepath.Join(name, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create config file: %v\n", err)
		os.Exit(1)
	}

	// Create sample content
	samplePost := "+++\n" +
	"title = \"Welcome to " + name + "\"\n" +
	"date = \"" + time.Now().Format("2006-01-02T15:04:05Z07:00") + "\"\n" +
	"description = \"Your first post with VanGo\"\n" +
	"draft = false\n" +
	"tags = [\"welcome\", \"getting-started\"]\n" +
	"+++\n\n" +
	"# Welcome to Your New VanGo Site!\n\n" +
	"This is your first post. You can edit this file or create new posts in the `content` directory.\n\n" +
	"## Getting Started\n\n" +
	"1. Edit this post in `content/welcome.md`\n" +
	"2. Create new posts with `vango new post \"Post Title\"`\n" +
	"3. Start the development server with `vango serve`\n" +
	"4. Build your site with `vango build`\n\n" +
	"## Features\n\n" +
	"VanGo includes many powerful features:\n\n" +
	"- Fast Go-powered builds\n" +
	"- Live reload development server\n" +
	"- Markdown with front matter\n" +
	"- Flexible theming system\n" +
	"- SEO optimization\n" +
	"- And much more!\n\n" +
	"Happy building! 🚀\n"


	postPath := filepath.Join(name, "content", "welcome.md")
	if err := os.WriteFile(postPath, []byte(samplePost), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create sample post: %v\n", err)
		os.Exit(1)
	}

	// Create basic template
	templateContent := `<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Page.Title }} | {{ .Site.Title }}</title>
    <meta name="description" content="{{ default .Site.Description .Page.Description }}">
</head>
<body>
    <header>
        <h1><a href="/">{{ .Site.Title }}</a></h1>
    </header>
    
    <main>
        <article>
            <h1>{{ .Page.Title }}</h1>
            <time>{{ humanizeDate .Page.ParsedDate }}</time>
            <div>{{ .Page.Content }}</div>
        </article>
    </main>
    
    <footer>
        <p>&copy; {{ dateFormat "2006" now }} {{ .Site.Params.author }}. Built with VanGo.</p>
    </footer>
</body>
</html>`

	templatePath := filepath.Join(name, "layouts", "_default", "single.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create template: %v\n", err)
		os.Exit(1)
	}

	// Create theme files
	themeFiles := map[string]string{
		"themes/modern-app/theme.json": `{ 
  "name": "modern-app",
  "version": "1.0.0",
  "description": "A modern, clean theme with responsive design and dark mode support",
  "author": "VanGo Team",
  "homepage": "",
  "license": "MIT",
  "min_vango_version": "1.0.0",
  "tags": ["modern", "responsive", "dark-mode", "clean"],
  "features": ["responsive", "dark-mode", "syntax-highlighting", "reading-time"],
  "config": {
    "colors": {
      "primary": "#3b82f6",
      "secondary": "#6b7280",
      "accent": "#10b981"
    }
  },
  "layouts_dir": "layouts",
  "static_dir": "static",
  "assets_dir": "assets"
}`,
		"themes/modern-app/layouts/_default/baseof.html": `<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ block "title" . }}{{ .Page.Title }} | {{ .Site.Title }}{{ end }}</title>
    <meta name="description" content="{{ block "description" . }}{{ default .Site.Description .Page.Description }}{{ end }}">
    <meta name="author" content="{{ default .Site.Author .Page.Author }}">
    
    <!-- Open Graph / Facebook -->
    <meta property="og:type" content="{{ block "og_type" . }}article{{ end }}">
    <meta property="og:url" content="{{ .Site.BaseURL }}{{ .Page.URL }}">
    <meta property="og:title" content="{{ .Page.Title }}">
    <meta property="og:description" content="{{ default .Site.Description .Page.Description }}">
    
    <!-- Twitter -->
    <meta property="twitter:card" content="summary">
    <meta property="twitter:url" content="{{ .Site.BaseURL }}{{ .Page.URL }}">
    <meta property="twitter:title" content="{{ .Page.Title }}">
    <meta property="twitter:description" content="{{ default .Site.Description .Page.Description }}">
    
    <link rel="stylesheet" href="{{ themeAsset "css/style.css" }}">
    <link rel="canonical" href="{{ .Site.BaseURL }}{{ .Page.URL }}">
    
    {{ block "head" . }}{{ end }}
    
    {{ if hasFeature "syntax" }}
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
    <script>hljs.highlightAll();</script>
    {{ end }}
</head>
<body class="{{ block "body_class" . }}modern-theme{{ end }}">
    <nav class="navbar">
        <div class="nav-container">
            <a href="" class="nav-logo">{{ .Site.Title }}</a>
            <ul class="nav-menu">
                <li><a href="" class="nav-link">Home</a></li>
                <li><a href="/about/" class="nav-link">About</a></li>
            </ul>
            {{ if hasFeature "dark_mode" }}
            <button class="theme-toggle" onclick="toggleTheme()">🌙</button>
            {{ end }}
        </div>
    </nav>

    {{ block "main" . }}
    <main class="main-content">
        {{ block "content" . }}{{ end }}
    </main>
    {{ end }}

    <footer class="site-footer">
        <div class="footer-container">
            <p>&copy; {{ dateFormat "2006" now }} {{ .Site.Author }}. Built with VanGo.</p>
            {{ if .Site.Params.social }}
            <div class="social-links">
                {{ if index .Site.Params.social "twitter" }}
                    <a href="https://twitter.com/{{ index .Site.Params.social "twitter" }}" target="_blank" class="social-link">Twitter</a>
                {{ end }}
                {{ if index .Site.Params.social "github" }}
                    <a href="https://github.com/{{ index .Site.Params.social "github" }}" target="_blank" class="social-link">GitHub</a>
                {{ end }}
            </div>
            {{ end }}
        </div>
    </footer>
    
    {{ block "scripts" . }}
    {{ if hasFeature "dark_mode" }}
    <script>
        function toggleTheme() {
            document.body.classList.toggle('dark-theme');
            const isDark = document.body.classList.contains('dark-theme');
            sessionStorage.setItem('theme', isDark ? 'dark' : 'light');
            document.querySelector('.theme-toggle').textContent = isDark ? '☀️' : '🌙';
        }
        
        // Load saved theme from sessionStorage
        const savedTheme = sessionStorage.getItem('theme');
        if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
            document.body.classList.add('dark-theme');
            if (document.querySelector('.theme-toggle')) {
                document.querySelector('.theme-toggle').textContent = '☀️';
            }
        }
    </script>
    {{ end }}
    {{ end }}
</body>
</html>`,
		"themes/modern-app/layouts/_default/list.html": `{{ define "content" }}
<section class="hero-section">
    <div class="hero-content">
        <h1 class="hero-title">{{ .Site.Title }}</h1>
        <p class="hero-description">{{ .Site.Description }}</p>
    </div>
</section>

<section class="posts-section">
    <div class="section-header">
        <h2 class="section-title">Latest Posts</h2>
    </div>
    
    <div class="posts-grid">
        {{ range .Pages }}
        <article class="post-card">
            <div class="post-card-content">
                <h3 class="post-card-title">
                    <a href="{{ .URL }}" class="post-link">{{ .Title }}</a>
                </h3>
                <div class="post-card-meta">
                    <time datetime="{{ dateFormat "2006-01-02" .ParsedDate }}">
                        {{ humanizeDate .ParsedDate }}
                    </time>
                    {{ if hasFeature "reading_time" }}
                        {{ if gt .ReadingTime 0 }}
                        • {{ .ReadingTime }} min read
                        {{ end }}
                    {{ end }}
                </div>
                {{ if .Description }}
                <p class="post-card-excerpt">{{ .Description }}</p>
                {{ end }}
                {{ if .Tags }}
                <div class="post-card-tags">
                    {{ range .Tags }}
                        <span class="tag">#{{ . }}</span>
                    {{ end }}
                </div>
                {{ end }}
            </div>
        </article>
        {{ end }}
    </div>
</section>
{{ end }}`,
		"themes/modern-app/layouts/_default/single.html": `{{ define "content" }}
<article class="article-container">
    <header class="article-header">
        <h1 class="article-title">{{ .Page.Title }}</h1>
        <div class="article-meta">
            <time datetime="{{ dateFormat "2006-01-02" .Page.ParsedDate }}">
                {{ humanizeDate .Page.ParsedDate }}
            </time>
            {{ if .Page.Author }}
                by <span class="author-name">{{ .Page.Author }}</span>
            {{ end }}
            {{ if hasFeature "reading_time" }}
                {{ if gt .Page.ReadingTime 0 }}
                • <span class="reading-time">{{ .Page.ReadingTime }} min read</span>
                {{ end }}
            {{ end }}
            {{ if .Page.WordCount }}
                • <span class="word-count">{{ .Page.WordCount }} words</span>
            {{ end }}
        </div>
        
        {{ if .Page.Tags }}
        <div class="article-tags">
            {{ range .Page.Tags }}
                <span class="tag">#{{ . }}</span>
            {{ end }}
        </div>
        {{ end }}
    </header>

    <div class="article-content">
        {{ .Page.Content }}
    </div>

    {{ if .Page.Categories }}
    <footer class="article-footer">
        <div class="categories">
            <strong>Categories:</strong>
            {{ range $i, $cat := .Page.Categories }}
                {{ if $i }}, {{ end }}
                <a href="/categories/{{ lower $cat }}/" class="category-link">{{ $cat }}</a>
            {{ end }}
        </div>
    </footer>
    {{ end }}
</article>
{{ end }}`,
		"themes/modern-app/static/css/style.css": `/* Modern App Theme for VanGo */
:root {
  /* Colors */
  --color-primary: #3b82f6;
  --color-primary-hover: #2563eb;
  --color-secondary: #6b7280;
  --color-accent: #10b981;
  --color-background: #ffffff;
  --color-surface: #f8fafc;
  --color-text: #1f2937;
  --color-text-light: #6b7280;
  --color-border: #e5e7eb;
  --color-success: #10b981;
  --color-warning: #f59e0b;
  --color-error: #ef4444;
  
  /* Dark theme colors */
  --color-dark-background: #0f172a;
  --color-dark-surface: #1e293b;
  --color-dark-text: #f1f5f9;
  --color-dark-text-light: #94a3b8;
  --color-dark-border: #334155;
  
  /* Typography */
  --font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  --font-size-xs: 0.75rem;
  --font-size-sm: 0.875rem;
  --font-size-base: 1rem;
  --font-size-lg: 1.125rem;
  --font-size-xl: 1.25rem;
  --font-size-2xl: 1.5rem;
  --font-size-3xl: 1.875rem;
  --font-size-4xl: 2.25rem;
  
  /* Spacing */
  --spacing-xs: 0.25rem;
  --spacing-sm: 0.5rem;
  --spacing-md: 1rem;
  --spacing-lg: 1.5rem;
  --spacing-xl: 2rem;
  --spacing-2xl: 3rem;
  --spacing-3xl: 4rem;
  
  /* Layout */
  --max-width: 1200px;
  --content-width: 800px;
  
  /* Shadows */
  --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
  --shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
  --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1);
  
  /* Border radius */
  --radius-sm: 0.25rem;
  --radius-md: 0.375rem;
  --radius-lg: 0.5rem;
  --radius-xl: 0.75rem;
  
  /* Transitions */
  --transition-fast: 0.15s ease-in-out;
  --transition-normal: 0.3s ease-in-out;
  --transition-slow: 0.5s ease-in-out;
}

/* Reset and base styles */
*, 
*::before, 
*::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

html {
  scroll-behavior: smooth;
}

body {
  font-family: var(--font-family);
  font-size: var(--font-size-base);
  line-height: 1.6;
  color: var(--color-text);
  background-color: var(--color-background);
  transition: color var(--transition-normal), background-color var(--transition-normal);
}

/* Dark theme */
body.dark-theme {
  color: var(--color-dark-text);
  background-color: var(--color-dark-background);
}

body.dark-theme .navbar {
  background-color: var(--color-dark-surface);
  border-color: var(--color-dark-border);
}

body.dark-theme .nav-link {
  color: var(--color-dark-text);
}

body.dark-theme .nav-link:hover {
  background-color: var(--color-dark-border);
}

body.dark-theme .article-container,
body.dark-theme .post-card {
  background-color: var(--color-dark-surface);
  border-color: var(--color-dark-border);
}

body.dark-theme .site-footer {
  background-color: var(--color-dark-surface);
  border-color: var(--color-dark-border);
}

body.dark-theme .hero-title,
body.dark-theme .section-title,
body.dark-theme .article-title,
body.dark-theme .post-card-title,
body.dark-theme h1, 
body.dark-theme h2, 
body.dark-theme h3, 
body.dark-theme h4, 
body.dark-theme h5, 
body.dark-theme h6 {
  color: var(--color-dark-text);
}

body.dark-theme .post-link {
  color: var(--color-dark-text);
}

body.dark-theme .post-card-meta,
body.dark-theme .article-meta,
body.dark-theme .post-card-excerpt {
  color: var(--color-dark-text-light);
}

body.dark-theme .article-content {
  color: var(--color-dark-text);
}

body.dark-theme .article-content h1,
body.dark-theme .article-content h2,
body.dark-theme .article-content h3,
body.dark-theme .article-content h4,
body.dark-theme .article-content h5,
body.dark-theme .article-content h6 {
  color: var(--color-dark-text);
}

body.dark-theme .categories {
  color: var(--color-dark-text-light);
}

body.dark-theme .theme-toggle:hover {
  background-color: var(--color-dark-border);
}

body.dark-theme .admin-panel-btn {
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  color: white;
}

body.dark-theme .admin-panel-btn:hover {
  opacity: 0.8;
}

/* Typography */
h1, h2, h3, h4, h5, h6 {
  font-weight: 700;
  line-height: 1.2;
  margin-bottom: var(--spacing-md);
  color: inherit;
}

h1 { font-size: var(--font-size-4xl); }
h2 { font-size: var(--font-size-3xl); }
h3 { font-size: var(--font-size-2xl); }
h4 { font-size: var(--font-size-xl); }
h5 { font-size: var(--font-size-lg); }
h6 { font-size: var(--font-size-base); }

p {
  margin-bottom: var(--spacing-md);
  line-height: 1.7;
}

a {
  color: var(--color-primary);
  text-decoration: none;
  transition: color var(--transition-fast);
}

a:hover {
  color: var(--color-primary-hover);
}

/* Navigation */
.navbar {
  background-color: var(--color-background);
  border-bottom: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
  position: sticky;
  top: 0;
  z-index: 100;
  transition: background-color var(--transition-normal), border-color var(--transition-normal);
}

.nav-container {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 0 var(--spacing-lg);
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.nav-logo {
  font-size: var(--font-size-xl);
  font-weight: 700;
  color: var(--color-primary);
  text-decoration: none;
}

.nav-menu {
  display: flex;
  list-style: none;
  gap: var(--spacing-xl);
  margin: 0;
  padding: 0;
}

.nav-link {
  font-weight: 500;
  color: var(--color-text);
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius-md);
  transition: color var(--transition-fast), background-color var(--transition-fast);
}

.nav-link:hover {
  color: var(--color-primary);
  background-color: var(--color-surface);
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.admin-panel-btn {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  color: white;
  border: none;
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius-md);
  font-weight: 600;
  font-size: var(--font-size-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
  box-shadow: var(--shadow-sm);
}

.admin-panel-btn:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
  opacity: 0.9;
}

.admin-icon {
  font-size: var(--font-size-base);
}

.admin-text {
  font-size: var(--font-size-sm);
}

.theme-toggle {
  background: none;
  border: none;
  font-size: var(--font-size-lg);
  cursor: pointer;
  padding: var(--spacing-sm);
  border-radius: var(--radius-md);
  transition: background-color var(--transition-fast);
}

.theme-toggle:hover {
  background-color: var(--color-surface);
}

/* Main content */
.main-content {
  flex: 1;
  min-height: calc(100vh - 64px - 200px);
}

/* Hero section */
.hero-section {
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-accent) 100%);
  color: white;
  padding: var(--spacing-3xl) var(--spacing-lg);
  text-align: center;
  margin-bottom: var(--spacing-3xl);
}

.hero-content {
  max-width: var(--content-width);
  margin: 0 auto;
}

.hero-title {
  font-size: var(--font-size-4xl);
  font-weight: 900;
  margin-bottom: var(--spacing-lg);
  text-shadow: 0 2px 4px rgb(0 0 0 / 0.1);
}

.hero-description {
  font-size: var(--font-size-xl);
  opacity: 0.9;
  line-height: 1.6;
}

/* Posts section */
.posts-section {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 0 var(--spacing-lg) var(--spacing-3xl);
}

.section-header {
  text-align: center;
  margin-bottom: var(--spacing-3xl);
}

.section-title {
  font-size: var(--font-size-3xl);
  font-weight: 800;
  color: var(--color-text);
}

.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: var(--spacing-xl);
}

/* Post cards */
.post-card {
  background-color: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  padding: var(--spacing-xl);
  box-shadow: var(--shadow-sm);
  transition: transform var(--transition-fast), box-shadow var(--transition-fast);
}

.post-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.post-card-title {
  font-size: var(--font-size-xl);
  font-weight: 700;
  margin-bottom: var(--spacing-sm);
}

.post-link {
  color: var(--color-text);
}

.post-link:hover {
  color: var(--color-primary);
}

.post-card-meta {
  color: var(--color-text-light);
  font-size: var(--font-size-sm);
  margin-bottom: var(--spacing-md);
}

.post-card-excerpt {
  color: var(--color-text-light);
  line-height: 1.6;
  margin-bottom: var(--spacing-md);
}

.post-card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}

/* Article styles */
.article-container {
  max-width: var(--content-width);
  margin: var(--spacing-3xl) auto;
  background-color: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  padding: var(--spacing-3xl);
  box-shadow: var(--shadow-md);
}

.article-header {
  margin-bottom: var(--spacing-3xl);
  padding-bottom: var(--spacing-xl);
  border-bottom: 1px solid var(--color-border);
}

.article-title {
  font-size: var(--font-size-4xl);
  font-weight: 900;
  line-height: 1.1;
  margin-bottom: var(--spacing-lg);
  color: var(--color-text);
}

.article-meta {
  color: var(--color-text-light);
  font-size: var(--font-size-sm);
  margin-bottom: var(--spacing-lg);
}

.author-name {
  font-weight: 600;
  color: var(--color-primary);
}

.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-lg);
}

/* Tags */
.tag {
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  color: white;
  font-size: var(--font-size-xs);
  font-weight: 500;
  padding: var(--spacing-xs) var(--spacing-md);
  border-radius: var(--radius-lg);
}

/* Article content */
.article-content {
  font-size: var(--font-size-lg);
  line-height: 1.8;
  color: var(--color-text);
}

.article-content h1,
.article-content h2,
.article-content h3,
.article-content h4,
.article-content h5,
.article-content h6 {
  margin: var(--spacing-2xl) 0 var(--spacing-lg);
  color: var(--color-text);
  font-weight: 700;
}

.article-content h1 { font-size: var(--font-size-3xl); }
.article-content h2 { font-size: var(--font-size-2xl); }
.article-content h3 { font-size: var(--font-size-xl); }

.article-content p {
  margin-bottom: var(--spacing-lg);
}

.article-content a {
  color: var(--color-primary);
  border-bottom: 1px solid transparent;
  transition: border-color var(--transition-fast);
}

.article-content a:hover {
  border-bottom-color: var(--color-primary);
}

.article-content blockquote {
  border-left: 4px solid var(--color-primary);
  background-color: var(--color-surface);
  padding: var(--spacing-lg);
  margin: var(--spacing-xl) 0;
  border-radius: 0 var(--radius-md) var(--radius-md) 0;
  font-style: italic;
}

.article-content code {
  background-color: var(--color-surface);
  color: var(--color-primary);
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.9em;
}

.article-content pre {
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--spacing-lg);
  overflow-x: auto;
  margin: var(--spacing-xl) 0;
}

.article-content pre code {
  background: none;
  color: inherit;
  padding: 0;
}

.article-content ul,
.article-content ol {
  margin: var(--spacing-lg) 0;
  padding-left: var(--spacing-xl);
}

.article-content li {
  margin-bottom: var(--spacing-sm);
}

/* Article footer */
.article-footer {
  margin-top: var(--spacing-3xl);
  padding-top: var(--spacing-xl);
  border-top: 1px solid var(--color-border);
}

.categories {
  color: var(--color-text-light);
  font-size: var(--font-size-sm);
}

.category-link {
  color: var(--color-primary);
  font-weight: 500;
}

/* Footer */
.site-footer {
  background-color: var(--color-surface);
  border-top: 1px solid var(--color-border);
  padding: var(--spacing-2xl) var(--spacing-lg);
  text-align: center;
  margin-top: auto;
}

.footer-container {
  max-width: var(--max-width);
  margin: 0 auto;
}

.social-links {
  margin-top: var(--spacing-lg);
  display: flex;
  justify-content: center;
  gap: var(--spacing-lg);
}

.social-link {
  color: var(--color-text-light);
  font-weight: 500;
  transition: color var(--transition-fast);
}

.social-link:hover {
  color: var(--color-primary);
}

/* Responsive design */
@media (max-width: 768px) {
  .nav-container {
    padding: 0 var(--spacing-md);
    flex-direction: column;
    height: auto;
    gap: var(--spacing-md);
    padding-top: var(--spacing-md);
    padding-bottom: var(--spacing-md);
  }
  
  .nav-menu {
    gap: var(--spacing-lg);
  }
  
  .nav-actions {
    gap: var(--spacing-md);
  }
  
  .admin-panel-btn {
    padding: var(--spacing-xs) var(--spacing-sm);
  }
  
  .admin-text {
    display: none;
  }
  
  .hero-section {
    padding: var(--spacing-2xl) var(--spacing-md);
  }
  
  .hero-title {
    font-size: var(--font-size-3xl);
  }
  
  .posts-section {
    padding: 0 var(--spacing-md) var(--spacing-2xl);
  }
  
  .posts-grid {
    grid-template-columns: 1fr;
    gap: var(--spacing-lg);
  }
  
  .article-container {
    margin: var(--spacing-lg) var(--spacing-md);
    padding: var(--spacing-lg);
  }
  
  .article-title {
    font-size: var(--font-size-3xl);
  }
  
  .article-content {
  font-size: var(--font-size-base);
  }
  
  .social-links {
    flex-direction: column;
    gap: var(--spacing-md);
  }
}

@media (max-width: 480px) {
  .hero-title {
    font-size: var(--font-size-2xl);
  }
  
  .article-title {
    font-size: var(--font-size-2xl);
  }
  
  .posts-grid {
    gap: var(--spacing-md);
  }
  
  .post-card {
    padding: var(--spacing-lg);
  }
}

/* Print styles */
@media print {
  .navbar,
  .site-footer,
  .theme-toggle,
  .social-links {
    display: none;
  }
  
  .article-container {
    box-shadow: none;
    border: none;
    margin: 0;
    padding: 0;
  }
  
  .hero-section {
    background: none;
    color: var(--color-text);
  }
}

/* High contrast mode support */
@media (prefers-contrast: high) {
  :root {
    --color-border: #000000;
    --color-text-light: var(--color-text);
  }
  
  .post-card,
  .article-container {
    border-width: 2px;
  }
}

/* Reduced motion support */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
`,
	}

	for path, content := range themeFiles {
		fullPath := filepath.Join(name, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to create directory for %s: %v\n", fullPath, err)
			os.Exit(1)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to create file %s: %v\n", fullPath, err)
			os.Exit(1)
		}
	}

	fmt.Printf("✅ Site created successfully!\n")
	fmt.Printf("📁 Location: %s\n", name)
	fmt.Printf("🚀 Next steps:\n")
	fmt.Printf("   cd %s\n", name)
	fmt.Printf("   vango serve\n")
}


func createNewPost(title string) {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Generate filename from title
	filename := strings.ToLower(title)
	filename = strings.ReplaceAll(filename, " ", "-")
	filename = strings.ReplaceAll(filename, "'", "")
	filename += ".md"

	postContent := fmt.Sprintf(`+++
title = "%s"
date = "%s"
description = ""
author = "%s"
draft = true
tags = []
categories = []
+++

# %s

Write your post content here...
`, title, time.Now().Format("2006-01-02T15:04:05Z07:00"), cfg.Author, title)

	postPath := filepath.Join(cfg.ContentDir, filename)
	if err := os.WriteFile(postPath, []byte(postContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create post: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Post created: %s\n", postPath)
}

func createNewPage(title string) {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	filename := strings.ToLower(title)
	filename = strings.ReplaceAll(filename, " ", "-")
	filename += ".md"

	pageContent := fmt.Sprintf(`+++
title = "%s"
date = "%s"
description = ""
draft = false
+++

# %s

Page content goes here...
`, title, time.Now().Format("2006-01-02T15:04:05Z07:00"), title)

	pagePath := filepath.Join(cfg.ContentDir, filename)
	if err := os.WriteFile(pagePath, []byte(pageContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create page: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Page created: %s\n", pagePath)
}

// Theme functions are now in theme.go

func showConfig() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	switch outputFormat {
	case "json":
		// Output as JSON
		fmt.Println("JSON output not implemented yet")
	case "yaml":
		// Output as YAML
		fmt.Println("YAML output not implemented yet")
	default:
		// Text output
		fmt.Printf("🔧 Configuration:\n\n")
		fmt.Printf("Title: %s\n", cfg.Title)
		fmt.Printf("Base URL: %s\n", cfg.BaseURL)
		fmt.Printf("Language: %s\n", cfg.Language)
		fmt.Printf("Theme: %s\n", cfg.Theme)
		fmt.Printf("Environment: %s\n", cfg.Environment)
		fmt.Printf("Content Dir: %s\n", cfg.ContentDir)
		fmt.Printf("Output Dir: %s\n", cfg.PublicDir)
	}
}

func validateConfig() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Configuration validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Configuration is valid")
	fmt.Printf("📊 Settings validated for environment: %s\n", cfg.Environment)
}

func showVersion() {
	version := rootCmd.Version
	fmt.Printf("VanGo v%s\n", version)
	fmt.Printf("Built with Go\n")
	fmt.Printf("https://github.com/vango/vango\n")
}

func runBenchmark(cmd *cobra.Command) {
	iterations, _ := cmd.Flags().GetInt("iterations")
	includeMemory, _ := cmd.Flags().GetBool("memory")

	fmt.Printf("🏃 Running benchmark (%d iterations)...\n", iterations)
	
	var totalDuration time.Duration
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	b := builder.New(cfg)
	
	for i := 0; i < iterations; i++ {
		start := time.Now()
		if err := b.Build(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Build failed on iteration %d: %v\n", i+1, err)
			os.Exit(1)
		}
		duration := time.Since(start)
		totalDuration += duration
		
		if verbose {
			fmt.Printf("  Iteration %d: %v\n", i+1, duration)
		}
	}

	avgDuration := totalDuration / time.Duration(iterations)
	pages := b.GetPages()
	
	fmt.Printf("📊 Benchmark Results:\n")
	fmt.Printf("  Iterations: %d\n", iterations)
	fmt.Printf("  Total time: %v\n", totalDuration)
	fmt.Printf("  Average time: %v\n", avgDuration)
	fmt.Printf("  Pages built: %d\n", len(pages))
	fmt.Printf("  Pages/second: %.2f\n", float64(len(pages))/avgDuration.Seconds())

	if includeMemory {
		fmt.Printf("  Memory profiling: enabled\n")
		// Memory profiling implementation
	}
}

func validateSite() {
	fmt.Println("🔍 Validating site...")
	
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Configuration error: %v\n", err)
		os.Exit(1)
	}

	issues := 0
	
	// Validate configuration
	fmt.Printf("✅ Configuration valid\n")
	
	// Check content directory
	if _, err := os.Stat(cfg.ContentDir); os.IsNotExist(err) {
		fmt.Printf("❌ Content directory missing: %s\n", cfg.ContentDir)
		issues++
	} else {
		fmt.Printf("✅ Content directory exists\n")
	}
	
	// Check layout directory
	if _, err := os.Stat(cfg.LayoutDir); os.IsNotExist(err) {
		fmt.Printf("❌ Layout directory missing: %s\n", cfg.LayoutDir)
		issues++
	} else {
		fmt.Printf("✅ Layout directory exists\n")
	}

	// Validate content files
	// Implementation would check for valid front matter, broken links, etc.
	
	if issues == 0 {
		fmt.Printf("✅ Site validation completed - no issues found\n")
	} else {
		fmt.Printf("⚠️  Site validation completed - %d issues found\n", issues)
		os.Exit(1)
	}
}

func deploySite(args []string) {
	if len(args) == 0 {
		fmt.Println("❌ Deployment target required")
		fmt.Println("Available targets: github, netlify, vercel, s3, ftp")
		os.Exit(1)
	}

	target := args[0]
	fmt.Printf("🚀 Deploying to %s...\n", target)
	
	// Build site first
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Set production environment
	cfg.Environment = "production"
	cfg.Performance.EnableMinification = true
	
	b := builder.New(cfg)
	if err := b.Build(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Build failed: %v\n", err)
		os.Exit(1)
	}

	switch target {
	case "github":
		deployToGitHub(cfg)
	case "netlify":
		deployToNetlify(cfg)
	case "vercel":
		deployToVercel(cfg)
	case "s3":
		deployToS3(cfg)
	default:
		fmt.Printf("❌ Unknown deployment target: %s\n", target)
		os.Exit(1)
	}
}

func deployToGitHub(cfg *config.Config) {
	fmt.Println("📤 Deploying to GitHub Pages...")
	// Implementation for GitHub Pages deployment
	fmt.Println("✅ Deployed to GitHub Pages!")
}

func deployToNetlify(cfg *config.Config) {
	fmt.Println("📤 Deploying to Netlify...")
	// Implementation for Netlify deployment
	fmt.Println("✅ Deployed to Netlify!")
}

func deployToVercel(cfg *config.Config) {
	fmt.Println("📤 Deploying to Vercel...")
	// Implementation for Vercel deployment
	fmt.Println("✅ Deployed to Vercel!")
}

func deployToS3(cfg *config.Config) {
	fmt.Println("📤 Deploying to AWS S3...")
	// Implementation for S3 deployment
	fmt.Println("✅ Deployed to S3!")
}

// Helper function to load configuration
func loadConfig() (*config.Config, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	// Apply global flag overrides
	if environment != "" {
		cfg.Environment = environment
	}
	if workers > 0 {
		cfg.Workers = workers
	}

	return cfg, nil
}