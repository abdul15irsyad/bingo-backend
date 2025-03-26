package handler

import (
	"bingo/lib"
	"bingo/model"
	"bingo/service"
	"bingo/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *SocketHandler) GameHandler(c *gin.Context) {
	authUser, ok := c.Get("authUser")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	totalPlayerQuery := c.Query("total-player")
	totalPlayer, err := strconv.Atoi(totalPlayerQuery)
	if err != nil {
		totalPlayer = 2
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
		client.Conn.Close()

		return nil
	})

	currentTime := time.Now()
	fmt.Printf("%s: %s connected\n", currentTime.Format("2006-01-02 15:04:05"), client.User.Name)

	h.socketService.Mutex.Lock()
	readyOpponents := util.FilterSlice(&h.socketService.Queues, func(queue *model.Queue) bool {
		return queue.GameTotalPlayer == totalPlayer
	})
	if len(readyOpponents)+1 < totalPlayer {
		// add to queue
		newUuid, _ := uuid.NewRandom()
		h.socketService.Queues = append(h.socketService.Queues, model.Queue{
			Id:              newUuid,
			GameTotalPlayer: totalPlayer,
			Client:          client,
			CreatedAt:       time.Now(),
		})
	} else {
		// start game
		players := util.FilterSlice(&h.socketService.Queues, func(queue *model.Queue) bool {
			return queue.GameTotalPlayer == totalPlayer && queue.Client.User.Id != client.User.Id
		})
		newUuid, _ := uuid.NewRandom()
		players = append(players, model.Queue{
			Id:              newUuid,
			GameTotalPlayer: totalPlayer,
			Client:          client,
			CreatedAt:       time.Now(),
		})
		h.socketService.Queues = util.FilterSlice(&h.socketService.Queues, func(queue *model.Queue) bool {
			return queue.GameTotalPlayer != totalPlayer && queue.Client.User.Id != client.User.Id
		})
		game, err := h.gameService.CreateGame(service.CreateGameDTO{
			TotalPlayer: totalPlayer,
			Users: util.MapSlice(players, func(player model.Queue) model.User {
				return *player.Client.User
			}),
		})
		if err != nil {
			c.Error(err)
			return
		}

		room := h.socketService.CreateRoom(&game)
		room.Clients = util.MapSlice(players, func(queue model.Queue) model.Client {
			return *queue.Client
		})

		err = h.socketService.BroadcastToRoom(room, model.Message{
			Type:      "game",
			Content:   "game started",
			CreatedAt: time.Now(),
		})
		if err != nil {
			c.Error(err)
			return
		}
	}
	h.socketService.Mutex.Unlock()

	for {
		var JSONMessage map[string]any
		if err := conn.ReadJSON(&JSONMessage); err != nil {
			fmt.Println(err)
			break
		}
		jsonData, _ := json.Marshal(JSONMessage)
		switch JSONMessage["type"] {
		case "message":
			var payload model.Message
			_ = json.Unmarshal(jsonData, &payload)
			payload.CreatedAt = time.Now()
			payload.User = client.User
			fmt.Println(*payload.User)
		default:
			fmt.Printf("type \"%s\" not found\n", JSONMessage["type"])
		}
	}
}
