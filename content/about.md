+++
title = "About VanGo"
date = "2025-08-05T00:00:00Z"
description = "Learn more about the VanGo static site generator and its philosophy"
author = "VanGo Team"
tags = ["about", "philosophy", "static-site-generator"]
categories = ["About"]
+++

## About VanGo

VanGo is a modern static site generator built with Go, designed to combine simplicity with powerful features. Our goal is to provide developers and content creators with a tool that's both easy to use and highly customizable.

### Philosophy

We believe that static site generators should be:

1. **Fast** - Both in build times and runtime performance
2. **Simple** - Easy to understand and get started with
3. **Flexible** - Capable of handling diverse use cases
4. **Reliable** - Consistent behavior across different environments

### Why Go?

We chose Go as our foundation because:

- **Performance** - Go's compiled nature provides excellent speed
- **Simplicity** - Go's syntax is clean and maintainable
- **Concurrency** - Built-in support for parallel processing
- **Cross-platform** - Single binary deployment across platforms
- **Strong typing** - Helps prevent runtime errors

### Architecture

VanGo follows a modular architecture with separate packages for:

- **Config** - Site configuration management
- **Content** - Markdown parsing and page handling
- **Template** - HTML template rendering
- **Builder** - Static site generation
- **Server** - Development server with live reload

### Features in Detail

#### Content Management
- TOML front matter support
- Markdown rendering with GitHub Flavored Markdown
- Automatic reading time calculation
- Draft and future post handling
- Tags and categories support

#### Template System
- Go's powerful `html/template` engine
- Rich set of built-in functions
- Custom function support
- Template inheritance
- Flexible layout system

#### Development Experience
- Live development server
- Automatic rebuilding on changes
- API endpoints for debugging
- Comprehensive error messages
- Hot reloading support

#### Build System
- Fast parallel processing
- Static asset copying
- Clean build options
- Configurable output directories
- Cross-platform compatibility

### Getting Help

If you need assistance with VanGo:

1. Check the documentation in your content files
2. Review the template examples
3. Examine the configuration options
4. Look at the source code for advanced usage

### Contributing

VanGo is designed to be extensible. The modular architecture makes it easy to:

- Add new template functions
- Create custom content processors
- Extend the configuration system
- Add new output formats

### Future Plans

We're continuously improving VanGo with features like:

- Theme system
- Plugin architecture
- Advanced templating features
- Performance optimizations
- Enhanced content management

---

*Built with ❤️ using Go and modern web technologies.*
