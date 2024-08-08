package handlers

import (
	"lru-cache-gin/utils"

	"github.com/gin-gonic/gin"
)

func WebSocketHandler(hub *utils.WebSocketHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		hub.ServeWs(c.Writer, c.Request)
	}
}
