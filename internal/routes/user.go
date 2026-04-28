package routes

import (
	"quiz/config/db"
	"quiz/internal/features/user"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup) {
	rp := user.NewPostgresRepository(db.DB)
	s := user.NewUserUseCase(rp)
	h := user.NewUserHandler(s)

	rg.POST("/user", h.Register)
	rg.GET("/user/:id", h.GetByID)
}
