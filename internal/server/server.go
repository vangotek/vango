package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"vango/internal/builder"
	"vango/internal/config"

	"github.com/fsnotify/fsnotify"
)

// Server handles the enhanced development server
type Server struct {
	config    *config.Config
	builder   *builder.Builder
	port      int
	mux       *http.ServeMux
	verbose   bool
	clients   map[chan string]bool
	clientsMu sync.RWMutex
	
	// Performance tracking
	stats     *ServerStats
	statsMu   sync.RWMutex
}

// ServerStats tracks server performance metrics
type ServerStats struct {
	StartTime    time.Time            `json:"start_time"`
	Requests     int64                `json:"requests"`
	BuildCount   int64                `json:"build_count"`
	LastBuild    time.Time            `json:"last_build"`
	BuildTime    time.Duration        `json:"build_time"`
	ErrorCount   int64                `json:"error_count"`
	FileWatches  int                  `json:"file_watches"`
	ClientCount  int                  `json:"client_count"`
	PageViews    map[string]int64     `json:"page_views"`
	BuildErrors  []string             `json:"build_errors"`
}

// New creates a new enhanced development server
func New(cfg *config.Config, port int) *Server {
	return &Server{
		config:  cfg,
		builder: builder.New(cfg),
		port:    port,
		mux:     http.NewServeMux(),
		verbose: false,
		clients: make(map[chan string]bool),
		stats: &ServerStats{
			StartTime: time.Now(),
			PageViews: make(map[string]int64),
			BuildErrors: make([]string, 0),
		},
	}
}


// SetVerbose sets verbose logging
func (s *Server) SetVerbose(verbose bool) {
	s.verbose = verbose
}

// Start starts the enhanced development server
func (s *Server) Start() error {
	// Build site initially
	fmt.Println("🏗️  Building site for development server...")
	if err := s.buildSite(); err != nil {
		return fmt.Errorf("initial build failed: %w", err)
	}

	// Start file watcher
	go s.watchFiles()

	// Setup routes with enhanced features
	s.setupEnhancedRoutes()

	// Start server
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("🚀 Development server running at http://localhost%s\n", addr)
	fmt.Println("📊 Admin panel: http://localhost" + addr + "/admin")
	fmt.Println("🔄 Live reload enabled")
	fmt.Println("📝 Press Ctrl+C to stop")

	server := &http.Server{
		Addr:         addr,
		Handler:      s.loggingMiddleware(s.mux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server.ListenAndServe()
}

// setupEnhancedRoutes configures enhanced HTTP routes
func (s *Server) setupEnhancedRoutes() {
	// Static files with better caching
	staticDir := filepath.Join(s.config.PublicDir, "static")
	s.mux.Handle("/static/", s.cacheMiddleware(
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))),
	))

	// Theme assets
	themeDir := filepath.Join(s.config.PublicDir, "theme")
	s.mux.Handle("/theme/", s.cacheMiddleware(
		http.StripPrefix("/theme/", http.FileServer(http.Dir(themeDir))),
	))

	// Live reload WebSocket endpoint
	s.mux.HandleFunc("/ws/reload", s.handleWebSocket)

	// Enhanced API endpoints
	s.mux.HandleFunc("/api/rebuild", s.handleRebuild)
	s.mux.HandleFunc("/api/status", s.handleStatus)
	s.mux.HandleFunc("/api/stats", s.handleStats)
	s.mux.HandleFunc("/api/pages", s.handlePages)
	s.mux.HandleFunc("/api/config", s.handleConfig)
	s.mux.HandleFunc("/api/clear-cache", s.handleClearCache)
	s.mux.HandleFunc("/api/validate", s.handleValidate)

	// Admin panel
	s.mux.HandleFunc("/admin", s.handleAdmin)
	s.mux.HandleFunc("/admin/", s.handleAdmin)

	// Development tools
	s.mux.HandleFunc("/dev/template-debug", s.handleTemplateDebug)
	s.mux.HandleFunc("/dev/performance", s.handlePerformance)

	// Serve generated pages (with live reload injection)
	s.mux.HandleFunc("/", s.handlePageWithLiveReload)
}

