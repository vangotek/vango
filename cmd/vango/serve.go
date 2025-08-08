package vango

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vango/internal/config"
	"vango/internal/server"
)

var (
	servePort int
	serveHost string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the development server",
	Long: `Start the development server with live reload.

The development server watches for changes in your content, layout,
and static files, automatically rebuilding and refreshing the site.

Server provides:
  • Live preview at http://localhost:1313
  • Automatic rebuilding on file changes  
  • API endpoints for debugging
  • Hot reload support`,
	Example: `  vango serve                     # Start server on default port (1313)
  vango serve -p 8080             # Start server on port 8080
  vango serve --host 0.0.0.0      # Bind to all interfaces
  vango serve -v                  # Start with verbose output`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("🚀 Starting development server...")
		}
		
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Override config with command line flags
		if servePort != 1313 {
			cfg.Port = servePort
		}
		if serveHost != "localhost" {
			cfg.Host = serveHost
		}

		if verbose {
			fmt.Printf("🏠 Site: %s\n", cfg.Title)
			fmt.Printf("🌐 Host: %s\n", cfg.Host)
			fmt.Printf("🔌 Port: %d\n", cfg.Port)
			fmt.Printf("🔄 Live Reload: %v\n", cfg.LiveReload)
		}

		s := server.New(cfg, cfg.Port)
		s.SetVerbose(verbose) // Pass verbose flag to server
		fmt.Printf("🎨 Development server starting...\n")
		fmt.Printf("🔗 Local: http://%s:%d\n", cfg.Host, cfg.Port)
		fmt.Println("📝 Press Ctrl+C to stop")
		if err := s.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Server failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 1313, "Port for development server")
	serveCmd.Flags().StringVar(&serveHost, "host", "localhost", "Host to bind to")
}

