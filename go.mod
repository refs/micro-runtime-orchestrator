module github.com/refs/whatever

go 1.13

replace google.golang.org/grpc => google.golang.org/grpc v1.24.0

require (
	github.com/asim/go-micro/plugins/client/grpc/v3 v3.0.0-20210217182006-0f0ace1a44a9 // indirect
	github.com/asim/go-micro/v3 v3.5.0
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.5
	github.com/grpc-ecosystem/grpc-gateway v1.12.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.2.0 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/micro/go-micro/v2 v2.4.0
	github.com/thejerf/suture v4.0.0+incompatible
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/protobuf v1.25.0
)
