package handler

import (
	"bingo/dto"
	"bingo/lib"
	"bingo/model"
	"bingo/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GameHandler struct {
	socketService *service.SocketService
	gameService   *service.GameService
}

func NewGameHandler(socketService *service.SocketService, gameService *service.GameService) *GameHandler {
	lib.Logger.Info("NewGameHandler initialized")
	return &GameHandler{socketService, gameService}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *GameHandler) Start(c *gin.Context) {
	authUser, ok := c.Get("authUser")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	var startDTO dto.StartDTO
	c.ShouldBind(&startDTO)
	validationErrors := lib.Validate(startDTO)
	if len(validationErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"errors":  validationErrors,
		})
		return
	}
	totalPlayer := startDTO.TotalPlayer

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
		client.Conn.Close()

		return nil
	})

	currentTime := time.Now()
	fmt.Printf("%s: %s connected\n", currentTime.Format("2006-01-02 15:04:05"), client.User.Name)

	h.gameService.QueuePlayer(client, totalPlayer)

	for {
		var JSONMessage map[string]any
		if err := conn.ReadJSON(&JSONMessage); err != nil {
			fmt.Println(err)
			break
		}
		jsonData, _ := json.Marshal(JSONMessage)
		switch JSONMessage["type"] {
		case model.PlayerReadyType:
			var gameReady model.Payload
			_ = json.Unmarshal(jsonData, &gameReady)
			gameId, _ := uuid.Parse(gameReady.Content.(string))
			game := h.gameService.GetGame(gameId)
			if game == nil {
				c.Error(fmt.Errorf("game %s not found", gameId.String()))
				continue
			}
			player := h.gameService.GetPlayerFromUserId(game, client.User.Id)

			isAllReady := h.gameService.PlayerReady(game.Id, player.Id)
			if isAllReady {
				room := h.socketService.GetRoomFromGame(game.Id)
				err := h.socketService.BroadcastToRoom(room, model.Payload{
					Type:      model.GameStartType,
					Content:   nil,
					CreatedAt: time.Now(),
				})
				if err != nil {
					c.Error(err)
					return
				}
			}
		case model.MessageType:
			var message model.Payload
			_ = json.Unmarshal(jsonData, &message)
			message.CreatedAt = time.Now()
			message.User = client.User
			h.socketService.Broadcast(message)
		default:
			fmt.Printf("type \"%s\" not found\n", JSONMessage["type"])
		}
	}
}
