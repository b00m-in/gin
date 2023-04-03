package main

import (
        "context"
        "encoding/pem"
        "flag"
        "fmt"
        "io/ioutil"
        "github.com/golang/glog"
        "github.com/gin-gonic/gin"
        "b00m.in/gin/routes"
        "b00m.in/gin/tmpls"
        "b00m.in/gin/glue"
        "b00m.in/subs"
        xdsr "b00m.in/xds/resource"
        xdss "b00m.in/xds/server"
        "b00m.in/crypto/util"
        "net/http"
        "os"
        "os/signal"
        "syscall"
        "strconv"
        "strings"
        "time"
)

var (
        cfg = flag.String("c", "config/b00m.config", "path to config file") // this isn't even used anywhere
        srv http.Server
        httpPort int
)

func main() {
	flag.Parse() // does this need to be called. there are no flags at this point - maybe only stderrthreshold?
        config, err := subs.ConfigureConfig(os.Args[1:])
        if err != nil {
                glog.Errorf("configureconfig %v \n", err)
        }
        pkb, err := ioutil.ReadFile(config.Keypair.Private)
        if err != nil {
                glog.Errorf("readfile key %v \n", err)
        }
        priv,_ := pem.Decode(pkb)
        if priv == nil || !strings.Contains(priv.Type, "PRIVATE") {
                glog.Infof("%s \n", "Nil rsa key")
        }
        /*pkb, err = ioutil.ReadFile(config.Keypair.Public)
        if err != nil {
                glog.Errorf("readfile key %v \n", err)
        }
        pub,_ := pem.Decode(pkb)
        if pub == nil || !strings.Contains(pub.Type, "PUBLIC") {
                glog.Infof("%s \n", "Nil rsa key")
        }*/
        httpPort, err = strconv.Atoi(config.HTTPPort)
        if err != nil {
                glog.Errorf("strconv httpport exiting %v \n", err)
                return
        }
        ctx, stop := signal.NotifyContext(context.Background(),syscall.SIGINT, syscall.SIGTERM)
        defer stop()
        r := gin.Default()
        routes.SetupRouter(r)
        errch := make(chan error, 1)
        if config.Xds != nil {
                v := "13.1"
                rr := xdsr.NewRR(config.Xds.RouteName, config.Xds.ClusterName, r.Routes())
                c := rr.GenerateRoutesSnapshot(v)
                s := xdss.NewServer(util.SDS{NodeId: config.Xds.NodeId, SdsPort: config.Xds.SdsPort, Debug: config.Xds.Debug})
                srv, err := s.NewSDSWithCache()
                if err != nil {
                        glog.Errorf("%v\n", err)
                }
                s.SetCacheSnapshot(c)
                go xdss.RunServerGo(errch, srv, uint(config.Xds.SdsPort))
        }
        tmpls.SetupTemplates(r, &config.Categories)
        if config.From != nil {
                glue.Setup(config.From.Fromemail, config.From.Fromname)
                go glue.LoadPubdeets(config.Debug)
        }
        //r.Run(fmt.Sprintf(":%d", p)) //(":8080")
        go startHttp(r)
        <-ctx.Done()
        stop()
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		glog.Fatalf("Server forced to shutdown: %v ", err)
	}

	glog.Infof("%s\n", "Server exiting")
}

func startHttp(r *gin.Engine) {
        //mux := http.NewServeMux()
        //mux.Handle("/debug/vars", expvar.Handler())
        /*mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                oks.Add(1)
                fmt.Printf("Visitor %v \n", *oks)
                fmt.Fprintf(w, "Nothing to see here \n")
        }))*/
        /*if redirectHttp {
                mux.Handle("/", http.HandlerFunc(RedirectHttp))
        } else {
                //mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
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
        }*/
        tr := http.TimeoutHandler(r, time.Second*6, "Request timed out")
        srv = http.Server{
                ReadTimeout: 10 * time.Second,
                WriteTimeout: time.Duration(10) * time.Second,
                Addr: fmt.Sprintf(":%d", httpPort),
                Handler: tr,
        }
        err := srv.ListenAndServe()
        if err != nil {
                fmt.Printf("Oops: %v \n", err)
        }
}
