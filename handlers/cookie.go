package handlers

import (
        "github.com/gin-gonic/gin"
        "github.com/gin-contrib/sessions"
        "github.com/golang/glog"
        "b00m.in/data"
        "encoding/gob"
)
var(
        cc = 1
        guest = &data.Sub{Name: "new", Email: "gst"}
)
func init() {

        gob.Register(&data.Sub{})
}

func ReadCookie() gin.HandlerFunc {
        return func(ctx *gin.Context) {
                var you *data.Sub
                var err error
                s := sessions.DefaultMany(ctx, "sub")
                cs := s.Get("count")
                if cs == nil {
                        glog.Infof("cs nil: %d\n", cs)
                        s.Set("count", cc)
                } else {
                        glog.Infof("cs: %d\n", cs)
                        cs = cs.(int)+1
                        s.Set("count", cs)
                        //cc = cs.(int)
                }
                us := s.Get("useremail")
                if us == nil || us == "gst" {
                        you = &data.Sub{}
                        s.Set("useremail", "gst")
                        glog.Infof("new user %v\n", us)
                } else {
                        uss := us.(string)
                        you, err = data.GetSubByEmail(uss) // expensive
                        if err != nil {
                                glog.Errorf("%v\n", err)
                                you = &data.Sub{Name: "new", Email:"gst"}
                                s.Set("useremail", "gst")
                        } else {
                                s.Set("useremail", you.Email)
                                glog.Infof("user session %v %v\n", you.Name, you.Email)
                        }
                }
                s.Set("user", you)
                if err := s.Save(); err != nil {
                        glog.Errorf("cookie save %v\n", err)
                }
        }
}

// WriteCookie strips out all user related data other then email from the session
func WriteCookie() gin.HandlerFunc {
        return func(ctx *gin.Context) {
                var you *data.Sub
                s := sessions.DefaultMany(ctx, "sub")
                us := s.Get("user")
                if us == nil {
                        s.Set("useremail", "gst")
                        glog.Infof("user nil cookie set: %v\n", "gst")
                } else {
                        you = us.(*data.Sub)
                        //s.Clear()
                        s.Set("useremail", you.Email)
                        glog.Infof("user %v cookie set: %v\n", us, you.Email)
                }
                if err := s.Save(); err != nil {
                        glog.Errorf("cookie save %v\n", err)
                }
        }
}

func ClearCookie() gin.HandlerFunc {
        return func(ctx *gin.Context) {
                s := sessions.DefaultMany(ctx, "sub")
                s.Clear()
                s.Save()
        }
}
