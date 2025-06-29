<p align="center">
  <img src="https://github.com/muratmirgun/socketeer/blob/main/internal/templates/logo.png" alt="Socketeer Logo" width="180" />
</p>

# Socketeer

**Modern, Swagger-Style WebSocket API Docs & Playground for Go**

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/muratmirgun/socketeer)](https://goreportcard.com/report/github.com/muratmirgun/socketeer)
[![Release](https://img.shields.io/github/v/release/muratmirgun/socketeer)](https://github.com/muratmirgun/socketeer/releases)

<p align="center">
  <img src="https://github.com/muratmirgun/socketeer/blob/main/banner.png" alt="Socketeer UI Showcase" width="100%" />
  <br/>
  <em>Modern, interactive WebSocket API documentation and playground UI</em>
</p>

Socketeer is an open-source tool that generates interactive, Swagger-like documentation and playgrounds for your WebSocket APIs in Go.  
It parses special annotations in your Go code and produces a `wsapi.yaml` spec, which is visualized in a beautiful, build-free frontend.

---

## 📑 Table of Contents

- [Features](#-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [CLI Commands](#-cli-commands)
- [Annotation Reference](#-annotation-reference)
- [Available Annotations](#️-available-annotations)
- [Advanced Usage](#-advanced-usage)
  - [Gin Middleware Integration](#gin-middleware-integration)
- [Development](#-development)
- [How It Works](#-how-it-works)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🚀 Features

- **Swagger-style API info annotations** (title, version, description, contact, license)
- **Parse Go code for custom WebSocket annotations** (`@WebSocket`, `@Message`, `@Payload`, `@Group`, etc.)
- **Struct-based payload support** (`@Payload MyStruct` or `@Payload dto.MyStruct`)
- **Generate `wsapi.yaml` or JSON spec**
- **Serve docs and playground via HTTP** (no build step required)
- **Multi-client playground** (test with multiple virtual clients in one UI)
- **Modern, responsive UI** (Swagger-inspired, with live playground)
- **Cobra-powered CLI** (`init`, `generate`, `serve`, `version`)
- **MIT licensed, easy to extend**

---

## 📦 Installation

### From GitHub Releases

Download the latest release for your platform from [GitHub Releases](https://github.com/muratmirgun/socketeer/releases).

### Using Go

```sh
go install github.com/muratmirgun/socketeer@latest
```

---

## 🏃‍♂️ Quick Start

```sh
# Install socketeer
go install github.com/muratmirgun/socketeer@latest

# Navigate to your Go project
cd your-go-project

# Initialize socketeer project
socketeer init

# Add annotations to your Go code (see below)
socketeer generate --src ./ --out ./wsdocs/wsapi.yaml

# Serve the documentation
socketeer serve

# Open http://localhost:8080 in your browser
```

---

## 📋 CLI Commands

### `socketeer init`
Initialize a new socketeer project by creating a `wsdocs/` directory with base files.

```sh
socketeer init
```

### `socketeer generate`
Generate `wsapi.yaml` spec from Go source annotations.

```sh
# Basic usage
socketeer generate

# With custom source and output
socketeer generate --src ./internal --out ./docs/wsapi.yaml

# Available flags:
#   --src string   Source directory to scan for Go files (default "./")
#   --out string   Output spec file (YAML) (default "wsdocs/wsapi.yaml")
```

### `socketeer serve`
Serve documentation and playground from a directory.

```sh
# Basic usage
socketeer serve

# With custom directory and port
socketeer serve --dir ./docs
SOCKETEER_PORT=3000 socketeer serve

# Available flags:
#   --dir string   Directory to serve static files from (default "wsdocs")
# Environment variables:
#   SOCKETEER_PORT  Port to serve on (default "8080")
```

### `socketeer version`
Show socketeer version information.

```sh
socketeer version
```

---

## 📝 Annotation Reference

### API Info Annotations

Add these annotations above your `main` function or at the top of your main Go file:

```go
// @title Socketeer WebSocket API Docs
// @version 1.0.0
// @description Real-time WebSocket API documentation with interactive playground
// @contact.name Murat Mirgun
// @contact.email murat@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
func main() {
    // ...
}
```

### WebSocket Endpoint Annotations

```go
// @WebSocket CompanySocket
// @Group Company Management
// @URL ws://localhost:8080/ws/company
// @Description Company management WebSocket channel for real-time operations
// @Tags company, admin, real-time

// @ConnectionParam name query string required User name for authentication
// @ConnectionParam token header string required JWT token

// @Message addCompany
// @Send
// @Description Add a new company to the system
// @Payload dto.ReqAddCompany
// @Tags company, create

// @Message companyAdded
// @Receive
// @Description Company added successfully
// @Payload dto.ResAddCompany
// @Tags company, response

// @Message updateCompany
// @Send
// @Description Update an existing company
// @Payload dto.ReqUpdateCompany
// @Tags company, update

// @Message companyUpdated
// @Receive
// @Description Company updated successfully
// @Payload dto.ResUpdateCompany
// @Tags company, response

// @Message deleteCompany
// @Send
// @Description Delete a company
// @Payload {"id": "string"}
// @Tags company, delete

// @Message companyDeleted
// @Receive
// @Description Company deleted successfully
// @Payload {"id": "string", "success": true}
// @Tags company, response

// @Message companyError
// @Receive
// @Description Error response for company operations
// @Payload {"error": "string", "code": "string"}
// @Error 400 Bad Request
// @Error 404 Company not found
// @Error 500 Internal Server Error
// @Tags company, error

func CompanySocketHandler(c *gin.Context) {
    // WebSocket handler implementation
}
```

### Struct-based Payload Example

```go
package dto

// ReqAddCompany represents a company creation request
type ReqAddCompany struct {
    // Name of the company
    Name string `json:"name" validate:"required,min=2,max=100,alpha_space" Example:"Acme Inc"`
    // Status of the company (1: active, 0: inactive)
    Status int64 `json:"status" validate:"required" Example:1`
    // Company description
    Description string `json:"description" Example:"Leading technology company"`
    // Founded year
    FoundedYear int `json:"founded_year" Example:2020`
}

// ResAddCompany represents a company creation response
type ResAddCompany struct {
    // Company ID
    ID string `json:"id" Example:"comp_123456"`
    // Company name
    Name string `json:"name" Example:"Acme Inc"`
    // Creation timestamp
    CreatedAt string `json:"created_at" Example:"2024-01-15T10:30:00Z"`
    // Success status
    Success bool `json:"success" Example:true`
}
```

---

## 🏷️ Available Annotations

### API Info Annotations
| Annotation | Description | Example |
|------------|-------------|---------|
| `@title` | API title | `@title My WebSocket API` |
| `@version` | API version | `@version 1.0.0` |
| `@description` | API description | `@description Real-time API for chat` |
| `@contact.name` | Contact name | `@contact.name John Doe` |
| `@contact.email` | Contact email | `@contact.email john@example.com` |
| `@license.name` | License name | `@license.name MIT` |
| `@license.url` | License URL | `@license.url https://opensource.org/licenses/MIT` |

### WebSocket Annotations
| Annotation | Description | Example |
|------------|-------------|---------|
| `@WebSocket` | WebSocket name | `@WebSocket ChatSocket` |
| `@Group` | Group name | `@Group Chat Management` |
| `@URL` | WebSocket URL | `@URL ws://localhost:8080/ws/chat` |
| `@Description` | WebSocket description | `@Description Real-time chat functionality` |
| `@Tags` | Tags for categorization | `@Tags chat, real-time, messaging` |
| `@ConnectionParam` | Connection parameters | `@ConnectionParam token header string required JWT token` |

### Message Annotations
| Annotation | Description | Example |
|------------|-------------|---------|
| `@Message` | Message type/name | `@Message sendMessage` |
| `@Send` | Send direction | `@Send` |
| `@Receive` | Receive direction | `@Receive` |
| `@Payload` | Message payload | `@Payload dto.ChatMessage` |
| `@Error` | Error response | `@Error 400 Bad Request` |
| `@Deprecated` | Mark as deprecated | `@Deprecated` |

### Struct Field Annotations
| Annotation | Description | Example |
|------------|-------------|---------|
| `Example:` | Field example value | `// Example: "Hello World"` |

---

## 🔧 Advanced Usage

### Custom Port Configuration

```sh
# Set custom port via environment variable
export SOCKETEER_PORT=3000
socketeer serve

# Or inline
SOCKETEER_PORT=3000 socketeer serve
```

### Multiple WebSocket Endpoints

```go
// Chat WebSocket
// @WebSocket ChatSocket
// @Group Chat
// @URL ws://localhost:8080/ws/chat
// @Description Real-time chat functionality
// @Tags chat, messaging

// @Message sendMessage
// @Send
// @Description Send a chat message
// @Payload dto.ChatMessage

// @Message messageReceived
// @Receive
// @Description Message received from another user
// @Payload dto.ChatMessage

func ChatSocketHandler(c *gin.Context) {
    // Chat handler
}

// Notification WebSocket
// @WebSocket NotificationSocket
// @Group Notifications
// @URL ws://localhost:8080/ws/notifications
// @Description Real-time notifications
// @Tags notifications, alerts

// @Message subscribe
// @Send
// @Description Subscribe to notifications
// @Payload {"user_id": "string"}

// @Message notification
// @Receive
// @Description Receive notification
// @Payload dto.Notification

func NotificationSocketHandler(c *gin.Context) {
    // Notification handler
}
```

### Error Handling

```go
// @Message userAction
// @Send
// @Description Perform user action
// @Payload dto.UserAction

// @Message actionResult
// @Receive
// @Description Action result
// @Payload dto.ActionResult

// @Message actionError
// @Receive
// @Description Action error
// @Payload {"error": "string", "code": "string"}
// @Error 400 Invalid input
// @Error 401 Unauthorized
// @Error 403 Forbidden
// @Error 404 Not found
// @Error 500 Internal server error
```

### Gin Middleware Integration

You can serve Socketeer documentation and playground directly from your Gin application using the built-in middleware:

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/muratmirgun/socketeer/pkg/socketeer"
)

func main() {
    r := gin.Default()

    // Serve Socketeer docs at /docs (default: wsdocs directory)
    r.Use(socketeer.GinMiddleware(&socketeer.Config{
        Path:       "/docs",              // URL path to serve docs
        StaticPath: "wsdocs",             // Directory for static files (index.html, logo.png, etc.)
        SpecPath:   "wsdocs/wsapi.yaml",  // Path to wsapi.yaml
        EnableCORS: true,                  // Enable CORS headers (optional)
    }))

    // Your WebSocket and API routes here
    r.GET("/ws/company", CompanySocketHandler)

    r.Run(":8080")
}
```

- Artık http://localhost:8080/docs adresinden Socketeer arayüzüne erişebilirsiniz.
- `socketeer.GinMiddleware(nil)` ile varsayılan ayarları da kullanabilirsiniz.

---

## 🛠️ Development

### Building from Source

```sh
git clone https://github.com/muratmirgun/socketeer.git
cd socketeer
go build -o socketeer .
```

### Running Tests

```sh
go test -v ./...
```

### Linting

```sh
golangci-lint run
```

---

## 📖 How It Works

1. **Annotate** your Go code with Swagger-style and WebSocket-specific comments
2. **Generate** the spec: `socketeer generate --src ./ --out ./wsdocs/wsapi.yaml`
3. **Serve** the docs and playground: `socketeer serve`
4. **Explore and test** your WebSocket API in the browser with a modern, interactive UI

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

## 📄 License

MIT License - see the [LICENSE](LICENSE) file for details.
