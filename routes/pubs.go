package routes

import (
	"net/http"
        "github.com/gin-gonic/gin"
        "b00m.in/gin/handlers"
)

func addPubRoutes(rg *gin.RouterGroup) {
	rg.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pubs")
	})
        rg.GET("/subs", handlers.HandlePubsGET)
}

