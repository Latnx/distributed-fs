package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"grpc-distributed-fs/metadata"
	pb "grpc-distributed-fs/proto/fs"

	"github.com/davecgh/go-spew/spew"
)

type Client struct {
	pb.FileSystemClient
}

func UploadFile(clients [](*Client), tree *metadata.FileTree, command []string, chunkSize int64) {
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
	parentPath := tree.Current.Metadata.Name
	fileID := parentPath + filename // 生成唯一文件标识符
	var fileChunks []metadata.FileChunk

	// 分片逻辑
	for i := int64(0); i < int64(len(data)); i += chunkSize {
		end := i + chunkSize
		if end > int64(len(data)) {
			end = int64(len(data))
		}

		chunkData := data[i:end]
		// 生成唯一的分片标识符
		chunkID := fmt.Sprintf("%s_chunk_%d", fileID, i/chunkSize)

		storageLocation := int(i) % len(clients)

		fileChunk := metadata.FileChunk{
			ChunkID:         chunkID,
			FileID:          fileID,
			ChunkNumber:     int(i / chunkSize),
			OriginalName:    filename,
			Size:            int64(len(chunkData)),
			StorageLocation: storageLocation,
			Replicas:        []string{},
		}

		fileChunks = append(fileChunks, fileChunk)

		// 上传分片
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err = clients[storageLocation].WriteFile(ctx, &pb.WriteRequest{
			Filename: chunkID,
			Data:     chunkData,
		})
		if err != nil {
			fmt.Printf("Failed to upload chunk %d: %v\n", i/chunkSize, err)
			return
		}

		fmt.Printf("Uploaded chunk %d successfully.\n", i/chunkSize)
	}

	// 更新元数据到文件树
	err = tree.AddFile(&metadata.FileMetadata{
		Name:             filename,
		IsDirectory:      false,
		Size:             int64(len(data)),
		CreationTime:     time.Now(),
		ModificationTime: time.Now(),
		Chunks:           fileChunks, // 记录所有分片
	})
	if err != nil {
		fmt.Printf("Failed to update metadata: %v\n", err)
		return
	}

	fmt.Printf("File '%s' uploaded successfully in %d chunks.\n", filename, len(fileChunks))
}

func DownloadFile(clients [](*Client), tree *metadata.FileTree, command []string) {
	if len(command) < 2 {
		fmt.Println("Usage: download <file-name>")
		return
	}
	filename := command[1]

	// 获取文件元数据
	fileMetadata, err := tree.GetFileMetadata(filename)
	if err != nil {
		fmt.Printf("Failed to find file metadata: %v\n", err)
		return
	}

	if fileMetadata.IsDirectory {
		fmt.Println("Cannot download a directory.")
		return
	}

	var fileData []byte
	for _, chunk := range fileMetadata.Chunks {
		// 下载分片
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		clientIndex := chunk.StorageLocation % len(clients)
		resp, err := clients[clientIndex].ReadFile(ctx, &pb.ReadRequest{Filename: chunk.ChunkID})
		if err != nil {
			fmt.Printf("Failed to download chunk %d: %v\n", chunk.ChunkNumber, err)
			return
		}

		fmt.Printf("Downloaded chunk %d successfully.\n", chunk.ChunkNumber)
		fileData = append(fileData, resp.Data...)
	}

	// 将拼接后的数据保存到本地
	err = ioutil.WriteFile(filename, fileData, 0644)
	if err != nil {
		fmt.Printf("Failed to save file locally: %v\n", err)
		return
	}

	fmt.Printf("File '%s' downloaded successfully.\n", filename)
}

func RemoveFile(clients [](*Client), tree *metadata.FileTree, command []string) {
	if len(command) < 2 {
		fmt.Println("Usage: rm <file-name>")
		return
	}
	filename := command[1]

	// 获取文件元数据
	fileMetadata, err := tree.GetFileMetadata(filename)
	if err != nil {
		fmt.Printf("Failed to find file metadata: %v\n", err)
		return
	}

	if fileMetadata.IsDirectory {
		fmt.Println("Cannot remove a directory with this function.")
		return
	}

	// 删除所有分片
	for _, chunk := range fileMetadata.Chunks {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		clientIndex := chunk.StorageLocation % len(clients)
		_, err := clients[clientIndex].DeleteFile(ctx, &pb.DeleteRequest{Filename: chunk.ChunkID})
		if err != nil {
			fmt.Printf("Failed to delete chunk %d: %v\n", chunk.ChunkNumber, err)
			return
		}

		fmt.Printf("Deleted chunk %d successfully.\n", chunk.ChunkNumber)
	}

	// 更新元数据
	err = tree.RemoveFile(filename)
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
	spew.Dump(meta)
	// fmt.Printf("Metadata for '%s':\n", filename)
	// fmt.Printf("  Size: %d bytes\n", meta.Size)
	// fmt.Printf("  Creation Time: %s\n", meta.CreationTime)
	// fmt.Printf("  Modification Time: %s\n", meta.ModificationTime)
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
