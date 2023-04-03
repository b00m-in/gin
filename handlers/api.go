package handlers

import (
        "encoding/json"
        "fmt"
        "io"
        "net/http"
        "strings"
        "strconv"
        "time"
        "github.com/gin-gonic/gin"
        "github.com/golang/glog"
        "b00m.in/data"
        "b00m.in/comms"
        "b00m.in/gin/glue"
)

func HandleApiGET(c *gin.Context) {
        //c.HTML(200, "user_content.html", gin.H{"title": "subs content"})
        r := c.Request
        w := c.Writer
        toks := strings.Split(r.URL.Path, "/")
        glog.Infof("handleApi %s %v \n", r.Method, toks)
        if len(toks) < 3 {
                glog.Infof("No get/post at %s \n", r.URL.Path)
                w.WriteHeader(http.StatusNoContent)
                return
        }
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
        case "pubs":
                //pubs, err := data.GetPubs(100)
                pubs, err := data.GetPubDummies(50)
                if err != nil {
                        str := fmt.Sprintf("Couldn't query results: %s", err)
                        http.Error(w, str, 500)
                        return
                }
                glog.Infof("Dummies: %v \n", pubs[0])
                err = json.NewEncoder(w).Encode(pubs)
                if err != nil {
                        str := fmt.Sprintf("Couldn't encode results: %s", err)
                        http.Error(w, str, 500)
                        return
                }
                return
        case "subs":
                return
        case "packets":
                return
        case "summary": // /api/summary/hourly/hash/2021-Jan-01/2021-Jan-31
                var err error
                freqp := c.Params.ByName("freq")
                hashp := c.Params.ByName("hash")
                startp := c.Params.ByName("start")
                endp := c.Params.ByName("end")
                hash := int64(13582)
                freq := "hourly"
                if freqp != "" {
                        freq = freqp
                }
                from := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
                to := time.Date(2021, time.January, 31, 0, 0, 0, 0, time.UTC)
                freq = toks[3]
                hash, err = strconv.ParseInt(hashp, 10, 64)
                if err != nil {
                        glog.Errorf("strconv parseint %v \n", err)
                }
                from, err = time.Parse(shortform, startp)
                if err != nil {
                        glog.Errorf("time.Parse %v \n", err)
                        from, err = time.Parse(shortform1, startp)
                        if err != nil {
                                glog.Errorf("time.Parse %v \n", err)
                        }
                }
                to, err = time.Parse(shortform, endp)
                if err != nil {
                        glog.Errorf("time.Parse %v \n", err)
                        to, err = time.Parse(shortform1, endp)
                        if err != nil {
                                glog.Errorf("time.Parse %v \n", err)
                        }
                }
                ss, err := data.GetSummariesByHash(from, to, freq, hash)
                if err != nil {
                        glog.Errorf("GetSummaries %v \n", err)
                }
                glog.Infof("handleApi get api/summary: %v %v %v %v \n", from, to, freq, len(ss))
                w.Header().Set("Access-Control-Allow-Origin", "*")
                err = json.NewEncoder(w).Encode(ss)
                if err != nil {
                        glog.Errorf("json.Encode %v \n", err)
                        str := fmt.Sprintf("Couldn't encode results: %s", err)
                        http.Error(w, str, 500)
                        return
                }
        }
}

func HandleApiRegisterPOST(c *gin.Context) {
        //r := c.Request
        w := c.Writer
        //if sub == "" {
                glog.Infof("handleApi post api/register w/o sub \n")
                fn := c.PostForm("fullname")
                e := c.PostForm("email")
                pw := c.PostForm("pswd")
                ph := c.PostForm("phone")
                cc := c.PostForm("confirm")
                s, err := data.GetSubByEmail(e)
                if s != nil || err == nil {
                        glog.Errorf("handleApi post register s: %s err: %v \n", e, err)
                        w.WriteHeader(http.StatusConflict)
                        io.WriteString(w, "Email already registered")
                        return
                }
                if pw != cc {
                        glog.Errorf("handleApi post register %s == %s \n", cc, pw)
                        w.WriteHeader(http.StatusConflict)
                        io.WriteString(w, "Password mismatch")
                        return
                }
                // success
                //glog.Infof("handleApi post login success %s \n", e)
                //http.SetCookie(w, &http.Cookie{Name: "sub", Value: e, Domain:"b00m.in", Path: "/", MaxAge: 300, HttpOnly: true, Expires: time.Now().Add(time.Second * 120)})
                glog.Infof("handleapi post register %s \n", e)
                sha1str := data.Sha1Str(e) // 16 characters of sha1 hash
                i, err := data.PutSub(&data.Sub{Email:e, Name:fn, Phone:ph, Pswd:pw, Verification: sha1str})
                if err != nil {
                        glog.Infof("handleapi post register %s %v \n", e, err)
                        w.WriteHeader(http.StatusServiceUnavailable)
                        io.WriteString(w, "Sorry couldn't register")
                        return
                }
                glue.Newregs <- comms.Entity{e, fn}
                w.WriteHeader(http.StatusOK)
                io.WriteString(w, "Registered " + strconv.Itoa(int(i)))
                return
        /*} else {
                glog.Infof("handlesubs post sub/lin w/ %s \n", sub)
                render := Render{Message: "Already " + sub, Categories: dflt_ctgrs}
                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                return
        }*/
}

func HandleApiLoginPOST(c *gin.Context) {
        w := c.Writer
        //if sub == "" {
                glog.Infof("handleApi post api/login w/o sub \n")
                e := c.PostForm("email")
                pw := c.PostForm("pswd")
                s, err := data.GetSubByEmail(e)
                if err != nil {
                        glog.Errorf("handleApi post login %v \n", err)
                        w.WriteHeader(http.StatusUnauthorized)
                        return
                }
                if !data.CheckPswd(e, pw) {
                        glog.Errorf("handleApi post putsubs %s == %s \n", s.Pswd, pw)
                        w.WriteHeader(http.StatusUnauthorized)
                        return
                }
                // success
                //glog.Infof("handleApi post login success %s \n", e)
                //http.SetCookie(w, &http.Cookie{Name: "sub", Value: e, Domain:"b00m.in", Path: "/", MaxAge: 300, HttpOnly: true, Expires: time.Now().Add(time.Second * 120)})
                glog.Infof("handleapi post login %s %s \n", s.Name, e)
                w.WriteHeader(http.StatusOK)
                io.WriteString(w, "Logged in")
                return
        /*} else {
                glog.Infof("handlesubs post sub/lin w/ %s \n", sub)
                render := Render{Message: "Already " + sub, Categories: dflt_ctgrs}
                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                return
        }*/
}
