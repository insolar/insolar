module github.com/insolar/insolar

go 1.12

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	github.com/ThreeDotsLabs/watermill v1.0.2
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cheggaaa/pb/v3 v3.0.1
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/dgraph-io/badger v1.6.0-rc1.0.20191024172150-efb9d9d15d7f
	github.com/dustin/go-humanize v1.0.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/getkin/kin-openapi v0.2.1-0.20191211203508-0d9caf80ada6
	github.com/gogo/protobuf v1.2.1
	github.com/gojuno/minimock/v3 v3.0.5
	github.com/golang/protobuf v1.3.2
	github.com/google/gofuzz v1.0.0
	github.com/google/gops v0.3.6
	github.com/grpc-ecosystem/grpc-gateway v1.9.6
	github.com/hashicorp/golang-lru v0.5.3
	github.com/insolar/component-manager v0.2.1-0.20191028200619-751a91771d2f
	github.com/insolar/go-actors v0.0.0-20190805151516-2fcc7bfc8ff9
	github.com/insolar/insconfig v0.0.0-20200227134411-011eca6dc866
	github.com/insolar/rpc v1.2.2-0.20190812143745-c27e1d218f1f
	github.com/insolar/x-crypto v0.0.0-20191031140942-75fab8a325f6
	github.com/jackc/pgx/v4 v4.2.1
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/olekukonko/tablewriter v0.0.1
	github.com/onrik/gomerkle v1.0.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/ory/dockertest/v3 v3.5.2
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/common v0.6.0 // indirect
	github.com/prometheus/procfs v0.0.4 // indirect
	github.com/rs/zerolog v1.15.0
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/tylerb/gls v0.0.0-20150407001822-e606233f194d
	github.com/tylerb/is v2.1.4+incompatible // indirect
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/uber/jaeger-client-go v2.19.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/ugorji/go v1.1.4
	go.opencensus.io v0.22.0
	go.uber.org/goleak v1.0.0
	golang.org/x/crypto v0.0.0-20190911031432-227b76d455e7
	golang.org/x/net v0.0.0-20191003171128-d98b1b443823
	golang.org/x/tools v0.0.0-20191108193012-7d206e10da11
	gonum.org/v1/gonum v0.0.0-20191018104224-74cb7b153f2c
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
	google.golang.org/grpc v1.21.0
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/insolar/insolar => ./
