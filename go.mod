module HorizonFS

// Consider updating the Go version if your toolchain supports it, e.g., go 1.21 or 1.22
go 1.23.0

toolchain go1.24.1

require (
	// Updated etcd internal packages to latest stable v3.5.x release train
	go.etcd.io/etcd/client/pkg/v3 v3.5.14
	go.etcd.io/etcd/raft/v3 v3.5.14 // Aligned with server/client for compatibility with example code structure
	go.etcd.io/etcd/server/v3 v3.5.14

	// Updated Zap to latest stable
	go.uber.org/zap v1.27.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect; indirect // Note: v0.3.1 is often used, check if newer needed
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/prometheus/client_golang v1.20.4 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.63.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/etcd/api/v3 v3.5.14 // indirect
	go.etcd.io/etcd/pkg/v3 v3.5.14 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

require (
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.59.0
)

require (
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
)
