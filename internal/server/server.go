package server

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Global chat room state (in-memory, not production safe)
var chatClients = struct {
	sync.Mutex
	clients map[*websocket.Conn]string // conn -> name
}{clients: make(map[*websocket.Conn]string)}

// Broadcast message to all connected chat clients
func broadcastChat(msg interface{}, exclude *websocket.Conn) {
	chatClients.Lock()
	defer chatClients.Unlock()
	for c := range chatClients.clients {
		if c != exclude {
			b, _ := json.Marshal(msg)
			c.WriteMessage(websocket.TextMessage, b)
		}
	}
}

// wsChatHandler: gelişmiş chat odası
func wsChatHandler(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.String(400, "name query param required")
		return
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	// Register client
	chatClients.Lock()
	chatClients.clients[ws] = name
	chatClients.Unlock()
	// Diğer clientlara joined mesajı
	broadcastChat(map[string]string{"user": "system", "text": name + " joined the chat"}, ws)
	// Okuma döngüsü
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		// Mesajı parse et
		var m map[string]interface{}
		if err := json.Unmarshal(msg, &m); err == nil {
			text, _ := m["text"].(string)
			if text != "" {
				// Tüm clientlara ilet
				broadcastChat(map[string]string{"user": name, "text": text}, nil)
			}
		}
	}
	// Bağlantı kopunca unregister ve left mesajı
	chatClients.Lock()
	delete(chatClients.clients, ws)
	chatClients.Unlock()
	broadcastChat(map[string]string{"user": "system", "text": name + " left the chat"}, ws)
}

// StartServer starts the HTTP server for docs, frontend, and WebSocket test endpoints.
func StartServer() {
	r := gin.Default()

	// Serve the API docs as JSON
	r.GET("/api/docs", func(c *gin.Context) {
		data, err := os.ReadFile("./public/wsapi.yaml")
		if err != nil {
			c.JSON(500, gin.H{"error": "Spec not found"})
			return
		}
		var spec interface{}
		if err := yaml.Unmarshal(data, &spec); err != nil {
			c.JSON(500, gin.H{"error": "Invalid YAML"})
			return
		}
		c.JSON(200, spec)
	})

	// WebSocket test endpoints for playground
	r.GET("/ws/chat", wsChatHandler)
	r.GET("/ws/notify", wsEchoHandler)
	r.GET("/ws/ping", wsPingPongHandler)

	// Serve static frontend
	r.Static("/static", "./public")
	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	r.Run(":8080")
}

// wsEchoHandler is a simple echo WebSocket handler for playground testing.
func wsEchoHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		// Echo back the received message
		ws.WriteMessage(mt, msg)
	}
}

// wsPingPongHandler handles ping-pong logic for playground demo.
func wsPingPongHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	defer ws.Close()
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		// Try to parse as JSON
		var m map[string]interface{}
		if err := yaml.Unmarshal(msg, &m); err == nil {
			if t, ok := m["type"]; ok && t == "ping" {
				m["type"] = "pong"
				pong, _ := yaml.Marshal(m)
				ws.WriteMessage(websocket.TextMessage, pong)
				continue
			}
		}
		// Default: echo
		ws.WriteMessage(websocket.TextMessage, msg)
	}
}
