package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"grpc-distributed-fs/metadata"
	pb "grpc-distributed-fs/proto/fs"

	"google.golang.org/grpc"
)

type Client struct {
	pb.FileSystemClient
}

// 初始化 gRPC 客户端
func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{pb.NewFileSystemClient(conn)}
}

// 上传文件
func UploadFile(client *Client, tree *metadata.FileTree, command []string) {
	if len(command) < 2 {
		fmt.Println("Usage: upload <local-file-path>")
		return
	}
	localPath := command[1]
	data, err := ioutil.ReadFile(localPath)
	if err != nil {
		fmt.Printf("Failed to read local file: %v\n", err)
		return
	}
	filename := getFileName(localPath)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.WriteFile(ctx, &pb.WriteRequest{
		Filename: filename,
		Data:     data,
	})
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return
	}

	err = tree.AddFile(filename, int64(len(data)))
	if err != nil {
		fmt.Printf("Failed to update metadata: %v\n", err)
		return
	}

	fmt.Printf("File '%s' uploaded successfully.\n", filename)
}

// 下载文件
func DownloadFile(client *Client, tree *metadata.FileTree, command []string) {
	if len(command) < 2 {
		fmt.Println("Usage: download <file-name>")
		return
	}
	filename := command[1]
	_, err := tree.GetFileMetadata(filename)
	if err != nil {
		fmt.Printf("File not found in namespace: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.ReadFile(ctx, &pb.ReadRequest{Filename: filename})
	if err != nil {
		fmt.Printf("Failed to download file: %v\n", err)
		return
	}

	err = ioutil.WriteFile(filename, resp.Data, 0644)
	if err != nil {
		fmt.Printf("Failed to save file locally: %v\n", err)
		return
	}

	fmt.Printf("File '%s' downloaded successfully.\n", filename)
}

// 删除文件
func RemoveFile(client *Client, tree *metadata.FileTree, command []string) {
	if len(command) < 2 {
		fmt.Println("Usage: rm <file-name>")
		return
	}
	filename := command[1]
	_, err := tree.GetFileMetadata(filename)
	if err != nil {
		fmt.Printf("File not found in namespace: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.DeleteFile(ctx, &pb.DeleteRequest{Filename: filename})
	if err != nil {
		fmt.Printf("Failed to delete file: %v\n", err)
		return
	}

	delete(tree.Current.Children, filename)
	fmt.Printf("File '%s' deleted successfully.\n", filename)
}

// 查看文件元数据
func ViewMetadata(tree *metadata.FileTree, command []string) {
	if len(command) < 2 {
		fmt.Println("Usage: meta <file-name>")
		return
	}
	filename := command[1]
	meta, err := tree.GetFileMetadata(filename)
	if err != nil {
		fmt.Printf("File not found: %v\n", err)
		return
	}
	fmt.Printf("Metadata for '%s':\n", filename)
	fmt.Printf("  Size: %d bytes\n", meta.Size)
	fmt.Printf("  Creation Time: %s\n", meta.CreationTime)
	fmt.Printf("  Modification Time: %s\n", meta.ModificationTime)
}

// 切换目录
func ChangeDirectory(tree *metadata.FileTree, command []string) {
	if len(command) < 2 {
		fmt.Println("Usage: cd <directory>")
		return
	}
	err := tree.Cd(command[1])
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// 创建目录
func MakeDirectory(tree *metadata.FileTree, command []string) {
	if len(command) < 2 {
		fmt.Println("Usage: mkdir <directory>")
		return
	}
	err := tree.Mkdir(command[1])
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// 工具函数
func getFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
