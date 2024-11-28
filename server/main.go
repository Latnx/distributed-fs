package main

import (
	"log"
	"net"

	pb "grpc-distributed-fs/proto/fs"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	fsServer := NewFileSystemServer("./data")
	pb.RegisterFileSystemServer(server, fsServer)

	log.Println("Server is running on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
