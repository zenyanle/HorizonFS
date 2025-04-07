package main

import (
	"HorizonFS/internal/node/kvstore"
	"HorizonFS/internal/node/service"
	"HorizonFS/pkg/logger"
	"HorizonFS/pkg/proto"
	"flag"
	"google.golang.org/grpc"
	"net"
	"path/filepath"
	"strconv"
)

func main() {
	kvport := flag.Int("port", 22380, "key-value server port")
	flag.Parse()
	
	// the key-value http handler will propose updates to raft
	// serveHttpKVAPI(kvs, *kvport, confChangeC, errorC)
	grpcServer := grpc.NewServer()

	kv, err := kvstore.NewBadgerStore(filepath.Join("storage", "test.db"))
	if err != nil {
		logger.Fatal(err)
	}

	datanodeService := &service.DataNodeService{
		Kv: kv,
	}

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*kvport))
	if err != nil {
		logger.Fatal(err)
	}

	proto.RegisterDataNodeServiceServer(grpcServer, datanodeService)

	logger.Printf("Server is running on :%d", *kvport)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
	logger.Info("Server start")
}
