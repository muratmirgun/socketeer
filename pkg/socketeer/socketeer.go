package socketeer

// SocketeerConfig holds configuration for the Socketeer middleware
// This is shared by all framework integrations (Gin, Fiber, Echo, ...)
type Config struct {
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

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Path:       "/docs",
		SpecPath:   "./wsdocs/wsapi.yaml",
		StaticPath: "./wsdocs",
		Title:      "WebSocket API Documentation",
		EnableCORS: true,
	}
} 