package glue

import (
        "fmt"
        "time"
        "github.com/golang/glog"
        "b00m.in/data"
        "b00m.in/comms"
)

var(
        newfaults chan Fault
        from comms.Entity
        Newregs chan comms.Entity
)

type Fault struct {
        Entity comms.Entity
        Hash int64
        Nickname string
}

func Setup(email, name string) {
        from = comms.Entity{email, name}
}

// loadPubDdeets also listens on channel to recieve data of new registrants to send email. It also listens on a channel to send email notifications of faults to subscribers
func LoadPubdeets(debug bool) {
        //from := comms.Entity{"rcs@b00m.in", "RCS"}
        to := make([]comms.Entity, 0)
        for {
                select {
                case newreg:= <-Newregs:
                        if !debug {
                                to = append(to, newreg)
                                sha1str := data.Sha1Str(newreg.Email)
                                err := comms.SendMail(from, to, "Welcome to B00M","Welcome. Please verify your email https://pv.b00m.in/subs/verify/" + sha1str)
                                if err != nil {
                                        glog.Infof("comms.SendEmail %v \n", err)
                                } else {
                                        glog.Infof("comms.SendEmail newreg %v \n", newreg.Email)
                                }
                                to = to[:0]
                        } else {
                                glog.Infof("simulate send reg email from %s to %s \n", from, to)
                        }
                case newfault:= <-newfaults:
                        if !debug {
                                to = append(to, newfault.Entity)
                                err := comms.SendMail(from, to, "B00M Fault!", fmt.Sprintf("Please check on %s - %d \n", newfault.Nickname, newfault.Hash))
                                if err != nil {
                                        glog.Infof("comms.SendMail %v \n", err)
                                } else {
                                        glog.Infof("comms.SendMail fault %v email %v \n", newfault.Nickname, newfault.Entity.Email)
                                        // email sent so update pubconfig.lastnotified time and pub.protected
                                        if err := data.UpdatePubStatus(newfault.Hash); err != nil {
                                                glog.Errorf("updatelastnotified %v \n", err)
                                        }
                                }
                                to = to[:0]
                        } else {
                                glog.Infof("simulate send fault email from %s to %s \n", from, to)
                        }
                case <-time.After(time.Duration(600 * time.Second)):
                        //if err := data.LoadDummyPubdeets(); err != nil {
                        if err := data.LoadDummiesForAll(); err != nil {
                                glog.Infof("data.LoadDummiesForAll %v \n", err)
                        }
                        faults, err := data.GetPubFaults(true)
                        if err!= nil {
                                glog.Infof("data.LoadFaults %v \n", err)
                        }
                        for _, fault := range faults {
                                newfaults <- Fault{comms.Entity{fault.Email, fault.Name}, fault.Hash, fault.Nickname}
                        }
                }
        }
}
