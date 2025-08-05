# VanGo - Static Site Generator

A fast, modern static site generator built with Go. VanGo combines simplicity with powerful features to help you create beautiful websites efficiently.

## Features

- 🚀 **Fast builds** - Powered by Go's performance
- 📝 **Markdown support** - Write content in Markdown with TOML front matter
- 🎨 **Template engine** - Flexible HTML templates with built-in functions
- 🔧 **Development server** - Live preview with hot reloading
- 📁 **Static asset handling** - Automatic copying and optimization
- 🏗️ **Modular architecture** - Clean, extensible codebase
- 🌙 **Dark mode support** - Built-in responsive design
- 📱 **Mobile-first** - Responsive templates out of the box

## Quick Start

### Prerequisites

- Go 1.22 or later

### Installation

1. Clone or download the VanGo project
2. Navigate to the project directory
3. Install dependencies:

```bash
go mod tidy
```

### Usage

#### Build your site
```bash
go run main.go
```

#### Start development server
```bash
go run main.go -mode serve
```

#### Custom port
```bash
go run main.go -mode serve -port 8080
```

#### Custom config file
```bash
go run main.go -config my-config.toml
```

#### Help
```bash
go run main.go -help
```

## Directory Structure

```
your-site/
├── config.toml          # Site configuration
├── content/             # Your content files (Markdown)
│   ├── hello.md
│   └── about.md
├── layouts/             # HTML templates
│   └── _default/
│       └── single.html
├── static/              # Static assets (CSS, JS, images)
│   └── style.css
├── public/              # Generated site (output)
├── internal/            # Go packages
│   ├── config/         # Configuration management
│   ├── content/        # Content parsing
│   ├── template/       # Template rendering
│   ├── builder/        # Site building
│   └── server/         # Development server
└── main.go             # Entry point
```

## Configuration

The `config.toml` file controls your site's behavior:

```toml
title = "My Awesome Site"
description = "A fantastic website built with VanGo"
baseURL = "https://mysite.com/"
language = "en"
author = "Your Name"

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

# Custom parameters
[params]
    version = "1.0.0"
    [params.social]
        twitter = "@username"
        github = "username"
```

## Content Format

Content files use Markdown with TOML front matter:

```markdown
+++
title = "My Post Title"
date = "2025-08-05"
description = "A brief description of the post"
author = "Author Name"
tags = ["tag1", "tag2"]
categories = ["Category"]
draft = false
+++

# Your Content Here

This is the content of your post written in **Markdown**.

## Features

- List items
- *Italic text*
- **Bold text**
- `code snippets`

```code
Code blocks are supported too!
```
```

## Templates

VanGo uses Go's `html/template` package with many built-in functions:

### Template Functions

- `{{ .Site.Title }}` - Site title
- `{{ .Page.Title }}` - Page title
- `{{ .Page.Content }}` - Rendered content
- `{{ .Page.Date }}` - Page date
- `{{ .Page.ReadingTime }}` - Calculated reading time
- `{{ .Page.WordCount }}` - Word count
- `{{ dateFormat "2006-01-02" .Page.Date }}` - Format dates
- `{{ humanizeDate .Page.Date }}` - Human-readable dates
- `{{ timeAgo .Page.Date }}` - Time since publication
- `{{ range .Page.Tags }}` - Loop through tags
- `{{ upper .Page.Title }}` - String manipulation
- `{{ default "default" .Page.Author }}` - Default values

### Template Structure

```html
<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
    <title>{{ .Page.Title }} | {{ .Site.Title }}</title>
    <meta name="description" content="{{ .Page.Description }}">
</head>
<body>
    <h1>{{ .Page.Title }}</h1>
    <time>{{ humanizeDate .Page.Date }}</time>
    <div>{{ .Page.Content }}</div>
</body>
</html>
```

## Development Server

The development server provides:

- Live preview at `http://localhost:1313`
- Automatic rebuilding on file changes
- API endpoints for debugging:
  - `/api/status` - Server status and statistics
  - `/api/rebuild` - Manual rebuild trigger
- Custom 404 page support
- Static file serving

## Architecture

VanGo is built with a modular architecture:

### Packages

- **config** - Configuration management and validation
- **content** - Markdown parsing and page structure
- **template** - HTML template rendering with functions
- **builder** - Site generation and static file handling
- **server** - Development server with live reload

### Key Components

1. **Config Parser** - Loads and validates TOML configuration
2. **Content Parser** - Processes Markdown files with front matter
3. **Template Engine** - Renders HTML with custom functions
4. **Site Builder** - Generates static site files
5. **Dev Server** - Provides live preview capabilities

## Building from Source

```bash
# Clone the repository
git clone <repository-url> vango
cd vango

# Install dependencies
go mod tidy

# Build binary
go build -o vango main.go

# Run
./vango -mode serve
```

## Performance

VanGo is designed for speed:

- Parallel content processing
- Efficient template caching
- Fast Markdown rendering with goldmark
- Minimal memory footprint
- Quick development server startup

## Contributing

VanGo welcomes contributions! The modular architecture makes it easy to:

- Add new template functions
- Create custom content processors
- Extend the configuration system
- Add new output formats
- Improve performance

## License

This project is open source. See the license file for details.

## Roadmap

Future features planned:

- [ ] Theme system
- [ ] Plugin architecture
- [ ] YAML front matter support
- [ ] Image optimization
- [ ] RSS feed generation
- [ ] Sitemap generation
- [ ] Multi-language support
- [ ] Content collections
- [ ] Advanced templating features

---

Built with ❤️ using Go and modern web technologies.
