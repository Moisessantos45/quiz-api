package user

import (
	"net/http"
	"quiz/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service UserService
}

type RegisterInput struct {
	Nickname string `json:"nickname"`
	GameID   uint64 `json:"game_id"`
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input: " + err.Error()})
		return
	}

	result, err := h.service.Register(input.Nickname, input.GameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": result, "message": "User created successfully"})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := utils.ValidateParamsId(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "message": "User fetched successfully"})
}
