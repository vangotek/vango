package main

import (
	"flag"
	"fmt"
	"os"

	"vango/internal/builder"
	"vango/internal/config"
	"vango/internal/server"
)

func main() {
	var (
		mode       = flag.String("mode", "build", "Operation mode: build or serve")
		configPath = flag.String("config", "config.toml", "Path to configuration file")
		port       = flag.Int("port", 1313, "Port for development server")
		help       = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	switch *mode {
	case "build":
		if err := buildSite(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Site built successfully!")

	case "serve":
		if err := serveSite(cfg, *port); err != nil {
			fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown mode: %s\n", *mode)
		fmt.Fprintf(os.Stderr, "Use 'build' or 'serve'\n")
		os.Exit(1)
	}
}

func buildSite(cfg *config.Config) error {
	b := builder.New(cfg)
	return b.Build()
}

func serveSite(cfg *config.Config, port int) error {
	s := server.New(cfg, port)
	return s.Start()
}

func printHelp() {
	fmt.Println("VanGo - A Static Site Generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  vango [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -mode string     Operation mode: build or serve (default \"build\")")
	fmt.Println("  -config string   Path to configuration file (default \"config.toml\")")
	fmt.Println("  -port int        Port for development server (default 1313)")
	fmt.Println("  -help            Show this help information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  vango                    # Build site")
	fmt.Println("  vango -mode serve        # Start development server")
	fmt.Println("  vango -mode serve -port 8080  # Serve on custom port")
}
