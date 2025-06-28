package middleware

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// SocketeerConfig holds configuration for the Socketeer middleware
type SocketeerConfig struct {
	// Path where the documentation will be served (default: "/docs")
	Path string
	// Path to the wsapi.yaml file (default: "./wsdocs/wsapi.yaml")
	SpecPath string
	// Path to the static files (default: "./wsdocs")
	StaticPath string
	// Title for the documentation page
	Title string
	// Whether to enable CORS for the documentation
	EnableCORS bool
}

// DefaultSocketeerConfig returns default configuration
func DefaultSocketeerConfig() *SocketeerConfig {
	return &SocketeerConfig{
		Path:       "/docs",
		SpecPath:   "./wsdocs/wsapi.yaml",
		StaticPath: "./wsdocs",
		Title:      "WebSocket API Documentation",
		EnableCORS: true,
	}
}

// Socketeer returns a Gin middleware that serves Socketeer documentation
func Socketeer(config *SocketeerConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultSocketeerConfig()
	}

	// Ensure paths end with slash for proper routing
	docsPath := strings.TrimSuffix(config.Path, "/")
	
	return func(c *gin.Context) {
		// Handle CORS if enabled
		if config.EnableCORS {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		}

		// Serve the main documentation page
		if c.Request.URL.Path == docsPath || c.Request.URL.Path == docsPath+"/" {
			serveDocsPage(c, config)
			return
		}

		// Serve the wsapi.yaml spec file at /docs/wsapi.yaml
		if c.Request.URL.Path == docsPath+"/wsapi.yaml" {
			serveSpecFile(c, config.SpecPath)
			return
		}

		// Serve the wsapi.yaml spec file at /wsapi.yaml (root)
		if c.Request.URL.Path == "/wsapi.yaml" {
			serveSpecFile(c, config.SpecPath)
			return
		}

		// Serve static files (logo.png, etc.)
		if strings.HasPrefix(c.Request.URL.Path, docsPath+"/") {
			serveStaticFile(c, config.StaticPath, strings.TrimPrefix(c.Request.URL.Path, docsPath+"/"))
			return
		}

		// Continue to next middleware/handler
		c.Next()
	}
}

// serveDocsPage serves the main documentation HTML page
func serveDocsPage(c *gin.Context, config *SocketeerConfig) {
	// Read the index.html template
	indexPath := filepath.Join(config.StaticPath, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Documentation template not found: %v", err),
			"hint":  "Run 'socketeer init' to create documentation files",
		})
		return
	}

	// Set content type
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}

// serveSpecFile serves the wsapi.yaml specification file
func serveSpecFile(c *gin.Context, specPath string) {
	// Check if spec file exists
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "WebSocket API specification not found",
			"hint":  "Run 'socketeer generate' to create the specification file",
		})
		return
	}

	// Read and serve the spec file
	content, err := os.ReadFile(specPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to read specification: %v", err),
		})
		return
	}

	// Set content type for YAML
	c.Header("Content-Type", "application/x-yaml")
	c.Data(http.StatusOK, "application/x-yaml", content)
}

// serveStaticFile serves static files like logo.png
func serveStaticFile(c *gin.Context, staticPath, filename string) {
	filePath := filepath.Join(staticPath, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("File not found: %s", filename),
		})
		return
	}

	// Determine content type based on file extension
	//contentType := "application/octet-stream"
	//switch filepath.Ext(filename) {
	//case ".png":
	//	contentType = "image/png"
	//case ".jpg", ".jpeg":
	//	contentType = "image/jpeg"
	//case ".gif":
	//	contentType = "image/gif"
	//case ".svg":
	//	contentType = "image/svg+xml"
	//case ".css":
	//	contentType = "text/css"
	//case ".js":
	//	contentType = "application/javascript"
	//case ".json":
	//	contentType = "application/json"
	//}

	// Serve the file
	c.File(filePath)
}

// SocketeerHandler returns a Gin handler function for serving Socketeer docs
// This is a convenience function that creates a new router with Socketeer middleware
func SocketeerHandler(config *SocketeerConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultSocketeerConfig()
	}

	// Create a sub-router for documentation
	docsRouter := gin.New()
	docsRouter.Use(Socketeer(config))

	return func(c *gin.Context) {
		docsRouter.ServeHTTP(c.Writer, c.Request)
	}
} 