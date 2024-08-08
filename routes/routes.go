package routes

import (
	"lru-cache-gin/handlers"
	"lru-cache-gin/models"
	"lru-cache-gin/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config {
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"GET", "POST", "DELETE"},
        AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
    }))

	// Initialize the LRU Cache and WebSocket Hub
	cache := models.NewLRUCache(20) // Set the cache size as needed
	hub := utils.NewWebSocketHub()
	go hub.Run()



	// Define routes and associate handlers
	r.GET("/get/:key", func(ctx *gin.Context) { handlers.GetHandler(ctx, cache, hub) })
	r.POST("/set", func(ctx *gin.Context) { handlers.SetHandler(ctx, cache, hub) })
	r.DELETE("/delete/:key", func(ctx *gin.Context) { handlers.DeleteHandler(ctx, cache, hub) })

	// WebSocket route
	r.GET("/webSocket", handlers.WebSocketHandler(hub))

	// Broadcast cache updates periodically
	go func() {
		for {
			time.Sleep(1 * time.Second)
			cache.BroadcastCacheUpdate(hub)
		}
	}()

	return r
}
