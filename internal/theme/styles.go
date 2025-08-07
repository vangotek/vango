package theme

// --- CSS String Functions (getBasicCSS, getBlogCSS, getPortfolioCSS, getDocsCSS) ---
// These functions return large CSS strings. They are kept here for now,
// but ideally, the CSS should be stored in separate .css files and read at runtime
// or embedded using Go 1.16+ embed directive.

// getBasicCSS returns basic theme CSS
func (tm *ThemeManager) getBasicCSS() string {
	return `/* VanGo Basic Theme */
:root {
    --color-primary: #007bff;
    --color-secondary: #6c757d;
    --color-background: #ffffff;
    --color-surface: #f8f9fa;
    --color-text: #333333;
    --color-text-muted: #6c757d;
    --color-border: #e9ecef;
    --font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    --max-width: 800px;
}
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}
body {
    font-family: var(--font-family);
    line-height: 1.6;
    color: var(--color-text);
    background-color: var(--color-background);
}
.site-header {
    background-color: var(--color-surface);
    border-bottom: 1px solid var(--color-border);
    padding: 1rem 0;
}
.nav-container {
    max-width: var(--max-width);
    margin: 0 auto;
    padding: 0 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.site-title {
    font-size: 1.5rem;
    font-weight: bold;
    text-decoration: none;
    color: var(--color-primary);
}
.nav-links {
    display: flex;
    list-style: none;
    gap: 2rem;
}
.nav-links a {
    text-decoration: none;
    color: var(--color-text);
    font-weight: 500;
    transition: color 0.2s ease;
}
.nav-links a:hover {
    color: var(--color-primary);
}
.main-content {
    max-width: var(--max-width);
    margin: 2rem auto;
    padding: 0 2rem;
}
.post {
    background: var(--color-background);
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    padding: 2rem;
}
.post-title {
    font-size: 2.5rem;
    font-weight: bold;
    margin-bottom: 1rem;
    line-height: 1.2;
}
.post-meta {
    color: var(--color-text-muted);
    font-size: 0.9rem;
    margin-bottom: 1rem;
}
.post-content {
    line-height: 1.8;
}
.site-footer {
    background-color: var(--color-surface);
    padding: 2rem 0;
    margin-top: 4rem;
}
.footer-container {
    max-width: var(--max-width);
    margin: 0 auto;
    padding: 0 2rem;
    text-align: center;
}
.home-hero {
    text-align: center;
    padding: 4rem 0;
    background: linear-gradient(135deg, var(--color-primary), #764ba2);
    color: white;
    margin: -2rem -2rem 3rem -2rem;
    border-radius: 0 0 12px 12px;
}
.hero-title {
    font-size: 3.5rem;
    font-weight: bold;
    margin-bottom: 1rem;
}
.posts-list h2 {
    font-size: 2rem;
    margin-bottom: 2rem;
    text-align: center;
}
.post-summary {
    background: var(--color-surface);
    padding: 1.5rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
}
.post-summary h3 {
    margin-bottom: 0.5rem;
}
.post-summary a {
    color: var(--color-text);
    text-decoration: none;
}
.post-summary a:hover {
    color: var(--color-primary);
}
@media (max-width: 768px) {
    .nav-container {
        flex-direction: column;
        gap: 1rem;
    }
    .main-content {
        padding: 0 1rem;
    }
    .hero-title {
        font-size: 2.5rem;
    }
}`
}

