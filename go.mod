module github.com/insolar/insolar

go 1.15

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/ThreeDotsLabs/watermill v1.0.2
	github.com/cheggaaa/pb/v3 v3.0.1
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/dgraph-io/badger v1.6.2
	github.com/dgraph-io/ristretto v0.0.3 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gbrlsnchs/jwt/v3 v3.0.0-rc.2
	github.com/getkin/kin-openapi v0.2.1-0.20191211203508-0d9caf80ada6
	github.com/gogo/protobuf v1.3.1
	github.com/gojuno/minimock/v3 v3.0.8
	github.com/golang/protobuf v1.4.2
	github.com/google/gofuzz v1.0.0
	github.com/google/gops v0.3.6
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.9.6
	github.com/hashicorp/golang-lru v0.5.3
	github.com/hexdigest/gowrap v1.1.10 // indirect
	github.com/insolar/component-manager v0.2.1-0.20191028200619-751a91771d2f
	github.com/insolar/gls v0.0.0-20200427111849-9a08a622625d
	github.com/insolar/go-actors v0.0.0-20190805151516-2fcc7bfc8ff9
	github.com/insolar/insconfig v0.0.0-20200925122209-0598492b58e8
	github.com/insolar/ratelimiter v1.0.0
	github.com/insolar/rpc v1.2.2-0.20200331123021-db57b8833f82
	github.com/insolar/x-crypto v0.0.0-20191031140942-75fab8a325f6
	github.com/jackc/pgx/v4 v4.2.1
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/olekukonko/tablewriter v0.0.1
	github.com/onrik/gomerkle v1.0.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/ory/dockertest/v3 v3.5.2
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/rs/zerolog v1.15.0
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.4.0
	github.com/tylerb/is v2.1.4+incompatible // indirect
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/uber/jaeger-client-go v2.19.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/ugorji/go/codec v1.1.8
	go.opencensus.io v0.22.0
	go.uber.org/goleak v1.0.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20200925080053-05aa5d4ee321
	golang.org/x/sys v0.0.0-20200923182605-d9f96fdee20d // indirect
	golang.org/x/tools v0.0.0-20200925150135-34b80a0a46ad
	gonum.org/v1/gonum v0.0.0-20191018104224-74cb7b153f2c
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.27.1
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/insolar/insolar => ./
