package routes

import (
	"quiz/config/db"
	"quiz/internal/features/room"

	"github.com/gin-gonic/gin"
)

func RoomRoutes(rg *gin.RouterGroup) {
	rp := room.NewPostgresRepository(db.DB)
	s := room.NewRoomUseCase(rp)
	h := room.NewRoomHandler(s)

	rg.POST("/room", h.CreateRoom)
	rg.POST("/room/join", h.JoinRoom)
	rg.POST("/room/:id/start", h.StartGame)
}

func WsRoutes(r *gin.Engine) {
	rp := room.NewPostgresRepository(db.DB)
	s := room.NewRoomUseCase(rp)
	h := room.NewRoomHandler(s)

	r.GET("/ws/:key", h.WsHandler)
}