// getBlogCSS returns blog theme CSS
func (tm *ThemeManager) getBlogCSS() string {
	return `/* VanGo Blog Theme */
:root {
    --color-primary: #2563eb;
    --color-secondary: #64748b;
    --color-accent: #f59e0b;
    --color-background: #ffffff;
    --color-surface: #f8fafc;
    --color-text: #1e293b;
    --color-text-muted: #64748b;
    --color-border: #e2e8f0;
    --font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    --max-width: 900px;
}
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}
body {
    font-family: var(--font-family);
    line-height: 1.6;
    color: var(--color-text);
    background-color: var(--color-background);
}
.site-header {
    background-color: var(--color-background);
    border-bottom: 2px solid var(--color-primary);
    padding: 1rem 0;
    position: sticky;
    top: 0;
    z-index: 100;
}
.nav-container {
    max-width: var(--max-width);
    margin: 0 auto;
    padding: 0 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.site-title {
    font-size: 1.8rem;
    font-weight: bold;
    text-decoration: none;
    color: var(--color-primary);
    font-family: 'Georgia', serif;
}
.nav-links {
    display: flex;
    list-style: none;
    gap: 2rem;
}
.nav-links a {
    text-decoration: none;
    color: var(--color-text);
    font-weight: 500;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    transition: all 0.2s ease;
}
.nav-links a:hover {
    background-color: var(--color-surface);
    color: var(--color-primary);
}
.main-content {
    max-width: var(--max-width);
    margin: 2rem auto;
    padding: 0 2rem;
}
.post {
    background: var(--color-background);
    border-radius: 12px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07);
    overflow: hidden;
}
.post-header {
    padding: 3rem 3rem 1rem;
    border-bottom: 1px solid var(--color-border);
}
.post-title {
    font-size: 2.75rem;
    font-weight: 800;
    margin-bottom: 1.5rem;
    line-height: 1.1;
    color: var(--color-text);
    font-family: 'Georgia', serif;
}
.post-meta {
    color: var(--color-text-muted);
    font-size: 0.95rem;
    margin-bottom: 1.5rem;
    display: flex;
    align-items: center;
    gap: 1rem;
}
.author {
    font-weight: 600;
    color: var(--color-primary);
}
.post-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-top: 1rem;
}
.tag {
    background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
    color: white;
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 500;
}
.post-content {
    padding: 2rem 3rem 3rem;
    line-height: 1.8;
    font-size: 1.1rem;
}
.post-content h1,
.post-content h2,
.post-content h3 {
    color: var(--color-text);
    margin: 2rem 0 1rem;
    font-weight: 700;
}
.post-content h1 { font-size: 2rem; }
.post-content h2 { font-size: 1.6rem; }
.post-content h3 { font-size: 1.3rem; }
.post-content p {
    margin-bottom: 1.5rem;
}
.post-content a {
    color: var(--color-primary);
    text-decoration: none;
    border-bottom: 1px solid transparent;
    transition: border-color 0.2s ease;
}
.post-content a:hover {
    border-bottom-color: var(--color-primary);
}
.post-content blockquote {
    border-left: 4px solid var(--color-primary);
    padding: 1rem 2rem;
    margin: 2rem 0;
    background-color: var(--color-surface);
    border-radius: 0 8px 8px 0;
    font-style: italic;
}
.post-content code {
    background-color: var(--color-surface);
    padding: 0.2rem 0.4rem;
    border-radius: 4px;
    font-family: 'Courier New', monospace;
    font-size: 0.9em;
    color: var(--color-primary);
}
.post-content pre {
    background-color: var(--color-surface);
    padding: 1.5rem;
    border-radius: 8px;
    overflow-x: auto;
    margin: 1.5rem 0;
    border: 1px solid var(--color-border);
}
.post-share {
    padding: 2rem 3rem;
    border-top: 1px solid var(--color-border);
}
.post-share h4 {
    margin-bottom: 1rem;
    color: var(--color-text);
}
.post-share a {
    display: inline-block;
    padding: 0.5rem 1rem;
    background-color: var(--color-primary);
    color: white;
    text-decoration: none;
    border-radius: 6px;
    margin-right: 1rem;
    transition: transform 0.2s ease;
}
.post-share a:hover {
    transform: translateY(-2px);
}
.blog-hero {
    text-align: center;
    padding: 5rem 2rem;
    background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-accent) 100%);
    color: white;
    margin-bottom: 4rem;
    border-radius: 0 0 20px 20px;
}
.hero-title {
    font-size: 4rem;
    font-weight: 900;
    margin-bottom: 1rem;
    font-family: 'Georgia', serif;
}
.hero-description {
    font-size: 1.3rem;
    opacity: 0.9;
    max-width: 600px;
    margin: 0 auto;
}
.posts-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
    gap: 2.5rem;
}
.post-card {
    background: var(--color-background);
    border-radius: 12px;
    padding: 2rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
    border: 1px solid var(--color-border);
}
.post-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 25px rgba(0, 0, 0, 0.1);
}
.post-card h2 {
    margin-bottom: 1rem;
    font-size: 1.5rem;
    font-weight: 700;
}
.post-card h2 a {
    color: var(--color-text);
    text-decoration: none;
    transition: color 0.2s ease;
}
.post-card h2 a:hover {
    color: var(--color-primary);
}
.post-card .post-meta {
    margin-bottom: 1rem;
    font-size: 0.9rem;
}
.post-excerpt {
    color: var(--color-text-muted);
    line-height: 1.6;
    margin-bottom: 1.5rem;
}
.site-footer {
    background: linear-gradient(135deg, var(--color-text) 0%, #374151 100%);
    color: white;
    padding: 3rem 0;
    margin-top: 6rem;
}
.footer-container {
    max-width: var(--max-width);
    margin: 0 auto;
    padding: 0 2rem;
    text-align: center;
}
@media (max-width: 768px) {
    .nav-container {
        flex-direction: column;
        gap: 1rem;
    }
    .post-header,
    .post-content,
    .post-share {
        padding-left: 1.5rem;
        padding-right: 1.5rem;
    }
    .hero-title {
        font-size: 2.8rem;
    }
    .posts-grid {
        grid-template-columns: 1fr;
        gap: 1.5rem;
    }
    .post-title {
        font-size: 2.2rem;
    }
}`
}

