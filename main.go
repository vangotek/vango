package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/yuin/goldmark"
)

type Site struct {
	Title   string
	BaseURL string
}

type Page struct {
	Title   string
	Date    string
	Content template.HTML
	Params  map[string]interface{}
}

type TemplateData struct {
	Site Site
	Page Page
}

func main() {

	mode := flag.String("mode", "build", "Mode: build or serve")
	flag.Parse()

	switch *mode {
	case "build":
		buildSite()
	case "serve":
		serveSite()
	default:
		log.Fatal("Unknown mode:", *mode)
	}

	site := Site{Title: "My Hugo Clone", BaseURL: "http://localhost:1313/"}
	contentDir := "content"
	layoutPath := "layouts/_default/single.html"
	publicDir := "public"
	staticDir := "static"

	tpl := template.Must(template.ParseFiles(layoutPath))
	os.RemoveAll(publicDir)
	os.MkdirAll(publicDir, 0755)

	// Copy static assets
	copyStaticFiles(staticDir, filepath.Join(publicDir, "static"))

	// Walk content
	filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(data)
		var frontMatter, body string

		if strings.HasPrefix(content, "+++") {
			parts := strings.SplitN(content, "+++", 3)
			if len(parts) == 3 {
				frontMatter = parts[1]
				body = parts[2]
			}
		}

		tree, err := toml.Load(frontMatter)
		if err != nil {
			return err
		}

		title := tree.Get("title").(string)
		date := tree.Get("date").(string)
		params := tree.ToMap()

		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(body), &buf); err != nil {
			return err
		}

		html := buf.String()

		page := Page{
			Title:   title,
			Date:    date,
			Content: template.HTML(html),
			Params:  params,
		}

		templateData := TemplateData{
			Site: site,
			Page: page,
		}

		relPath, _ := filepath.Rel(contentDir, path)
		outDir := filepath.Join(publicDir, strings.TrimSuffix(relPath, ".md"))
		outputPath := filepath.Join(outDir, "index.html")

		os.MkdirAll(outDir, 0755)
		outFile, _ := os.Create(outputPath)
		defer outFile.Close()

		if err := tpl.Execute(outFile, templateData); err != nil {
			return err
		}

		fmt.Println("Generated:", outputPath)
		return nil
	})
}

func copyStaticFiles(src, dest string) {
	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(src, path)
		outPath := filepath.Join(dest, relPath)
		os.MkdirAll(filepath.Dir(outPath), 0755)

		from, _ := os.Open(path)
		defer from.Close()
		to, _ := os.Create(outPath)
		defer to.Close()
		io.Copy(to, from)

		return nil
	})
}
