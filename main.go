package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	settings "simpleBlog/config"
	"simpleBlog/internal"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"gopkg.in/yaml.v2"
)

type Article struct {
    Title   string
    Author  string
    Content template.HTML
    Path    string
    Date    time.Time
}

type Metadata struct {
    Title  string `yaml:"title"`
    Author string `yaml:"author"`
    Date   string `yaml:"date"`
}

type IndexData struct {
	Articles []Article
	Title   string
    Links []internal.MenuLink
}

func main() {


    // POPULATE PUBLIC FOLDER USING CONFIGURATION FILE  
    internal.InitPublic()

    // CREATE THE DATA TO POPULATE THE MENU IN THE NAVBAR
    linksList := internal.InitMenu()

    // Load templates
    articleTmpl := template.Must(template.ParseFiles("templates/base.html", "templates/article.html", "templates/navbar.html", "templates/footer.html", "templates/head.html"))
    indexTmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html", "templates/navbar.html", "templates/footer.html", "templates/head.html"))

    // Read articles
    files, err := os.ReadDir("Articles")
    if err != nil {
        log.Fatal(err)
    }

    test := filepath.Base("public/articles")    
    var articles []Article
    for _, file := range files {
        if filepath.Ext(file.Name()) == ".md" {
            content, err := os.ReadFile(filepath.Join("Articles", file.Name()))
            if err != nil {
                log.Fatal(err)
            }

            metadata, body := extractMetadata(content)
            htmlContent := template.HTML(markdown.ToHTML([]byte(body), nil, nil))
			articlePath := strings.TrimSuffix(file.Name(), ".md") + ".html"
            article := Article{
                Title:   metadata.Title,
                Author:  metadata.Author,
                Content: htmlContent,
                Path:    "/" + test + "/" + articlePath,
                Date:    parseDate(metadata.Date),
            }
            articles = append(articles, article) 

            // Generate HTML file for each article
            f, err := os.Create(filepath.Join("public", article.Path))
            if err != nil {
                log.Fatal(err)
            }
            articleTmpl.Execute(f, article )
            f.Close()
        }
    }
    
    // Generate index.html
    indexData := IndexData {
        Articles: articles,
        Title:    settings.BlogName,
        Links: linksList,
    }
    f, err := os.Create(filepath.Join("public", "index.html"))
    if err != nil {
        log.Fatal(err)
    }
    indexTmpl.Execute(f, indexData)
    f.Close()

	staticDir := "public"		
	cssDir := "static"
	cs := http.FileServer(http.Dir(cssDir))
	fs := http.FileServer(http.Dir(staticDir))

	http.Handle("/static/", http.StripPrefix("/static/", cs))
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        path := filepath.Join(staticDir, r.URL.Path)
        if _, err := filepath.Abs(path); err != nil {
            http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
            return
        }

        http.StripPrefix("/", fs).ServeHTTP(w, r)
    })

	http.ListenAndServe(":8080", nil)
}

func extractMetadata(content []byte) (Metadata, string) {
    parts := strings.SplitN(string(content), "---", 3)
    var metadata Metadata
    yaml.Unmarshal([]byte(parts[1]), &metadata)
    return metadata, parts[2]
}

func parseDate(dateStr string) time.Time {
    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        log.Fatal(err)
    }
    return date
}