// getPortfolioCSS returns portfolio theme CSS
func (tm *ThemeManager) getPortfolioCSS() string {
	return `/* VanGo Portfolio Theme */
:root {
    --color-primary: #6366f1;
    --color-secondary: #8b5cf6;
    --color-accent: #06b6d4;
    --color-background: #0f172a;
    --color-surface: #1e293b;
    --color-text: #f1f5f9;
    --color-text-muted: #94a3b8;
    --color-border: #334155;
    --font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    --max-width: 1200px;
}
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}
body {
    font-family: var(--font-family);
    line-height: 1.6;
    color: var(--color-text);
    background-color: var(--color-background);
}
.portfolio-nav {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    background: rgba(15, 23, 42, 0.9);
    backdrop-filter: blur(10px);
    border-bottom: 1px solid var(--color-border);
    padding: 1rem 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    z-index: 100;
}
.nav-logo {
    font-size: 1.5rem;
    font-weight: bold;
    color: var(--color-primary);
    text-decoration: none;
}
.nav-menu {
    display: flex;
    list-style: none;
    gap: 2rem;
}
.nav-menu a {
    color: var(--color-text);
    text-decoration: none;
    font-weight: 500;
    transition: color 0.2s ease;
}
.nav-menu a:hover {
    color: var(--color-primary);
}
.portfolio-main {
    margin-top: 80px;
    min-height: calc(100vh - 80px);
}
.hero-section {
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, var(--color-background) 0%, var(--color-surface) 100%);
    position: relative;
    overflow: hidden;
}
.hero-section::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: radial-gradient(circle at 50% 50%, rgba(99, 102, 241, 0.1) 0%, transparent 70%);
}
.hero-content {
    text-align: center;
    z-index: 2;
    max-width: 800px;
    padding: 2rem;
}
.hero-title {
    font-size: 4rem;
    font-weight: 900;
    margin-bottom: 1.5rem;
    background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}
.hero-subtitle {
    font-size: 1.3rem;
    color: var(--color-text-muted);
    margin-bottom: 2rem;
    line-height: 1.6;
}
.projects-section {
    padding: 5rem 2rem;
    max-width: var(--max-width);
    margin: 0 auto;
}
.projects-section h2 {
    font-size: 2.5rem;
    text-align: center;
    margin-bottom: 3rem;
    color: var(--color-text);
}
.projects-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
    gap: 2.5rem;
}
.project-card {
    background: var(--color-surface);
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    border: 1px solid var(--color-border);
}
.project-card:hover {
    transform: translateY(-8px);
    box-shadow: 0 16px 48px rgba(99, 102, 241, 0.2);
}
.project-image {
    height: 200px;
    background: linear-gradient(135deg, var(--color-primary), var(--color-secondary));
    position: relative;
    overflow: hidden;
}
.project-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}
.project-placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 4rem;
    font-weight: bold;
    color: white;
    height: 100%;
}
.project-info {
    padding: 2rem;
}
.project-info h3 {
    font-size: 1.3rem;
    margin-bottom: 1rem;
    color: var(--color-text);
}
.project-info h3 a {
    color: inherit;
    text-decoration: none;
    transition: color 0.2s ease;
}
.project-info h3 a:hover {
    color: var(--color-primary);
}
.project-description {
    color: var(--color-text-muted);
    margin-bottom: 1.5rem;
    line-height: 1.6;
}
.project-tech {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
}
.tech-tag {
    background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
    color: white;
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 500;
}
.project-detail {
    max-width: 900px;
    margin: 0 auto;
    padding: 2rem;
}
.project-header {
    text-align: center;
    margin-bottom: 3rem;
}
.project-title {
    font-size: 3rem;
    margin-bottom: 2rem;
    background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}
.project-content {
    color: var(--color-text-muted);
    font-size: 1.1rem;
    line-height: 1.8;
    margin-bottom: 3rem;
}
.project-links {
    display: flex;
    gap: 1rem;
    justify-content: center;
}
.btn {
    display: inline-block;
    padding: 1rem 2rem;
    text-decoration: none;
    border-radius: 8px;
    font-weight: 600;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.btn-primary {
    background: linear-gradient(135deg, var(--color-primary), var(--color-secondary));
    color: white;
}
.btn-secondary {
    background: transparent;
    color: var(--color-primary);
    border: 2px solid var(--color-primary);
}
.btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(99, 102, 241, 0.3);
}
@media (max-width: 768px) {
    .portfolio-nav {
        flex-direction: column;
        gap: 1rem;
        padding: 1rem;
    }
    .nav-menu {
        gap: 1rem;
    }
    .hero-title {
        font-size: 2.8rem;
    }
    .projects-grid {
        grid-template-columns: 1fr;
        gap: 1.5rem;
    }
    .project-links {
        flex-direction: column;
        align-items: center;
    }
}`
}

