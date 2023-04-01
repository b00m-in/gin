package routes

import (
	"net/http"
        "github.com/gin-gonic/gin"
        "b00m.in/gin/handlers"
        //"b00m.in/data"
)

func addPubRoutes(rg *gin.RouterGroup) {
	rg.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pubs")
	})
        rg.GET("/subs", handlers.HandlePubsGET)
}

func addDataRoutes(rg *gin.RouterGroup) {
        //rg.GET("/publishers", handlers.ReadCookie(), gin.WrapF(data.PubDeetsHandler))
        rg.GET("/publishers", handlers.ReadCookie(), handlers.PubDeetsHandler)
}
