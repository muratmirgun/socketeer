<p align="center">
  <img src="https://github.com/muratmirgun/socketeer/blob/main/internal/templates/logo.png" alt="Socketeer Logo" width="180" />
</p>

# Socketeer

**Modern, Swagger-Style WebSocket API Docs & Playground for Go**

Socketeer is an open-source tool that generates interactive, Swagger-like documentation and playgrounds for your WebSocket APIs in Go.  
It parses special annotations in your Go code and produces a `wsapi.yaml` spec, which is visualized in a beautiful, build-free frontend.

---

## Installation

### From GitHub Releases

Download the latest release for your platform from [GitHub Releases](https://github.com/muratmirgun/socketeer/releases).

### Using Go

```sh
go install github.com/muratmirgun/socketeer@latest
```

### Using Homebrew (macOS/Linux)

```sh
brew install muratmirgun/tap/socketeer
```

### Using Docker

```sh
docker pull ghcr.io/muratmirgun/socketeer:latest
docker run -p 8080:8080 ghcr.io/muratmirgun/socketeer:latest
```

---

## Features

- **Swagger-style API info annotations** (title, version, description, contact, license)
- **Parse Go code for custom WebSocket annotations** (`@WebSocket`, `@Message`, `@Payload`, `@Group`, etc.)
- **Struct-based payload support** (`@Payload MyStruct` or `@Payload dto.MyStruct`)
- **Generate `wsapi.yaml` or JSON spec**
- **Serve docs and playground via HTTP** (no build step required)
- **Multi-client playground** (test with multiple virtual clients in one UI)
- **Modern, responsive UI** (Swagger-inspired, with live playground)
- **Cobra-powered CLI** (`init`, `generate`, `serve`, `validate`, `fmt`, `version`)
- **MIT licensed, easy to extend**

---

## Quick Start

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

## Example: API Info Annotation

Add these annotations above your `main` function (or at the top of your main Go file):

```go
// @title Socketeer WebSocket API Docs
// @version 1.0.0
// @description Real-time WebSocket API documentation
// @contact.name Murat
// @contact.email murat@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
func main() {
    // ...
}
```

---

## Example: WebSocket Endpoint Annotation

```go
// @WebSocket CompanySocket
// @Group Company
// @URL ws://localhost:8080/ws/company
// @Description Company management WebSocket channel
// @Tags company, admin
// @ConnectionParam name query string required User name

// @Message addCompany
// @Direction send
// @Description Add a new company
// @Payload dto.ReqAddCompany
// @Example
// {
//   "name": "Acme Inc",
//   "status": 1
// }

// @Message companyAdded
// @Direction receive
// @Description Company added successfully
// @Payload dto.ReqAddCompany
func CompanySocketHandler(c *gin.Context) {
    // WebSocket handler
}
```

---

## Example: Struct-based Payload

```go
package dto

// ReqAddCompany represents a company creation request
type ReqAddCompany struct {
    // Name of the company
    Name string `json:"name" validate:"required,min=2,max=100,alpha_space"`
    // Status of the company
    Status int64 `json:"status" validate:"required"`
}
```

---

## CLI Usage

```sh
socketeer init                    # Initialize a new socketeer project
socketeer generate --src ./ --out ./wsdocs/wsapi.yaml  # Generate spec from Go code
socketeer serve                   # Serve documentation and playground
socketeer validate                # Validate wsapi.yaml file
socketeer fmt                     # Format wsapi.yaml file
socketeer version                 # Show version information
```

---

## Development

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

## How It Works

1. **Annotate** your Go code with Swagger-style and WebSocket-specific comments.
2. **Generate** the spec:  
   `socketeer generate --src ./ --out ./wsdocs/wsapi.yaml`
3. **Serve** the docs and playground:  
   `socketeer serve`
4. **Explore and test** your WebSocket API in the browser with a modern, interactive UI.

---

## License

MIT
