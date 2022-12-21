package handlers

import (
        //"encoding/json"
        //"fmt"
        //"io"
        "net/http"
        "strings"
        "strconv"
        "time"
        "github.com/gin-gonic/gin"
        "github.com/golang/glog"
        "github.com/gin-contrib/sessions"
        "b00m.in/data"
        "b00m.in/comms"
        "b00m.in/tmpl"
)

var (
        newregs chan comms.Entity
        //dflt_ctgrs = subs.Category{}
        dflt1_ctgrs = []string{"News", "Docs", "Gridwatch", "Leaderboard", "Community", "Github"}
        shortform = "2006-Jan-02"
        shortform1 = "2006-01-02"
)

func HandleSubsGET(c *gin.Context) {
        r := c.Request
        w := c.Writer
        toks := strings.Split(r.URL.Path, "/")
        glog.Infof("handleSubs %s %v \n", r.Method, toks)
        if len(toks) < 3 {
                glog.Infof("No get/post at %s \n", r.URL.Path)
                w.WriteHeader(http.StatusNoContent)
                return
        }
        sub := ""
        var you *data.Sub
        session := sessions.DefaultMany(c, "sub")
        if v := session.Get("user"); v != nil {
                you = v.(*data.Sub)
                if you != nil && you.Name != "new" {
                        sub = you.Email
                }
        }
        switch r.Method{
        case "GET":
                switch toks[2] {
                case "confo":
                        devicename := toks[3] // "Movprov"
                        ssid := toks[4] // "M0V"
                        confo, err := data.GetLastConfo(devicename, ssid)
                        if err != nil {
                                w.WriteHeader(http.StatusOK)
                                w.Write([]byte("Not found"))
                                return
                        }
                        w.WriteHeader(http.StatusOK)
                        w.Write([]byte(strconv.FormatInt(confo.Hash, 10)))
                        return
                case "faults":
                        if sub == "" {
                                glog.Infof("handlenew get faults no cookie \n")
                                render := tmpl.Render {Message: "Sign in", Categories: dflt1_ctgrs, User: you.Name}
                                c.HTML(200, "subs_login.html", render)
                                return
                        } else {
                                glog.Infof("handlenew get faults %s : %d \n", sub, you.Id)
                                pubs, err := data.GetPubFaultsForSub(you.Id)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := tmpl.Render{Message: "No Faults", Categories: dflt1_ctgrs, User: you.Name}
                                        c.HTML(200, "bodyfaults.html", rendere)
                                        return
                                }
                                render := tmpl.Render {Message: "Faults", Pubs: pubs, Categories: dflt1_ctgrs, User: you.Name}
                                c.HTML(200, "bodyfaults.html", render)
                                return
                        }
                case "pubs":
                        if sub == "" {
                                c.HTML(200, "subs_login.html", gin.H{})
                                return
                        } else {
                                if len(toks) == 3 { // /subs/pubs
                                        glog.Infof("handlenewsubs get pubs %s : %d \n", sub, you.Id)
                                        pubs, err := data.GetPubsForSub(you.Id)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := tmpl.Render{Message: "No pubs for you", Categories: dflt1_ctgrs, User: you.Name}
                                                c.HTML(200, "bodypubs.html", rendere)
                                                return
                                        }
                                        render := tmpl.Render {Message: "Welcome " + you.Name + " - Devices", Pubs: pubs, Categories: dflt1_ctgrs, User: you.Name}
                                        c.HTML(200, "bodypubs.html", render)
                                        return
                                } else { // /subs/pubs/:id
                                        var id int64
                                        idp := c.Params.ByName("hash")
                                        id, err := strconv.ParseInt(idp, 10, 64)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := tmpl.Render{Message: "No such pub", Categories: dflt1_ctgrs, User: you.Name}
                                                c.HTML(200, "bodypubs.html", rendere)
                                                return
                                        }
                                        pub, err := data.GetPubById(id)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := tmpl.Render{Message: "No pubs for you", Categories: dflt1_ctgrs, User: you.Name}
                                                c.HTML(200, "bodypubs.html", rendere)
                                                return
                                        }
                                        var render tmpl.RenderOne
                                        devicename,err := data.GetPubDeviceName(pub.Hash)
                                        if err != nil {
                                                glog.Errorf("Https handle subs/pubs/%d : %v \n", id, err)
                                                pc := &data.PubConfig{Kwp: 0, Kwpmake: "unknown", Kwr: 0, Kwrmake: "unknown"}
                                                render = tmpl.RenderOne{Message: "Unknown", Pub: pub, PubConfig: pc, Categories: dflt1_ctgrs, User: you.Name}
                                        } else {
                                                pc, err := data.GetPubConfigByHash(pub.Hash)
                                                if err != nil {
                                                        pc = &data.PubConfig{Kwp: 0, Kwpmake: "unknown", Kwr: 0, Kwrmake: "unknown"}
                                                }
                                                render = tmpl.RenderOne{Message: devicename, Pub: pub, PubConfig: pc, Categories: dflt1_ctgrs, User: you.Name}
                                        }
                                        c.HTML(200, "bodyconfig.html", render)
                                        return
                                }
                        }
                case "subs":
                        if sub == "" {
                                c.HTML(200, "subs_login.html", gin.H{})
                                return
                        } else {
                                render := tmpl.GetSummaryRender(you)
                                c.HTML(200, "bodysummary.html", gin.H{"Total": render.Total, "Online": render.Online, "Protected": render.Protected})
                                return
                        }
                case "login":
                        if sub == "" {
                                c.HTML(200, "subs_login.html", gin.H{})
                                return
                        } else {
                                render := tmpl.GetSummaryRender(you)
                                c.HTML(200, "bodysummary.html", gin.H{"Total": render.Total, "Online": render.Online, "Protected": render.Protected})
                                return
                        }
                case "logout":
                        session.Clear()
                        session.Save()
                        c.HTML(200, "subs_login.html", gin.H{})
                        return
                case "packets":
                        if sub == "" {
                                c.HTML(200, "subs_login.html", gin.H{})
                                return
                        } else { // /subs/packets/:id
                                idp := c.Params.ByName("hash")
                                id, err := strconv.ParseInt(idp, 10, 64)
                                if err != nil {
                                        glog.Infof("strconv: %v \n", err)
                                        rendere := tmpl.Render{Message: "No such pub", Categories: dflt1_ctgrs, User: you.Name}
                                        c.HTML(200, "bodypubs.html", rendere)
                                        return
                                }
                                np := 100
                                pks, err := data.GetLastPackets(id, np)
                                if err != nil || pks == nil || len(pks) == 0 {
                                        glog.Infof("Https %v \n", err)
                                        rendere := tmpl.Render {Message: "Packets error", Categories: dflt1_ctgrs, User: you.Name}
                                        c.HTML(200, "bodypubs.html", rendere)
                                        return
                                }
                                end := pks[0].Timestamp.Format("2006-Jan-02")
                                start := pks[0].Timestamp.Add(time.Hour * -48).Format("2006-Jan-02")
                                render := tmpl.Rendern {Message: "Most Recent Readings", Id: id, Packets: pks, Start: start, End: end, Freq: "hourly", Categories: dflt1_ctgrs, User: you.Name}
                                c.HTML(200, "bodypackets.html", render)
                                return
                        }
                case "register":
                        c.HTML(200, "subs_new.html", gin.H{})
                        return
                }
        case "POST":
                switch toks[2] {
                case "new":
                        if sub == "" {
                                var ds data.Sub
                                if err := c.ShouldBind(&ds); err != nil {
                                        glog.Errorf("bind %v \n",err)
                                        c.HTML(200, "subs_new.html", gin.H{})
                                        return

                                }
                                if s, err := data.GetSubByEmail(ds.Email); s != nil {
                                        glog.Errorf("handlesubs post putsubs %v \n", err)
                                        c.HTML(200, "subs_new.html", gin.H{})
                                        return
                                }
                                sha1str := data.Sha1Str(ds.Email) // 16 characters of sha1 hash
                                ds.Verification = sha1str
                                id, err := data.PutSub(&ds)
                                if err != nil {
                                        glog.Errorf("handlesubs post putsubs %v \n", err)
                                        c.HTML(200, "subs_new.html", gin.H{})
                                        return
                                }
                                // success
                                //newregs <- comms.Entity{e, n} // put in channel to send email
                                glog.Infof("handleSubs set cookie %s \n", ds.Email)
                                //http.SetCookie(w, &http.Cookie{Name: "sub", Value: e, Domain:"b00m.in", Path: "/", MaxAge: 600, HttpOnly: true, Expires: time.Now().Add(time.Second * 600)})
                                glog.Infof("handlesubs post putsubs %d \n", id)
                                c.HTML(200, "subs_login.html", gin.H{})
                                return
                        } else {
                                glog.Infof("handlesubs post sub/new with %s \n", sub)
                                c.HTML(200, "subs_login.html", gin.H{"user": you.Name})
                                return
                        }
                case "login":
                        if sub == "" {
                                glog.Infof("handlesubs post sub/login w/o sub \n")
                                var ds data.Sub
                                if err := c.ShouldBind(&ds); err != nil {
                                        glog.Infof("agreed: %v \n", err)
                                        c.HTML(200, "subs_login.html", gin.H{})
                                        return
                                }
                                s, err := data.GetSubByEmail(ds.Email)
                                if err != nil {
                                        glog.Errorf("handlesubs post login %v \n", err)
                                        c.HTML(200, "subs_login.html", gin.H{})
                                        return
                                }
                                if !data.CheckPswd(ds.Email, ds.Pswd) {
                                        glog.Errorf("handlesubs post putsubs %s == %s \n", ds.Pswd, s.Pswd)
                                        c.HTML(200, "subs_login.html", gin.H{})
                                        return
                                }
                                // success
                                glog.Infof("handleSubs post login %s \n", ds.Email)
                                you = s
                                you.Pswd = "****"
                                session.Set("user", you)
                                session.Set("useremail", you.Email)
                                if err := session.Save(); err != nil {
                                        glog.Errorf("cookie save %v\n", err)
                                }
                                glog.Infof("handlesubs post login %s %s \n", s.Name, ds.Email)
                                render := tmpl.GetSummaryRender(you)
                                c.HTML(200, "bodysummary.html", gin.H{"Total": render.Total, "Online": render.Online, "Protected": render.Protected})
                                return
                        } else {
                                glog.Infof("handlesubs post sub/lin w/ %s \n", sub)
                                c.HTML(200, "subs_login.html", gin.H{})
                                return
                        }
                case "packets":
                        if sub == "" {
                                glog.Infof("handlesubs post sub/packets w/o sub \n")
                                c.HTML(200, "subs_login.html", gin.H{})
                                return
                        } else {
                                hashp := c.Params.ByName("hash")
                                id, err := strconv.ParseInt(hashp, 10, 64)
                                if err != nil {
                                        glog.Infof("strconv: %v \n", err)
                                }
                                fr := "hourly"
                                var pps data.PostFormData
                                if err := c.ShouldBind(&pps); err != nil {
                                        glog.Infof("form bind: %v \n", err)
                                }
                                start, err := time.Parse(shortform, pps.Start)
                                if err != nil {
                                        glog.Infof("time.Parse: %v \n", err)
                                        start, err = time.Parse(shortform1, pps.Start)
                                        if err != nil {
                                                glog.Infof("time.Parse: %v \n", err)
                                                start = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
                                        }
                                }
                                end,err := time.Parse(shortform, pps.End)
                                if err != nil {
                                        glog.Infof("time.Parse: %v \n", err)
                                        end, err = time.Parse(shortform1, pps.End)
                                        if err != nil {
                                                glog.Infof("time.Parse: %v \n", err)
                                                end = time.Date(2021, time.January, 31, 0, 0, 0, 0, time.UTC)
                                        }
                                }
                                if end.Sub(start) > (time.Hour * 24) {
                                        fr = "daily"
                                }
                                pks := make([]*data.Packet, 0)
                                render := tmpl.Rendern {Message: "Summaries", Id: id, Packets: pks, Start: pps.Start, End: pps.End, Freq:fr, Categories: dflt1_ctgrs, User: you.Name}
                                c.HTML(200, "bodypackets.html", render)
                                return
                        }
                default:
                        glog.Infof("No posting at %s \n", r.URL.Path)
                        w.WriteHeader(http.StatusNotFound)
                        return
                }
        }
}

func HandleSubsPOST(c *gin.Context) {
}

func HandleSubsGet(w http.ResponseWriter, r *http.Request){
        sub := ""
        var you *data.Sub
        ctx := r.Context()
        if v := ctx.Value("user"); v != nil {
                //fmt.Println("found value:", v)
                you = v.(*data.Sub)
                if you != nil && you.Name != "new" {
                        sub = you.Email
                }
        }
        glog.Infof("User: %s\n", sub)
        http.SetCookie(w, &http.Cookie{Name: "sub", Value: "gst", Domain:"", SameSite: http.SameSiteLaxMode, Path: "/", MaxAge: 600, HttpOnly: true, Expires: time.Now().Add(time.Second * 600)})
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("HandleSubsGet"))
}
