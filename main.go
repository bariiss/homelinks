package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var (
	configPath = filepath.Join(os.Getenv("HOME"), ".config", "homelinks.json")
	links      []*Link
)

// Link represents a hyperlink with styling and text.
type Link struct {
	Name      string `json:"name"`
	Text      string `json:"text"`
	URL       string `json:"url"`
	Color     string `json:"color"`
	TextColor string `json:"textColor"`
	AltText   string `json:"altText"`
}

func init() {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Sample links
		links = []*Link{
			{"google", "Google", "https://www.google.com", "#4285F4", "white", "Search with Google"},
			{"facebook", "Facebook", "https://www.facebook.com", "#F48042", "black", "Connect with Facebook"},
			{"twitter", "Twitter", "https://www.twitter.com", "#F44242", "white", "Tweet with Twitter"},
		}
	} else {
		file, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Error reading the links file: %v", err)
		}

		err = json.Unmarshal(file, &links)
		if err != nil {
			log.Fatalf("Error parsing the links file: %v", err)
		}
	}
}

func main() {
	r := gin.Default()

	setupTmpl(r)
	setupAssets(r)
	setupGlobals(r)
	setupRoutes(r)

	r.Run(":8080")
}

// setupTmpl parses HTML templates from the filesystem.
func setupTmpl(r *gin.Engine) {
	templ, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}
	r.SetHTMLTemplate(templ)
}

// setupAssets serves static assets from the filesystem.
func setupAssets(r *gin.Engine) {
	r.Static("/assets", "./assets")
}

// setupGlobals sets global variables for the router.
func setupGlobals(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Set("bootstrapCss", "/assets/bootstrap/css/bootstrap.min.css")
		c.Set("bootstrapJs", "/assets/bootstrap/js/bootstrap.bundle.min.js")
		c.Set("customCss", "/assets/custom/css/custom.css")
		c.Next()
	})
}

// setupRoutes initializes the routes for the application.
func setupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "links":        links,
            "bootstrapCss": c.MustGet("bootstrapCss"),
            "bootstrapJs":  c.MustGet("bootstrapJs"),
            "customCss":    c.MustGet("customCss"),
        })
	})

	r.NoRoute(func(c *gin.Context) {
        c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
            "bootstrapCss": c.MustGet("bootstrapCss"),
            "bootstrapJs":  c.MustGet("bootstrapJs"),
            "customCss":    c.MustGet("customCss"),
        })
	})
}