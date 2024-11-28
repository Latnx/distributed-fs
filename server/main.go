package main

import (
	"context"
	"log"
	"net"

	pb "grpc-distributed-fs/proto/fs"
	"grpc-distributed-fs/storage"

	"google.golang.org/grpc"
)

type serverImpl struct {
	pb.UnimplementedFileSystemServer
	db *storage.FileDB
}

// 创建新文件（包括目录路径）
func (s *serverImpl) WriteFile(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	err := s.db.WriteFile(req.Filename, "/", req.Data) // 假设路径是根目录
	if err != nil {
		log.Printf("Error inserting file: %v", err)
		return nil, err
	}
	return &pb.WriteResponse{}, nil
}

// 读取文件（包括路径）
func (s *serverImpl) ReadFile(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	data, err := s.db.ReadFile(req.Filename, "/") // 假设路径是根目录
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return nil, err
	}
	return &pb.ReadResponse{Data: data}, nil
}

// 删除文件（包括路径）
func (s *serverImpl) DeleteFile(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := s.db.DeleteFile(req.Filename, "/") // 假设路径是根目录
	if err != nil {
		log.Printf("Error deleting file: %v", err)
		return nil, err
	}
	return &pb.DeleteResponse{}, nil
}

func main() {
	// 初始化数据库
	db := storage.NewFileDB("data")
	defer db.Close()

	// 启动 gRPC 服务
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterFileSystemServer(grpcServer, &serverImpl{db: db})

	log.Println("Server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
