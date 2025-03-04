module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver

go 1.18

require (
	github.com/census-instrumentation/opencensus-proto v0.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.60.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.60.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.60.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.60.0
	github.com/rs/cors v1.8.2
	github.com/soheilhy/cmux v0.1.5
	github.com/stretchr/testify v1.8.0
	go.opentelemetry.io/collector v0.60.1-0.20220916163348-84621e483dfb
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.35.0
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
)

require (
	go.opentelemetry.io/collector/pdata v0.60.1-0.20220916163348-84621e483dfb
	go.opentelemetry.io/otel/sdk v1.10.0
)

require (
	cloud.google.com/go/compute v1.10.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/knadh/koanf v1.4.3 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.1.17 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/collector/semconv v0.60.1-0.20220916163348-84621e483dfb // indirect
	go.opentelemetry.io/otel/metric v0.32.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
	golang.org/x/net v0.0.0-20220909164309-bea034e7d591 // indirect
	golang.org/x/oauth2 v0.0.0-20220909003341-f21342109be1 // indirect
	golang.org/x/sys v0.0.0-20220808155132-1c4a2a72c664 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220915135415-7fd63a7952de // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => ../../internal/common

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal => ../../internal/coreinternal

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent => ../../internal/sharedcomponent

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus => ../../pkg/translator/opencensus