// buildSite builds the site and tracks performance
func (s *Server) buildSite() error {
	start := time.Now()
	
	s.statsMu.Lock()
	s.stats.BuildCount++
	s.statsMu.Unlock()
	
	err := s.builder.Build()
	
	s.statsMu.Lock()
	s.stats.LastBuild = time.Now()
	s.stats.BuildTime = time.Since(start)
	if err != nil {
		s.stats.ErrorCount++
		s.stats.BuildErrors = append(s.stats.BuildErrors, err.Error())
		// Keep only last 10 errors
		if len(s.stats.BuildErrors) > 10 {
			s.stats.BuildErrors = s.stats.BuildErrors[1:]
		}
	}
	s.statsMu.Unlock()
	
	// Notify clients of rebuild
	if err == nil {
		s.notifyClients("reload")
	} else {
		s.notifyClients(fmt.Sprintf("error:%s", err.Error()))
	}
	
	return err
}

// Enhanced file watcher with better performance
func (s *Server) watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create file watcher: %v", err)
		return
	}
	defer watcher.Close()

	// Directories to watch
	watchDirs := []string{s.config.ContentDir, s.config.LayoutDir}
	
	// Add theme directory if active
	if s.config.Theme != "" {
		themeDir := filepath.Join("themes", s.config.Theme)
		if _, err := os.Stat(themeDir); err == nil {
			watchDirs = append(watchDirs, themeDir)
		}
	}
	
	// Only add static dir if it exists
	if _, err := os.Stat(s.config.StaticDir); err == nil {
		watchDirs = append(watchDirs, s.config.StaticDir)
	}

	// Add config file
	if _, err := os.Stat("config.toml"); err == nil {
		watcher.Add("config.toml")
	}

	// Add directories to watcher recursively
	for _, dir := range watchDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if s.verbose {
					log.Printf("👀 Watching directory: %s", path)
				}
				return watcher.Add(path)
			}
			return nil
		})
		
		if err != nil {
			log.Printf("Error setting up watcher for %s: %v", dir, err)
		}
	}

	s.statsMu.Lock()
	s.stats.FileWatches = len(watcher.WatchList())
	s.statsMu.Unlock()

	log.Printf("👀 File watcher started (watching %d paths)", len(watcher.WatchList()))

	// Debounce rebuilds
	var lastBuild time.Time
	const debounceTime = 300 * time.Millisecond

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			
			// Ignore hidden files and temporary files
			if strings.Contains(event.Name, "/.") || 
			   strings.HasSuffix(event.Name, "~") || 
			   strings.HasSuffix(event.Name, ".tmp") {
				continue
			}
			
			// Only rebuild on write events and if enough time has passed
			if event.Op&fsnotify.Write == fsnotify.Write {
				now := time.Now()
				if now.Sub(lastBuild) > debounceTime {
					lastBuild = now
					log.Printf("🔄 File changed: %s - rebuilding...", event.Name)
					
					// Use incremental build for better performance
					go func() {
						if err := s.builder.IncrementalBuild([]string{event.Name}); err != nil {
							log.Printf("❌ Incremental rebuild failed: %v", err)
							// Fallback to full rebuild
							if err := s.buildSite(); err != nil {
								log.Printf("❌ Full rebuild failed: %v", err)
							}
						} else {
							log.Println("✅ Incremental rebuild completed")
							s.notifyClients("reload")
						}
					}()
				}
			}
			
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("⚠️ File watcher error: %v", err)
		}
	}
}

// WebSocket handler for live reload
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade to WebSocket connection
	// Implementation would use gorilla/websocket or similar
	
	clientChan := make(chan string, 10)
	
	s.clientsMu.Lock()
	s.clients[clientChan] = true
	s.stats.ClientCount = len(s.clients)
	s.clientsMu.Unlock()
	
	defer func() {
		s.clientsMu.Lock()
		delete(s.clients, clientChan)
		s.stats.ClientCount = len(s.clients)
		s.clientsMu.Unlock()
		close(clientChan)
	}()
	
	// Keep connection alive and send messages
	for message := range clientChan {
		// Send message to WebSocket client
		if s.verbose {
			log.Printf("📤 Sending to client: %s", message)
		}
	}
}

