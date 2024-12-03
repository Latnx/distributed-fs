package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"grpc-distributed-fs/metadata"
	pb "grpc-distributed-fs/proto/fs"
)

type Client struct {
	pb.FileSystemClient
}

// 上传文件
func UploadFile(clients [](*Client), tree *metadata.FileTree, command []string) {
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

	// 使用树结构的当前路径作为父目录路径
	parentPath := tree.Current.Metadata.Name

	// 调用服务端上传
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = clients[0].WriteFile(ctx, &pb.WriteRequest{
		Filename: parentPath + filename,
		Data:     data,
	})
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return
	}

	err = tree.AddFile(filename, int64(len(data))) // 更新元数据
	if err != nil {
		fmt.Printf("Failed to update metadata: %v\n", err)
		return
	}

	fmt.Printf("File '%s' uploaded successfully.\n", filename)
}

// 下载文件
func DownloadFile(clients [](*Client), tree *metadata.FileTree, command []string) {
	if len(command) < 2 {
		fmt.Println("Usage: download <file-name>")
		return
	}
	filename := command[1]

	// 使用树结构的当前路径作为父目录路径
	parentPath := tree.Current.Metadata.Name

	// 请求服务器读取文件
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := clients[0].ReadFile(ctx, &pb.ReadRequest{Filename: parentPath + filename})
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
func RemoveFile(clients [](*Client), tree *metadata.FileTree, command []string) {
	if len(command) < 2 {
		fmt.Println("Usage: rm <file-name>")
		return
	}
	filename := command[1]

	// 使用树结构的当前路径作为父目录路径
	parentPath := tree.Current.Metadata.Name

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := clients[0].DeleteFile(ctx, &pb.DeleteRequest{Filename: parentPath + filename})
	if err != nil {
		fmt.Printf("Failed to delete file: %v\n", err)
		return
	}

	err = tree.RemoveFile(filename) // 更新元数据
	if err != nil {
		fmt.Printf("Failed to update metadata: %v\n", err)
		return
	}

	fmt.Printf("File '%s' deleted successfully.\n", filename)
}

// 工具函数
func getFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
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
	// fmt.Printf("Metadata for '%s':\n", filename)
	// fmt.Printf("  Size: %d bytes\n", meta.Size)
	// fmt.Printf("  Creation Time: %s\n", meta.CreationTime)
	// fmt.Printf("  Modification Time: %s\n", meta.ModificationTime)
	fmt.Printf("Metadata:\n %+v\n", meta)
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
