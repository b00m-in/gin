package routes

import (
        //"embed"
	"net/http"
        "github.com/gin-gonic/gin"
        "github.com/gin-contrib/sessions"
        "github.com/gin-contrib/sessions/cookie"
        //"github.com/golang/glog"
)

var db = make(map[string]string)
////go:embed assets/*
//var f embed.FS

func SetupRouter(r *gin.Engine) {
        store := cookie.NewStore([]byte("secretkey"))
        store.Options(sessions.Options{Secure: true})
        sessionnames := []string{"sub", "cart"}
        r.Use(sessions.SessionsMany(sessionnames, store))
        s := r.Group("/subs")
        p := r.Group("/pubs")
        a := r.Group("/api")
        addSubRoutes(s)
        addPubRoutes(p)
        addApiRoutes(a)
        r.StaticFS("/assets", http.FileSystem(http.Dir("assets"))) //http.FS(f))
        r.StaticFS("/static", http.FileSystem(http.Dir("assets"))) //http.FS(f))
        r.StaticFS("/images", http.FileSystem(http.Dir("assets"))) //http.FS(f))
        /*r.GET("favicon.ico", func(c *gin.Context) {
                icon, err := f.ReadFile("assets/favicon.ico")
                if err != nil {
                        glog.Errorf("Readfile: %v\n", err)
                }
                c.Data(http.StatusOK, "image/x-icon", icon)
        })*/
        r.GET("/user/:name", func(c *gin.Context) {
                user := c.Params.ByName("name")
                value, ok := db[user]
                if ok {
                        c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
                }  else {
                        c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
                }
        })
        r.GET("/user/content", func(c *gin.Context) {
                c.HTML(200, "user_content.html", gin.H{"title": "user content"})
        })
        authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
                "foo": "bar",
                "manu": "123",
        }))
        authorized.POST("admin", func(c *gin.Context) {
                user := c.MustGet(gin.AuthUserKey).(string)
                var json struct {
                        Value string `json:"value" binding:"required"`
                }
                if c.Bind(&json) == nil {
                        db[user] = json.Value
                        c.JSON(http.StatusOK, gin.H{"status": "ok"})
                }
        })
}

/*
mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(fileHandler)))
mux.Handle("/api/", http.HandlerFunc(handleAPI))
mux.Handle("/subs/nuevo/", data.BuildChain(handleNewSubs, authChain...))
mux.Handle("/subs/", http.HandlerFunc(handleSubs))
//mux.Handle("/pubs/new", http.HandlerFunc(handlePubs))
mux.Handle("/pubs/", http.HandlerFunc(handlePubs))
mux.Handle("/data/subway-stations", http.HandlerFunc(data.SubwayStationsHandler))
mux.Handle("/data/subway-lines", http.HandlerFunc(subwayLinesHandler))
mux.Handle("/data/publishers", http.HandlerFunc(data.PubDeetsHandler))
mux.Handle("/gridwatch/", http.HandlerFunc(indexHandler))
mux.Handle("/", http.HandlerFunc(handleRoot))
*/
