package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"grpc-distributed-fs/metadata"
	pb "grpc-distributed-fs/proto/fs"

	"google.golang.org/grpc"
)

// 初始化 gRPC 客户端
func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{pb.NewFileSystemClient(conn)}
}
func NewConn(port string) *Client {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	client := NewClient(conn)
	return client
}
func main() {
	clients := make([](*Client), 2)
	clients[0] = NewConn(":50051")
	clients[1] = NewConn(":50052")

	tree := metadata.NewFileTree() // 初始化文件树

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the Distributed File System!")
	for {
		fmt.Printf("%s> ", tree.Current.Metadata.Name)
		input, _ := reader.ReadString('\n')
		command := strings.Fields(strings.TrimSpace(input))
		if len(command) == 0 {
			continue
		}

		switch command[0] {
		case "ls":
			fmt.Println("Contents:", tree.Ls())
		case "cd":
			ChangeDirectory(tree, command)
		case "mkdir":
			MakeDirectory(tree, command)
		case "upload":
			UploadFile(clients, tree, command, 512)
		case "download":
			DownloadFile(clients, tree, command)
		case "rm":
			RemoveFile(clients, tree, command)
		case "meta":
			ViewMetadata(tree, command)
		case "exit":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Unknown command:", command[0])
		}
	}
}
