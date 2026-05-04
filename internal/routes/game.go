package routes

import (
	"quiz/config/db"
	"quiz/internal/features/game"

	"github.com/gin-gonic/gin"
)

func GameRoutes(rg *gin.RouterGroup) {
	rp := game.NewPostgresRepository(db.DB)
	s := game.NewGameUseCase(rp)
	h := game.NewGameHandler(s)

	rg.GET("/game/:game_id/question", h.GetQuestion)
	rg.POST("/game/answer", h.SubmitAnswer)
	rg.GET("/game/:game_id/scoreboard", h.GetScoreboard)
	rg.POST("/game/:game_id/finish", h.FinishGame)
}
