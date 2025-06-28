package spec

// Contact represents contact information for the API.
type Contact struct {
	Name  string `yaml:"name,omitempty" json:"name,omitempty"`
	Email string `yaml:"email,omitempty" json:"email,omitempty"`
}

// License represents license information for the API.
type License struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
	URL  string `yaml:"url,omitempty" json:"url,omitempty"`
}

// Info holds general API information.
type Info struct {
	Title       string  `yaml:"title" json:"title"`
	Version     string  `yaml:"version" json:"version"`
	Description string  `yaml:"description" json:"description"`
	Contact     Contact `yaml:"contact,omitempty" json:"contact,omitempty"`
	License     License `yaml:"license,omitempty" json:"license,omitempty"`
}

// Spec is the root of the WebSocket API documentation.
type Spec struct {
	Info    Info     `yaml:"info" json:"info"`
	Sockets []Socket `yaml:"sockets" json:"sockets"`
}

// Socket represents a WebSocket endpoint.
type Socket struct {
	Name             string            `yaml:"name" json:"name"`
	URL              string            `yaml:"url" json:"url"`
	Description      string            `yaml:"description" json:"description"`
	Group            string            `yaml:"group,omitempty" json:"group,omitempty"`
	Tags             []string          `yaml:"tags,omitempty" json:"tags,omitempty"`
	ConnectionParams []ConnectionParam `yaml:"connectionParams,omitempty" json:"connectionParams,omitempty"`
	Messages         []Message         `yaml:"messages" json:"messages"`
}

// ConnectionParam represents a connection parameter for a WebSocket endpoint.
type ConnectionParam struct {
	Name        string `yaml:"name" json:"name"`
	In          string `yaml:"in" json:"in"` // query, header
	Type        string `yaml:"type" json:"type"`
	Required    bool   `yaml:"required" json:"required"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

// Message represents a message type for a WebSocket endpoint.
type Message struct {
	Type        string      `yaml:"type" json:"type"`
	Direction   string      `yaml:"direction" json:"direction"` // send | receive
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	Payload     interface{} `yaml:"payload" json:"payload"`
	Example     interface{} `yaml:"example,omitempty" json:"example,omitempty"`
	Errors      []Error     `yaml:"errors,omitempty" json:"errors,omitempty"`
	Deprecated  bool        `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Tags        []string    `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// Error represents an error type for a WebSocket endpoint.
type Error struct {
	Code        string      `yaml:"code" json:"code"`
	Description string      `yaml:"description" json:"description"`
	Example     interface{} `yaml:"example,omitempty" json:"example,omitempty"`
}
