package handlers

import (
	"fmt"
	"lru-cache-gin/models"
	"lru-cache-gin/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetHandler(ctx *gin.Context, cache *models.LRUCache, hub *utils.WebSocketHub) {
	key := ctx.Param("key")

	if value, found := cache.Get(key, hub); found {
		ctx.JSON(http.StatusOK, gin.H{"value": value})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
	}
}

func SetHandler(ctx *gin.Context, cache *models.LRUCache, hub *utils.WebSocketHub) {
	var data struct {
		Key    string        `json:"key"`
		Value  interface{}   `json:"value"`
		Expiry time.Duration `json:"expiry"` // in seconds
	}

	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	cache.Set(data.Key, data.Value, time.Second*data.Expiry, hub)
	fmt.Print(time.Second*data.Expiry)
	ctx.Status(http.StatusOK)
}

func DeleteHandler(ctx *gin.Context, cache *models.LRUCache, hub *utils.WebSocketHub) {
	key := ctx.Param("key")

	cache.Delete(key, hub)
	ctx.Status(http.StatusOK)
}

