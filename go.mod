module b00m-in/gin

require (
	b00m.in/crypto/util v0.0.0-00010101000000-000000000000
	b00m.in/gin/routes v0.0.0-00010101000000-000000000000
	b00m.in/gin/tmpls v0.0.0-00010101000000-000000000000
	b00m.in/subs v0.0.0-00010101000000-000000000000
	b00m.in/xds/resource v0.0.0-00010101000000-000000000000
	b00m.in/xds/server v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.8.1
	github.com/golang/glog v1.0.0
)

require (
	b00m.in/comms v0.0.0-00010101000000-000000000000 // indirect
	b00m.in/data v0.0.0-00010101000000-000000000000 // indirect
	b00m.in/gin/handlers v0.0.0-00010101000000-000000000000 // indirect
	b00m.in/tmpl v0.0.0-00010101000000-000000000000 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cncf/xds/go v0.0.0-20220314180256-7f1daf1720fc // indirect
	github.com/dhconnelly/rtreego v1.1.0 // indirect
	github.com/envoyproxy/go-control-plane v0.10.3-0.20221219165740-8b998257ff09 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.9.1 // indirect
	github.com/gin-contrib/multitemplate v0.0.0-20220829131020-8c2a8441bc2b // indirect
	github.com/gin-contrib/sessions v0.0.5 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/paulmach/go.geojson v1.4.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/plutov/paypal/v4 v4.0.0 // indirect
	github.com/smira/go-point-clustering v1.0.1 // indirect
	github.com/stripe/stripe-go/v72 v72.122.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/willf/bitset v1.1.10 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20220822174746-9e6da59bd2fc // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace b00m.in/gin/handlers => ./handlers

replace b00m.in/gin/routes => ./routes

replace b00m.in/gin/tmpls => ./tmpls

replace b00m.in/data => ../data

replace b00m.in/crypto/util => ../crypto/util

replace b00m.in/xds/resource => ../xds/resource

replace b00m.in/xds/server => ../xds/server

replace b00m.in/comms => ../comms

replace b00m.in/subs => ../subs

replace b00m.in/tmpl => ../tmpl

go 1.19
