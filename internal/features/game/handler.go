package game

import (
	"encoding/json"
	"net/http"
	"quiz/internal/shared/hub"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	service GameService
}

func NewGameHandler(service GameService) *GameHandler {
	return &GameHandler{service: service}
}

type SubmitAnswerInput struct {
	DepartureID  uint64 `json:"departure_id"`
	QuestionID   uint64 `json:"question_id"`
	AnswerID     uint64 `json:"answer_id"`
	ResponseTime int    `json:"response_time"`
	GameKey      string `json:"game_key"`
}

func (h *GameHandler) GetQuestion(c *gin.Context) {
	idStr := c.Param("game_id")
	gameID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "game_id invalido"})
		return
	}

	question, err := h.service.GetQuestion(c.Request.Context(), gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": question})
}

func (h *GameHandler) SubmitAnswer(c *gin.Context) {
	var input SubmitAnswerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input: " + err.Error()})
		return
	}

	result, err := h.service.SubmitAnswer(c.Request.Context(), input.DepartureID, input.QuestionID, input.AnswerID, input.ResponseTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if input.GameKey != "" {
		scoreboard, _ := h.service.GetScoreboard(c.Request.Context(), result.Departure.GameID)
		msg, _ := json.Marshal(gin.H{"event": "score_update", "data": scoreboard})
		hub.Global.Broadcast(input.GameKey, msg)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *GameHandler) GetScoreboard(c *gin.Context) {
	idStr := c.Param("game_id")
	gameID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "game_id invalido"})
		return
	}

	scoreboard, err := h.service.GetScoreboard(c.Request.Context(), gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": scoreboard})
}

func (h *GameHandler) FinishGame(c *gin.Context) {
	idStr := c.Param("game_id")
	gameID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "game_id invalido"})
		return
	}

	gameKey := c.Query("key")

	if err := h.service.FinishGame(c.Request.Context(), gameID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	scoreboard, _ := h.service.GetScoreboard(c.Request.Context(), gameID)
	if gameKey != "" {
		msg, _ := json.Marshal(gin.H{"event": "game_finished", "data": scoreboard})
		hub.Global.Broadcast(gameKey, msg)
	}

	c.JSON(http.StatusOK, gin.H{"data": scoreboard, "message": "game finished"})
}