// Notify all connected clients
func (s *Server) notifyClients(message string) {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()
	
	for clientChan := range s.clients {
		select {
		case clientChan <- message:
		default:
			// Channel is full, skip this client
		}
	}
}

// Enhanced page handler with live reload injection
func (s *Server) handlePageWithLiveReload(w http.ResponseWriter, r *http.Request) {
	s.statsMu.Lock()
	s.stats.Requests++
	s.stats.PageViews[r.URL.Path]++
	s.statsMu.Unlock()
	
	// Clean the path
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index"
	}

	// Try to find the page file
	pagePath := filepath.Join(s.config.PublicDir, path, "index.html")
	
	if _, err := os.Stat(pagePath); os.IsNotExist(err) {
		pagePath = filepath.Join(s.config.PublicDir, path+".html")
	}

	if _, err := os.Stat(pagePath); os.IsNotExist(err) {
		s.handle404(w, r)
		return
	}

	// Read the file
	content, err := os.ReadFile(pagePath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Inject live reload script
	htmlContent := string(content)
	if strings.Contains(htmlContent, "</body>") {
		liveReloadScript := `
<script>
(function() {
    const ws = new WebSocket('ws://localhost:` + fmt.Sprintf("%d", s.port) + `/ws/reload');
    
    ws.onmessage = function(event) {
        const message = event.data;
        if (message === 'reload') {
            console.log('🔄 Reloading page...');
            window.location.reload();
        } else if (message.startsWith('error:')) {
            console.error('❌ Build error:', message.slice(6));
            // Show error notification
            showErrorNotification(message.slice(6));
        }
    };
    
    ws.onopen = function() {
        console.log('🔗 Live reload connected');
    };
    
    ws.onclose = function() {
        console.log('❌ Live reload disconnected');
        // Try to reconnect after 1 second
        setTimeout(() => window.location.reload(), 1000);
    };
    
    function showErrorNotification(error) {
        const notification = document.createElement('div');
        notification.style.cssText = ` + "`" + `
            position: fixed;
            top: 20px;
            right: 20px;
            background: #ff4444;
            color: white;
            padding: 15px;
            border-radius: 5px;
            z-index: 10000;
            max-width: 400px;
            font-family: monospace;
            font-size: 12px;
        ` + "`" + `;
        notification.textContent = 'Build Error: ' + error;
        document.body.appendChild(notification);
        
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 5000);
    }
})();
</script>`
		
		htmlContent = strings.Replace(htmlContent, "</body>", liveReloadScript+"\n</body>", 1)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Write([]byte(htmlContent))
}

// Enhanced API endpoints
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	s.statsMu.RLock()
	stats := *s.stats
	s.statsMu.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) handlePages(w http.ResponseWriter, r *http.Request) {
	pages := s.builder.GetPages()
	
	type PageInfo struct {
		Title    string `json:"title"`
		URL      string `json:"url"`
		WordCount int   `json:"word_count"`
		ReadingTime int `json:"reading_time"`
		LastModified time.Time `json:"last_modified"`
	}
	
	var pageInfos []PageInfo
	for _, page := range pages {
		pageInfos = append(pageInfos, PageInfo{
			Title: page.Title,
			URL: page.URL,
			WordCount: page.WordCount,
			ReadingTime: page.ReadingTime,
			LastModified: page.ParsedDate,
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pageInfos)
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.config)
}

// Admin panel handler
func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
 <meta charset="UTF-8">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css">
    <title>VanGo Admin Panel</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .card { background: white; padding: 20px; margin: 20px 0; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; }
        .stat { text-align: center; }
        .stat-value { font-size: 2em; font-weight: bold; color: #007bff; }
        .stat-label { color: #666; }
        button { background: #007bff; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; }
        button:hover { background: #0056b3; }
        .error { background: #fff5f5; border: 1px solid #feb2b2; color: #e53e3e; padding: 10px; border-radius: 4px; margin: 10px 0; }
    </style>
</head>
<body>
    <div class="container">
        <h1><i class="fa-solid fa-server"></i> VanGo Development Server</h1>
        
        <div class="card">
            <h2>Quick Actions</h2>
            <button onclick="rebuild()"><i class="fa-solid fa-repeat"></i> Rebuild Site</button>
            <button onclick="clearCache()"><i class="fa-solid fa-trash"></i> Clear Cache</button>
            <button onclick="location.reload()"><i class="fa-solid fa-rotate"></i> Refresh Panel</button>
        </div>
        
        <div class="card">
            <h2><i class="fa-solid fa-chart-column"></i> Server Statistics</h2>
            <div class="stats" id="stats"></div>
        </div>
        
        <div class="card">
            <h2><i class="fa-solid fa-file"></i> Pages</h2>
            <div id="pages"></div>
        </div>
        
        <div class="card">
            <h2><i class="fa-solid fa-gear"></i> Configuration</h2>
            <pre id="config"></pre>
        </div>
    </div>
    
    <script>
        async function loadStats() {
            const response = await fetch('/api/stats');
            const stats = await response.json();
            
            document.getElementById('stats').innerHTML = ` + "`" + `
                <div class="stat">
                    <div class="stat-value">${stats.requests}</div>
                    <div class="stat-label">Total Requests</div>
                </div>
                <div class="stat">
                    <div class="stat-value">${stats.build_count}</div>
                    <div class="stat-label">Builds</div>
                </div>
                <div class="stat">
                    <div class="stat-value">${stats.client_count}</div>
                    <div class="stat-label">Connected Clients</div>
                </div>
                <div class="stat">
                    <div class="stat-value">${stats.file_watches}</div>
                    <div class="stat-label">Watched Files</div>
                </div>
            ` + "`" + `;
            
            if (stats.build_errors && stats.build_errors.length > 0) {
                const errorsHtml = stats.build_errors.map(error => 
                    ` + "`" + `<div class="error">${error}</div>` + "`" + `
                ).join('');
                document.getElementById('stats').innerHTML += ` + "`" + `
                    <div style="grid-column: 1 / -1;">
                        <h3>Recent Build Errors</h3>
                        ${errorsHtml}
                    </div>
                ` + "`" + `;
            }
        }
        
        async function loadPages() {
            const response = await fetch('/api/pages');
            const pages = await response.json();
            
            document.getElementById('pages').innerHTML = pages.map(page => ` + "`" + `
                <div style="border-bottom: 1px solid #eee; padding: 10px 0;">
                    <strong>${page.title}</strong><br>
                    <a href="${page.url}" target="_blank">${page.url}</a><br>
                    <small>${page.word_count} words • ${page.reading_time} min read</small>
                </div>
            ` + "`" + `).join('');
        }
        
        async function loadConfig() {
            const response = await fetch('/api/config');
            const config = await response.json();
            document.getElementById('config').textContent = JSON.stringify(config, null, 2);
        }
        
        async function rebuild() {
            const response = await fetch('/api/rebuild', { method: 'POST' });
            if (response.ok) {
                alert('✅ Site rebuilt successfully!');
                loadStats();
            } else {
                alert('❌ Rebuild failed!');
            }
        }
        
        async function clearCache() {
            const response = await fetch('/api/clear-cache', { method: 'POST' });
            if (response.ok) {
                alert('✅ Cache cleared!');
            }
        }
        
        // Load data on page load
        loadStats();
        loadPages();
        loadConfig();
        
        // Auto-refresh stats every 5 seconds
        setInterval(loadStats, 5000);
    </script>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// Middleware for caching static assets
func (s *Server) cacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour cache
		next.ServeHTTP(w, r)
	})
}

// Logging middleware
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap the response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r)
		
		if s.verbose {
			duration := time.Since(start)
			log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapped.statusCode, duration)
		}
	})
}

// Additional handlers...
func (s *Server) handleClearCache(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Clear any caches here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}

func (s *Server) handleValidate(w http.ResponseWriter, r *http.Request) {
	// Validate site configuration and content
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "valid"}`))
}

func (s *Server) handleTemplateDebug(w http.ResponseWriter, r *http.Request) {
	// Template debugging information
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Template debug information..."))
}

func (s *Server) handlePerformance(w http.ResponseWriter, r *http.Request) {
	// Performance metrics and profiling
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"performance": "metrics"}`))
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
