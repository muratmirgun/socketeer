package socketeer

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// GinMiddleware returns a Gin middleware that serves Socketeer documentation
func GinMiddleware(config *Config) gin.HandlerFunc {
	if config == nil {
		config = DefaultConfig()
	}

	docsPath := strings.TrimSuffix(config.Path, "/")

	return func(c *gin.Context) {
		if config.EnableCORS {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		}

		if c.Request.URL.Path == docsPath || c.Request.URL.Path == docsPath+"/" {
			serveDocsPageGin(c, config)
			return
		}
		if c.Request.URL.Path == docsPath+"/wsapi.yaml" {
			serveSpecFileGin(c, config.SpecPath)
			return
		}
		if c.Request.URL.Path == "/wsapi.yaml" {
			serveSpecFileGin(c, config.SpecPath)
			return
		}
		if strings.HasPrefix(c.Request.URL.Path, docsPath+"/") {
			serveStaticFileGin(c, config.StaticPath, strings.TrimPrefix(c.Request.URL.Path, docsPath+"/"))
			return
		}
		c.Next()
	}
}

func serveDocsPageGin(c *gin.Context, config *Config) {
	indexPath := filepath.Join(config.StaticPath, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Documentation template not found: %v", err),
			"hint":  "Run 'socketeer init' to create documentation files",
		})
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}

func serveSpecFileGin(c *gin.Context, specPath string) {
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "WebSocket API specification not found",
			"hint":  "Run 'socketeer generate' to create the specification file",
		})
		return
	}
	content, err := os.ReadFile(specPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to read specification: %v", err),
		})
		return
	}
	c.Header("Content-Type", "application/x-yaml")
	c.Data(http.StatusOK, "application/x-yaml", content)
}

func serveStaticFileGin(c *gin.Context, staticPath, filename string) {
	filePath := filepath.Join(staticPath, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("File not found: %s", filename),
		})
		return
	}
	c.File(filePath)
} 