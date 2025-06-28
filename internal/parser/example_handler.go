package parser

import "github.com/gin-gonic/gin"

// @WebSocket ChatSocket
// @Group Messaging
// @URL ws://localhost:8080/ws/chat
// @Description Kullanıcıların mesajlaştığı WebSocket kanalı
// @Tags chat, public
// @ConnectionParam name query string required Kullanıcı adı

// @Message sendMessage
// @Direction send
// @Description Kullanıcıdan mesaj gönderme
// @Payload
// {
//   "text": "Merhaba"
// }
// @Example
// {
//   "text": "Merhaba"
// }
// @Error unauthorized Kullanıcı yetkisiz
// @Error invalid_format Geçersiz mesaj formatı

// @Message messageReceived
// @Direction receive
// @Description Sunucudan gelen mesaj
// @Payload
//
//	{
//	  "user": "Ali",
//	  "text": "Merhaba"
//	}
//
// @Example
//
//	{
//	  "user": "Ali",
//	  "text": "Merhaba"
//	}
func ChatSocketHandler(c *gin.Context) {
	// WebSocket bağlantısı
}

// @WebSocket NotificationSocket
// @Group Notification
// @URL ws://localhost:8080/ws/notify
// @Description Gerçek zamanlı bildirim kanalı
// @Tags notification, private
// @ConnectionParam session_id query string required Oturum ID

// @Message subscribe
// @Direction send
// @Description Bildirim kanalına abone ol
// @Payload
// {
//   "topics": ["news", "alerts"]
// }
// @Example
// {
//   "topics": ["news", "alerts"]
// }

// @Message notification
// @Direction receive
// @Description Sunucudan gelen bildirim
// @Payload
//
//	{
//	  "topic": "news",
//	  "content": "Yeni haber var!"
//	}
//
// @Example
//
//	{
//	  "topic": "news",
//	  "content": "Yeni haber var!"
//	}
func NotificationSocketHandler(c *gin.Context) {
	// WebSocket bildirim kanalı
}

// @WebSocket PingPongSocket
// @URL /ws/ping
// @Description Ping-pong test endpoint. Ping gönder, pong cevabı al.
// @Tags test, pingpong

// @Message ping
// @Direction send
// @Description Ping mesajı gönder
// @Payload
// {
//   "type": "ping",
//   "payload": "merhaba"
// }
// @Example
// {
//   "type": "ping",
//   "payload": "merhaba"
// }

// @Message pong
// @Direction receive
// @Description Pong cevabı
// @Payload
//
//	{
//	  "type": "pong",
//	  "payload": "merhaba"
//	}
//
// @Example
//
//	{
//	  "type": "pong",
//	  "payload": "merhaba"
//	}
func PingPongSocketHandler(c *gin.Context) {
	// WebSocket ping-pong
}

// DTO örneği (normalde başka bir dosyada/pakette olurdu)
type ReqAddCompany struct {
	// Name of the company
	// in: string
	Name string `json:"name" validate:"required,min=2,max=100,alpha_space"`
	// Status of the company
	// in: int64
	Status int64 `json:"status" validate:"required"`
}

// @WebSocket CompanySocket
// @Group Company
// @URL ws://localhost:8080/ws/company
// @Description Şirket ekleme WebSocket kanalı
// @Tags company, admin
// @ConnectionParam name query string required Kullanıcı adı

// @Message addCompany
// @Direction send
// @Description Yeni şirket ekle
// @Payload ReqAddCompany
// @Example
// {
//   "name": "Acme Inc",
//   "status": 1
// }

// @Message companyAdded
// @Direction receive
// @Description Şirket başarıyla eklendi
// @Payload ReqAddCompany
func CompanySocketHandler(c *gin.Context) {
	// WebSocket şirket ekleme
}
