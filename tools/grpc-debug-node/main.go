package main

import (
	"HorizonFS/internal/node/kvstore"
	"HorizonFS/internal/node/metadata"
	"HorizonFS/internal/node/raft"
	"HorizonFS/internal/node/service"
	"HorizonFS/pkg/logger"
	"HorizonFS/pkg/proto"
	"flag"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"google.golang.org/grpc"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	kvport := flag.Int("port", 9121, "key-value server port")
	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()

	proposeC := make(chan string)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)

	// raft provides a commit stream for the proposals from the http api
	var kvs *metadata.Metastore
	getSnapshot := func() ([]byte, error) { return kvs.GetSnapshot() }
	commitC, errorC, snapshotterReady := raft.NewRaftNode(*id, strings.Split(*cluster, ","), *join, getSnapshot, proposeC, confChangeC)

	kvs = metadata.NewMetaSore(<-snapshotterReady, proposeC, commitC, errorC)

	// the key-value http handler will propose updates to raft
	// serveHttpKVAPI(kvs, *kvport, confChangeC, errorC)
	grpcServer := grpc.NewServer()

	metadataService := &service.MetadataService{
		Store: kvs,
	}

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

	proto.RegisterMetadataServiceServer(grpcServer, metadataService)
	proto.RegisterDataNodeServiceServer(grpcServer, datanodeService)

	logger.Println("Server is running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
	logger.Info("Server start")
}
