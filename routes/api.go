package routes

import (
        "github.com/gin-gonic/gin"
        "b00m.in/gin/handlers"
)

//
func addApiRoutes(rg *gin.RouterGroup) {
	/*rg.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "subs")
	})*/
        rg.GET("/", handlers.HandleApiGET)
        rg.GET("/confo/:devicename/:ssid", handlers.HandleApiGET)
        rg.GET("/pubs", handlers.HandleApiGET)
        //rg.POST("/", handlers.HandleApiPOST)
        rg.POST("/register", handlers.HandleApiRegisterPOST)
        rg.POST("/login", handlers.HandleApiLoginPOST)
        rg.GET("/summary/:freq/:hash/:start/:end", handlers.HandleApiGET)
}