// getDocsCSS returns documentation theme CSS
func (tm *ThemeManager) getDocsCSS() string {
	return `/* VanGo Docs Theme */
:root {
    --color-primary: #059669;
    --color-secondary: #6b7280;
    --color-accent: #3b82f6;
    --color-background: #ffffff;
    --color-surface: #f9fafb;
    --color-text: #111827;
    --color-text-muted: #6b7280;
    --color-border: #e5e7eb;
    --color-sidebar: #f3f4f6;
    --font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    --sidebar-width: 280px;
}
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}
body {
    font-family: var(--font-family);
    line-height: 1.6;
    color: var(--color-text);
    background-color: var(--color-background);
}
.docs-layout {
    min-height: 100vh;
}
.docs-nav {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: 60px;
    background-color: var(--color-background);
    border-bottom: 1px solid var(--color-border);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 2rem;
    z-index: 100;
}
.nav-brand a {
    font-size: 1.3rem;
    font-weight: bold;
    color: var(--color-primary);
    text-decoration: none;
}
.nav-search input {
    width: 300px;
    padding: 0.5rem 1rem;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    font-size: 0.9rem;
    background-color: var(--color-surface);
}
.docs-container {
    display: flex;
    margin-top: 60px;
    min-height: calc(100vh - 60px);
}
.docs-sidebar {
    width: var(--sidebar-width);
    background-color: var(--color-sidebar);
    border-right: 1px solid var(--color-border);
    padding: 2rem 0;
    position: fixed;
    left: 0;
    top: 60px;
    bottom: 0;
    overflow-y: auto;
}
.sidebar-nav {
    padding: 0 2rem;
}
.sidebar-nav h3 {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
    margin-bottom: 1rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}
.sidebar-nav ul {
    list-style: none;
}
.sidebar-nav li {
    margin-bottom: 0.5rem;
}
.sidebar-nav a {
    color: var(--color-text-muted);
    text-decoration: none;
    display: block;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    font-size: 0.9rem;
    transition: all 0.2s ease;
}
.sidebar-nav a:hover,
.sidebar-nav a.active {
    background-color: var(--color-primary);
    color: white;
}
.docs-main {
    flex: 1;
    margin-left: var(--sidebar-width);
    padding: 2rem;
    max-width: calc(100% - var(--sidebar-width));
}
.docs-article {
    max-width: 800px;
    margin: 0 auto;
}
.docs-header {
    margin-bottom: 3rem;
    padding-bottom: 2rem;
    border-bottom: 1px solid var(--color-border);
}
.docs-header h1 {
    font-size: 2.5rem;
    font-weight: 800;
    margin-bottom: 1rem;
    color: var(--color-text);
}
.docs-description {
    font-size: 1.2rem;
    color: var(--color-text-muted);
    line-height: 1.6;
}
.docs-toc {
    background-color: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 1.5rem;
    margin-bottom: 2rem;
}
.docs-toc h4 {
    margin-bottom: 1rem;
    color: var(--color-text);
    font-size: 1rem;
}
.docs-content {
    font-size: 1rem;
    line-height: 1.8;
    color: var(--color-text);
}
.docs-content h1,
.docs-content h2,
.docs-content h3,
.docs-content h4 {
    margin: 2rem 0 1rem;
    color: var(--color-text);
    font-weight: 600;
}
.docs-content h1 { font-size: 2rem; border-bottom: 2px solid var(--color-border); padding-bottom: 0.5rem; }
.docs-content h2 { font-size: 1.6rem; }
.docs-content h3 { font-size: 1.3rem; }
.docs-content h4 { font-size: 1.1rem; }
.docs-content p {
    margin-bottom: 1.5rem;
}
.docs-content a {
    color: var(--color-primary);
    text-decoration: none;
    border-bottom: 1px solid transparent;
    transition: border-color 0.2s ease;
}
.docs-content a:hover {
    border-bottom-color: var(--color-primary);
}
.docs-content blockquote {
    border-left: 4px solid var(--color-primary);
    background-color: var(--color-surface);
    padding: 1rem 1.5rem;
    margin: 2rem 0;
    border-radius: 0 8px 8px 0;
}
.docs-content code {
    background-color: var(--color-surface);
    padding: 0.2rem 0.4rem;
    border-radius: 4px;
    font-family: 'Courier New', monospace;
    font-size: 0.9em;
    color: var(--color-primary);
    border: 1px solid var(--color-border);
}
.docs-content pre {
    background-color: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 1.5rem;
    overflow-x: auto;
    margin: 1.5rem 0;
    font-size: 0.9rem;
}
.docs-content pre code {
    background: none;
    border: none;
    padding: 0;
    color: var(--color-text);
}
.docs-content table {
    width: 100%;
    border-collapse: collapse;
    margin: 1.5rem 0;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    overflow: hidden;
}
.docs-content th,
.docs-content td {
    padding: 1rem;
    text-align: left;
    border-bottom: 1px solid var(--color-border);
}
.docs-content th {
    background-color: var(--color-surface);
    font-weight: 600;
}
.docs-footer {
    margin-top: 4rem;
    padding-top: 2rem;
    border-top: 1px solid var(--color-border);
}
.docs-home {
    max-width: 1000px;
    margin: 0 auto;
}
.docs-hero {
    text-align: center;
    padding: 4rem 0;
    margin-bottom: 3rem;
}
.docs-hero h1 {
    font-size: 3rem;
    margin-bottom: 1.5rem;
    color: var(--color-text);
}
.docs-hero p {
    font-size: 1.3rem;
    color: var(--color-text-muted);
    max-width: 600px;
    margin: 0 auto;
}
.docs-sections {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
}
.docs-card {
    background-color: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 2rem;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.docs-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
.docs-card h2 {
    margin-bottom: 1rem;
    color: var(--color-text);
}
.docs-card h2 a {
    color: inherit;
    text-decoration: none;
    transition: color 0.2s ease;
}
.docs-card h2 a:hover {
    color: var(--color-primary);
}
@media (max-width: 1024px) {
    .docs-sidebar {
        transform: translateX(-100%);
        transition: transform 0.3s ease;
    }
    .docs-main {
        margin-left: 0;
        max-width: 100%;
    }
    .nav-search input {
        width: 200px;
    }
}
@media (max-width: 768px) {
    .docs-nav {
        flex-direction: column;
        height: auto;
        padding: 1rem;
    }
    .nav-search {
        margin-top: 1rem;
    }
    .nav-search input {
        width: 100%;
    }
    .docs-container {
        margin-top: 100px;
    }
    .docs-main {
        padding: 1rem;
    }
    .docs-hero h1 {
        font-size: 2.2rem;
    }
    .docs-sections {
        grid-template-columns: 1fr;
        gap: 1.5rem;
    }
}`
}







