+++
title = "Welcome to VanGo"
date = "2025-08-05T00:00:00Z"
description = "A comprehensive guide to getting started with VanGo static site generator"
author = "VanGo Team"
tags = ["introduction", "guide", "static-site"]
categories = ["Documentation"]
+++

# Welcome to VanGo!

VanGo is a fast, modern static site generator built with Go. It's designed to be simple yet powerful, allowing you to create beautiful websites with minimal configuration.

## Features

- **Fast builds** - Powered by Go's performance
- **Markdown support** - Write content in Markdown with TOML front matter
- **Template engine** - Flexible HTML templates with built-in functions
- **Development server** - Live preview with automatic rebuilding
- **Static asset handling** - Automatic copying and optimization
- **Modular architecture** - Clean, extensible codebase

## Getting Started

### Building Your Site

To build your site, simply run:

```bash
go run main.go
```

Or with custom configuration:

```bash
go run main.go -config custom-config.toml
```

### Development Server

Start the development server for live preview:

```bash
go run main.go -mode serve
```

The server will start at `http://localhost:1313` by default.

### Directory Structure

```
your-site/
├── config.toml          # Site configuration
├── content/             # Your content files
│   ├── hello.md
│   └── about.md
├── layouts/             # HTML templates
│   └── _default/
│       └── single.html
├── static/              # Static assets (CSS, JS, images)
│   └── style.css
└── public/              # Generated site (output)
```

## Template Functions

VanGo provides many built-in template functions:

- `{{ dateFormat "2006-01-02" .Page.Date }}` - Format dates
- `{{ humanizeDate .Page.Date }}` - Human-readable dates
- `{{ .Page.ReadingTime }}` - Calculate reading time
- `{{ range .Page.Tags }}` - Loop through arrays
- `{{ lower "TEXT" }}` - String manipulation

## Configuration

The `config.toml` file controls your site's behavior:

```toml
title = "My Site"
baseURL = "https://mysite.com"
description = "A fantastic website"

[params]
    author = "John Doe"
    version = "1.0"
```

## Next Steps

1. Customize your templates in the `layouts/` directory
2. Add your content in the `content/` directory
3. Style your site with CSS in the `static/` directory
4. Configure your site settings in `config.toml`

Happy building with VanGo! 🚀
