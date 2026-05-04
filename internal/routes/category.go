package routes

import (
	"quiz/config/db"
	"quiz/internal/features/category"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(rg *gin.RouterGroup) {
	rp := category.NewPostgresRepository(db.DB)
	s := category.NewCategoryUseCase(rp)
	h := category.NewCategoryHandler(s)

	rg.GET("/categories", h.GetAll)
}
