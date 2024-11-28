package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"grpc-distributed-fs/metadata"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := NewClient(conn)
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
			UploadFile(client, tree, command)
		case "download":
			DownloadFile(client, tree, command)
		case "rm":
			RemoveFile(client, tree, command)
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
