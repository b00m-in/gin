package routes

import (
        "github.com/gin-gonic/gin"
        "b00m.in/gin/handlers"
        //"b00m.in/data"
)

func addSubRoutes(rg *gin.RouterGroup) {
        /*var authChain = []data.Middleware {
                data.AuthMiddleware,
        }
        rg.GET("/", gin.WrapF(data.BuildChain(handlers.HandleSubsGet, authChain...)))
        rg.POST("/", gin.WrapF(data.BuildChain(handlers.HandleSubsGet, authChain...)))*/
	/*rg.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "subs")
	})*/
        rg.GET("/", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/pubs", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/pubs/:hash", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/packets/:hash", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.POST("/packets/:hash", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/faults", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/subs", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/login", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/logout", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/register", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.POST("/login", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.POST("/new", handlers.ReadCookie(), handlers.HandleSubsGET)
}

