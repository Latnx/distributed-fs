package main

import (
	"context"
	"log"

	pb "grpc-distributed-fs/proto/fs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileSystemClient(conn)

	// 上传文件
	_, err = client.WriteFile(context.Background(), &pb.WriteRequest{
		Filename: "example2.txt",
		Data:     []byte("Hello, gRPC!!!!!"),
	})
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}
	log.Println("File uploaded successfully")

	// 列出文件
	listResp, err := client.ListFiles(context.Background(), &pb.ListRequest{})
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}
	log.Printf("Files: %v", listResp.Files)
}
