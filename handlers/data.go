package handlers

import (
        "encoding/json"
        "fmt"
        "net/http"
        "b00m.in/tmpl"
	"b00m.in/data"
        "github.com/gin-gonic/gin"
        "github.com/gin-gonic/gin/binding"
        "github.com/golang/glog"
        "github.com/gin-contrib/sessions"
)

var once bool

func init() {
        once = false
        /*if err := data.LoadDummyPubdeets(); err != nil {
                glog.Infof("%v \n", err)
        }*/

}

// pubdeetsHandler reads r for a "viewport" query parameter
// and writes a GeoJSON response of the features contained in
// that viewport into w.
func PubDeetsHandler(c *gin.Context) {
        if !once {
                if err := data.LoadDummiesForAll(); err != nil {
                        glog.Errorf("%v \n", err)
                }
                once = true
        }
        //r := c.Request
        w := c.Writer
        sub := ""
        var you *data.Sub
        session := sessions.DefaultMany(c, "sub")
        if v := session.Get("user"); v != nil {
                you = v.(*data.Sub)
                if you != nil && you.Name != "new" {
                        sub = you.Email
                }
        }
        glog.Infof("pubdeetshandler: %s\n", sub)
	w.Header().Set("Content-type", "application/json")
	//vp := r.FormValue("viewport")
        var pfvp tmpl.PostFormVP
        if err := c.ShouldBindWith(&pfvp, binding.Query); err != nil {
                glog.Infof("query bind: %v \n", err)
                pfvp = tmpl.PostFormVP{"", "10", ""}
        }
        if you == nil {
                glog.Infof("you nil: zoom \n", pfvp.Zoom)
                str := fmt.Sprintf("You nil: %s", "boo")
                glog.Error("PubDeetshandler couldn't encode results %v \n", "boo")
                http.Error(w, str, 500)
                return
        } else {
                yid := fmt.Sprintf("%d", you.Id)
                fc, err := data.FilterDummiesForSub(pfvp.Viewport, pfvp.Zoom, pfvp.SubId)
                if err != nil {
                        glog.Errorf("filterdummiesforsub %s %v \n", yid, err)
                }
                err = json.NewEncoder(w).Encode(fc)
                if err != nil {
                        str := fmt.Sprintf("Couldn't encode results: %s", err)
                        glog.Error("PubDeetshandler couldn't encode results %v \n", err)
                        http.Error(w, str, 500)
                        return
                }
        }
}
