package main

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"vango/internal/config"
	vangoTemplate "vango/internal/template"
)

// Example of how to extend VanGo with custom template functions

// CustomFunctions demonstrates how to add custom template functions
func CustomFunctions() template.FuncMap {
	return template.FuncMap{
		// Custom date formatting
		"customDate": func(date time.Time) string {
			return date.Format("Mon, Jan 2, 2006")
		},
		
		// Word truncation
		"truncate": func(text string, length int) string {
			words := strings.Fields(text)
			if len(words) <= length {
				return text
			}
			return strings.Join(words[:length], " ") + "..."
		},
		
		// Reading time with custom calculation
		"customReadingTime": func(wordCount int, wordsPerMinute int) int {
			if wordsPerMinute == 0 {
				wordsPerMinute = 200 // default
			}
			return (wordCount + wordsPerMinute - 1) / wordsPerMinute
		},
		
		// Slug generation
		"slugify": func(text string) string {
			text = strings.ToLower(text)
			text = strings.ReplaceAll(text, " ", "-")
			// Remove special characters (simplified)
			replacer := strings.NewReplacer(
				"!", "", "?", "", ".", "", ",", "", ":", "", ";", "",
				"'", "", "\"", "", "(", "", ")", "", "[", "", "]", "",
				"{", "", "}", "", "/", "", "\\", "", "|", "",
			)
			return replacer.Replace(text)
		},
		
		// Excerpt generation
		"excerpt": func(content template.HTML, length int) string {
			// Strip HTML tags (simplified)
			text := string(content)
			text = strings.ReplaceAll(text, "<", " <")
			// This is a very basic HTML stripper - in production use a proper library
			words := strings.Fields(text)
			var cleanWords []string
			for _, word := range words {
				if !strings.HasPrefix(word, "<") {
					cleanWords = append(cleanWords, word)
				}
			}
			
			if len(cleanWords) <= length {
				return strings.Join(cleanWords, " ")
			}
			return strings.Join(cleanWords[:length], " ") + "..."
		},
		
		// Social sharing URLs
		"twitterShare": func(url, title string) string {
			return fmt.Sprintf("https://twitter.com/intent/tweet?url=%s&text=%s", url, title)
		},
		
		"facebookShare": func(url string) string {
			return fmt.Sprintf("https://www.facebook.com/sharer/sharer.php?u=%s", url)
		},
		
		// URL manipulation
		"absURL": func(baseURL, path string) string {
			if strings.HasSuffix(baseURL, "/") && strings.HasPrefix(path, "/") {
				return baseURL + path[1:]
			} else if !strings.HasSuffix(baseURL, "/") && !strings.HasPrefix(path, "/") {
				return baseURL + "/" + path
			}
			return baseURL + path
		},
		
		// Math functions
		"percentage": func(part, total int) float64 {
			if total == 0 {
				return 0
			}
			return float64(part) / float64(total) * 100
		},
		
		// String utilities
		"removeHTML": func(s string) string {
			// Very basic HTML removal - use a proper library in production
			result := ""
			skip := false
			for _, char := range s {
				if char == '<' {
					skip = true
				} else if char == '>' {
					skip = false
				} else if !skip {
					result += string(char)
				}
			}
			return strings.TrimSpace(result)
		},
		
		// Array/slice utilities
		"first": func(limit int, slice interface{}) interface{} {
			// This would need proper type handling in a real implementation
			return slice // Simplified for example
		},
		
		"last": func(limit int, slice interface{}) interface{} {
			// This would need proper type handling in a real implementation
			return slice // Simplified for example
		},
		
		// JSON handling
		"toJSON": func(v interface{}) string {
			// Would need proper JSON marshaling
			return fmt.Sprintf("%v", v) // Simplified for example
		},
	}
}

// Example of how to create a custom template engine with additional functions
func NewCustomEngine(cfg *config.Config) *vangoTemplate.Engine {
	engine := vangoTemplate.NewEngine(cfg)
	
	// In a real implementation, you would need to modify the template engine
	// to accept custom functions. This is just an example structure.
	
	return engine
}

// Example usage in a custom build script
func ExampleCustomBuild() {
	fmt.Println("This is an example of how to extend VanGo with custom functionality")
	fmt.Println("You can:")
	fmt.Println("1. Add custom template functions")
	fmt.Println("2. Create custom content processors")
	fmt.Println("3. Add new configuration options")
	fmt.Println("4. Extend the server with custom endpoints")
	fmt.Println("5. Add custom output formats")
}

// Custom content processor example
type CustomContentProcessor struct {
	// Custom fields
}

func (p *CustomContentProcessor) ProcessCustomFormat(content string) (string, error) {
	// Example: process a custom content format
	// This could be for special markdown extensions, custom shortcodes, etc.
	return content, nil
}

// Example of extending configuration
type CustomConfig struct {
	*config.Config
	
	// Custom configuration options
	CustomField1 string            `toml:"customField1"`
	CustomField2 int               `toml:"customField2"`
	CustomParams map[string]string `toml:"customParams"`
}

func main() {
	ExampleCustomBuild()
}
