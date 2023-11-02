package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*.tmpl
var tmplFS embed.FS

//go:embed assets
var assetsFS embed.FS

// Link represents a hyperlink with styling and text.
type Link struct {
	Name      string
	Text      string
	URL       string
	Color     string
	TextColor string
	AltText   string
}

var links []*Link

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	configPath := filepath.Join(usr.HomeDir, ".config", "homelinks.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
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
	r.Use(basicAuthPubIP())

	setupTmpl(r)
	setupAssets(r)
	setupGlobals(r)
	setupRoutes(r)

	r.Run(":8080")
}

// basicAuthPubIP requires basic authentication for requests from public IP addresses.
func basicAuthPubIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		if clientIP := c.ClientIP(); isPubIP(net.ParseIP(clientIP)) {
			basicAuth(c)
		} else {
			c.Next()
		}
	}
}

// basicAuth prompts for basic authentication.
func basicAuth(c *gin.Context) {
	user, pass, ok := c.Request.BasicAuth()
	if !ok || user != "user" || pass != "pass" {
		c.Header("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Next()
}

// setupTmpl parses HTML templates from the embedded filesystem.
func setupTmpl(r *gin.Engine) {
	templ, err := template.ParseFS(tmplFS, "templates/*.tmpl")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}
	r.SetHTMLTemplate(templ)
}

// setupAssets serves static assets from the embedded filesystem.
func setupAssets(r *gin.Engine) {
	aFS, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		log.Fatal("Failed to create sub filesystem for assets:", err)
	}
	r.StaticFS("/assets", http.FS(aFS))
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
		render(c, "index.tmpl", "My Precious Links")
	})

	r.NoRoute(func(c *gin.Context) {
		render(c, "404.tmpl", "404")
	})
}

// render sends an HTML response using a template.
func render(c *gin.Context, tmplName, title string) {
	c.HTML(http.StatusOK, tmplName, gin.H{
		"title":        title,
		"links":        links,
		"bootstrapCss": c.MustGet("bootstrapCss"),
		"bootstrapJs":  c.MustGet("bootstrapJs"),
		"customCss":    c.MustGet("customCss"),
	})
}

// isPubIP checks if an IP is a public address.
func isPubIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		}
	}
	return true
}