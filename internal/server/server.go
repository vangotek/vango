package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vango/internal/builder"
	"vango/internal/config"
)

// Server handles the development server
type Server struct {
	config  *config.Config
	builder *builder.Builder
	port    int
	mux     *http.ServeMux
}

// New creates a new development server
func New(cfg *config.Config, port int) *Server {
	return &Server{
		config:  cfg,
		builder: builder.New(cfg),
		port:    port,
		mux:     http.NewServeMux(),
	}
}

// Start starts the development server
func (s *Server) Start() error {
	// Build site initially
	fmt.Println("Building site for development server...")
	if err := s.builder.Build(); err != nil {
		return fmt.Errorf("initial build failed: %w", err)
	}

	// Setup routes
	s.setupRoutes()

	// Start server
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Development server running at http://localhost%s\n", addr)
	fmt.Println("Press Ctrl+C to stop")

	server := &http.Server{
		Addr:         addr,
		Handler:      s.mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	// Serve static files
	staticDir := filepath.Join(s.config.PublicDir, "static")
	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// Serve generated pages
	s.mux.HandleFunc("/", s.handlePage)

	// API endpoints for development
	s.mux.HandleFunc("/api/rebuild", s.handleRebuild)
	s.mux.HandleFunc("/api/status", s.handleStatus)
}

// handlePage serves individual pages
func (s *Server) handlePage(w http.ResponseWriter, r *http.Request) {
	// Clean the path
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index"
	}

	// Try to find the page file
	pagePath := filepath.Join(s.config.PublicDir, path, "index.html")
	
	// If not found, try without subdirectory
	if _, err := os.Stat(pagePath); os.IsNotExist(err) {
		pagePath = filepath.Join(s.config.PublicDir, path+".html")
	}

	// If still not found, try with index.html
	if _, err := os.Stat(pagePath); os.IsNotExist(err) {
		pagePath = filepath.Join(s.config.PublicDir, "index.html")
	}

	// Serve the file if it exists
	if _, err := os.Stat(pagePath); err == nil {
		http.ServeFile(w, r, pagePath)
		return
	}

	// Return 404 if page not found
	s.handle404(w, r)
}

// handleRebuild rebuilds the site
func (s *Server) handleRebuild(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Rebuilding site...")
	if err := s.builder.Build(); err != nil {
		http.Error(w, fmt.Sprintf("Build failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "Site rebuilt successfully"}`))
}

// handleStatus returns server status
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	pages := s.builder.GetPages()
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := fmt.Sprintf(`{
		"status": "running",
		"pages": %d,
		"config": {
			"title": "%s",
			"baseURL": "%s"
		}
	}`, len(pages), s.config.Title, s.config.BaseURL)
	
	w.Write([]byte(response))
}

// handle404 serves a 404 page
func (s *Server) handle404(w http.ResponseWriter, r *http.Request) {
	// Try to serve custom 404 page
	notFoundPath := filepath.Join(s.config.PublicDir, "404.html")
	if _, err := os.Stat(notFoundPath); err == nil {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, notFoundPath)
		return
	}

	// Serve default 404
	w.WriteHeader(http.StatusNotFound)
	html := `<!DOCTYPE html>
<html>
<head>
    <title>404 - Page Not Found</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
        h1 { color: #333; }
        p { color: #666; }
        a { color: #007bff; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <h1>404 - Page Not Found</h1>
    <p>The page you requested could not be found.</p>
    <p><a href="/">← Back to Home</a></p>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// Logger middleware for request logging
func (s *Server) logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap the response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next(wrapped, r)
		
		duration := time.Since(start)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapped.statusCode, duration)
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
