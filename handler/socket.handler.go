package handler

import (
	"bingo/lib"
	"bingo/model"
	"bingo/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type SocketHandler struct {
	socketService *service.SocketService
	gameService   *service.GameService
}

func NewSocketHandler(socketService *service.SocketService, gameService *service.GameService) *SocketHandler {
	lib.Logger.Info("NewSocketHandler initialized")
	return &SocketHandler{socketService, gameService}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *SocketHandler) StartGameHandler(c *gin.Context) {
	authUser, ok := c.Get("authUser")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}
	defer conn.Close()

	authUserClient := authUser.(model.User)
	client := h.socketService.CreateClient(conn, &authUserClient)

	conn.SetCloseHandler(func(code int, text string) error {
		currentTime := time.Now()
		defer fmt.Printf("%s: %s disconnected\n", currentTime.Format("2006-01-02 15:04:05"), client.User.Name)
		h.socketService.RemoveClient(client)
		return nil
	})

	currentTime := time.Now()
	fmt.Printf("%s: %s connected\n", currentTime.Format("2006-01-02 15:04:05"), client.User.Name)

	for {
		var JSONMessage map[string]any
		if err := conn.ReadJSON(&JSONMessage); err != nil {
			fmt.Println(err)
			break
		}
		switch JSONMessage["type"] {
		case "message":
			jsonData, _ := json.Marshal(JSONMessage)
			var payload service.Message
			_ = json.Unmarshal(jsonData, &payload)
			payload.CreatedAt = time.Now()
			payload.Client = client
			fmt.Println(*payload.Client.User)
		default:
			fmt.Printf("type \"%s\" not found\n", JSONMessage["type"])
		}
	}
}
