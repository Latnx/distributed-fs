package main

import (
	"context"
	"log"
	"time"

	pb "grpc-distributed-fs/proto/fs"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileSystemClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 上传文件
	_, err = client.WriteFile(ctx, &pb.WriteRequest{
		Filename: "example2.txt",
		Data:     []byte("Hello, gRPC!"),
	})
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}
	log.Println("File uploaded successfully")

	// 列出文件
	listResp, err := client.ListFiles(ctx, &pb.ListRequest{})
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}
	log.Printf("Files: %v", listResp.Files)
}
