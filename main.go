package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"quiz/config/db"
	"quiz/internal/middleware"
	"quiz/internal/routes"
	"quiz/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var ctx = context.Background()

var allowedOrigins map[string]bool

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return allowedOrigins[origin]
	},
}

func wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		log.Printf("mensaje: %s\n", message)

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println("write error:", err)
			break
		}
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Advertencia: no se encontró .env, se usan variables de entorno del sistema")
	}

	if err := db.Connect(); err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	if err := db.InitializeDatabase(); err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	ip := utils.GetOutboundIP()
	log.Printf("IP del servidor: %s\n", ip)

	HOST_URL_DEV := os.Getenv("HOST_URL_DEV")
	HOST_URL_PROD := os.Getenv("HOST_URL_PROD")
	HOST_URL_PROD_WWW := os.Getenv("HOST_URL_PROD_WWW")

	allowedOrigins = map[string]bool{
		HOST_URL_DEV:      true,
		HOST_URL_PROD:     true,
		HOST_URL_PROD_WWW: true,
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{HOST_URL_DEV, HOST_URL_PROD, HOST_URL_PROD_WWW},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWebSockets:  true,
	}))

	r.Use(middleware.RateLimiterMiddleware(rate.Every(time.Minute/10), 10))
	r.GET("/ws", wsHandler)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the CV API"})
	})

	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK", "ip": ip})
	})

	api := r.Group("/api/v1")
	{
		routes.UserRoutes(api)
	}

	var wg sync.WaitGroup
	cleanupCtx, cancelCleanup := context.WithCancel(context.Background())
	wg.Go(func() {
		middleware.StartCleanup(cleanupCtx)
	})

	srv := &http.Server{
		Addr:    "0.0.0.0:4100",
		Handler: r,
	}

	go func() {
		log.Println("Server starting on :4100...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	cancelCleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	wg.Wait()
	log.Println("Server exiting")
}
