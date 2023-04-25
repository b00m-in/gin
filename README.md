# b00m.in/gin


## Contents

Contains a simple [gin](https://github.com/gin-gonic/gin) server, router and templates for display of data.

`b00m.in/xds` is used to also launch a discovery service to publish RoutesInfo to a proxy like envoy.


## Missing from repo:

`assets/code` should contain highcharts library related files

## Login/out

`handlers/cookie.go` has a `map[string]*data.Sub` named `current` which stores the current users logged in. This prevents expensive identifying database calls for requests from users already logged in.

`HandleSubsLogout` should delete the user's entry in `current` and  should also expire the session `session.Options(sessions.Options{MaxAge:-1})` so the session cookie gets deleted by the browser or else it persists and pollutes further requests

Newly registered users are sent an email to verify their registration. This is done using the `Newregs` channel in `glue`. This channel only gets set up if the config file provided at start-up contains relevant email api info. If not the channel never gets setup and sends to it may panic or block forever. To avoid this a `send or drop` select pattern is used:

```
select { // don't block - send or drop
        case glue.Newregs <- comms.Entity{ds.Email, ds.Name}: // put in channel to send email
        default:
}
```
## Service Discovery

[XDS](https://github.com/b00m-in/xds) encapsulates a [grpc xds server](https://github.com/envoyproxy/go-control-plane) which can be used to publish secrets and routes to [envoyproxy](https://envoyproxy.io).

This can be used to enable gin to publish it's routes as an RDS (route discovery service) which envoyproxy can then be configured to dynamically load whenever new routes are available or old ones are disabled. 
