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
        rg.POST("/pubs/:hash", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/packets/:hash", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.POST("/packets/:hash", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/faults", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/faults/map", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/summary", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/you", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.POST("/you", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/login", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.POST("/login", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/logout", handlers.ReadCookie(), handlers.HandleSubsLogout)
        rg.GET("/register", handlers.ReadCookie(), handlers.HandleSubsRegister)
        rg.POST("/register", handlers.ReadCookie(), handlers.HandleSubsRegister)
        rg.GET("/verify/:hash", handlers.ReadCookie(), handlers.HandleSubsRegister)
        rg.GET("/help", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/groups", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/connections", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/leaderboard", handlers.ReadCookie(), handlers.HandleSubsGET)
        rg.GET("/messages", handlers.ReadCookie(), handlers.HandleSubsGET)
}

