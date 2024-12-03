package metadata

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
)

// SplitFile 分割文件并生成分片信息
func SplitFile(filePath string, chunkSize int64, storageDir string) (*FileMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %v", err)
	}

	// 初始化文件元数据
	metadata := &FileMetadata{
		Name:             fileInfo.Name(),
		IsDirectory:      fileInfo.IsDir(),
		Size:             fileInfo.Size(),
		CreationTime:     time.Now(), // 假设当前时间为文件创建时间
		ModificationTime: fileInfo.ModTime(),
		Chunks:           []FileChunk{},
	}

	// 分片处理
	buffer := make([]byte, chunkSize)
	chunkID := 0
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		if bytesRead == 0 {
			break
		}

		// 计算分片校验值
		checksum := md5.Sum(buffer[:bytesRead])
		checksumStr := hex.EncodeToString(checksum[:])

		// 创建分片文件名和存储路径
		chunkFileName := fmt.Sprintf("%s_chunk_%d", fileInfo.Name(), chunkID)
		storagePath := fmt.Sprintf("%s/%s", storageDir, chunkFileName)

		// 写入分片到存储位置
		err = writeChunkToFile(storagePath, buffer[:bytesRead])
		if err != nil {
			return nil, fmt.Errorf("failed to write chunk: %v", err)
		}

		// 创建分片元数据
		chunk := FileChunk{
			ChunkID:         fmt.Sprintf("%s-%d", fileInfo.Name(), chunkID),
			FileID:          fileInfo.Name(),
			ChunkNumber:     chunkID,
			OriginalName:    fileInfo.Name(),
			Size:            int64(bytesRead),
			Checksum:        checksumStr,
			StorageLocation: storagePath,
			Replicas:        []string{},
		}

		metadata.Chunks = append(metadata.Chunks, chunk)
		chunkID++
	}

	return metadata, nil
}

// writeChunkToFile 将分片写入指定位置
func writeChunkToFile(path string, data []byte) error {
	chunkFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer chunkFile.Close()

	_, err = chunkFile.Write(data)
	if err != nil {
		return err
	}
	return nil
}
