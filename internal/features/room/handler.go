package room

import (
	"encoding/json"
	"net/http"
	"quiz/internal/shared/hub"
	"quiz/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type RoomHandler struct {
	service RoomService
}

func NewRoomHandler(service RoomService) *RoomHandler {
	return &RoomHandler{service: service}
}

type CreateRoomInput struct {
	Name       string `json:"name"`
	Nickname   string `json:"nickname"`
	AvatarURL  string `json:"avatar_url"`
	CategoryID uint64 `json:"category_id"`
}

type JoinRoomInput struct {
	Key       string `json:"key"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var input CreateRoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input: " + err.Error()})
		return
	}

	result, err := h.service.CreateRoom(c.Request.Context(), input.Name, input.Nickname, input.AvatarURL, input.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": result})
}

func (h *RoomHandler) JoinRoom(c *gin.Context) {
	var input JoinRoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input: " + err.Error()})
		return
	}

	result, err := h.service.JoinRoom(c.Request.Context(), input.Key, input.Nickname, input.AvatarURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	msg, _ := json.Marshal(gin.H{"event": "player_joined", "data": result})
	hub.Global.Broadcast(input.Key, msg)

	c.JSON(http.StatusCreated, gin.H{"data": result})
}

func (h *RoomHandler) StartGame(c *gin.Context) {
	idStr := c.Param("id")
	gameID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id invalido"})
		return
	}

	game, err := h.service.StartGame(c.Request.Context(), gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	msg, _ := json.Marshal(gin.H{"event": "game_started", "data": game})
	hub.Global.Broadcast(game.Key, msg)

	c.JSON(http.StatusOK, gin.H{"message": "game started"})
}

func (h *RoomHandler) WsHandler(c *gin.Context) {
	key := c.Param("key")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &hub.Client{Conn: conn, RoomKey: key}
	hub.Global.Register(client)
	defer hub.Global.Unregister(client)
	defer conn.Close()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	_, err := utils.ValidateParamsId(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
}
