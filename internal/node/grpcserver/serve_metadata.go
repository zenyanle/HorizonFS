package grpcserver

import (
	"HorizonFS/internal/node/metadata"
	"HorizonFS/internal/node/service"
	"HorizonFS/pkg/logger"
	"google.golang.org/grpc"
	"log"
)

func ServeMetadata(store *metadata.Metastore) {
	grpcServer := grpc.NewServer()

	metadataService := &service.MetadataService{
		Store: store,
	}
	proto.RegisterTransfrServiceServer(grpcServer, metadataService)
	log.Println("Server is running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	logger.Info("Server start")
}
