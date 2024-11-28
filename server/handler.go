package main

import (
	"context"

	pb "grpc-distributed-fs/proto/fs"
	"grpc-distributed-fs/storage"
)

type FileSystemServer struct {
	pb.UnimplementedFileSystemServer
	storage *storage.LocalStorage
}

func NewFileSystemServer(baseDir string) *FileSystemServer {
	return &FileSystemServer{
		storage: storage.NewLocalStorage(baseDir),
	}
}

func (s *FileSystemServer) WriteFile(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	err := s.storage.WriteFile(req.Filename, req.Data)
	if err != nil {
		return nil, err
	}
	return &pb.WriteResponse{Message: "File written successfully"}, nil
}

func (s *FileSystemServer) ReadFile(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	data, err := s.storage.ReadFile(req.Filename)
	if err != nil {
		return nil, err
	}
	return &pb.ReadResponse{Data: data}, nil
}

func (s *FileSystemServer) DeleteFile(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := s.storage.DeleteFile(req.Filename)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteResponse{Message: "File deleted successfully"}, nil
}

func (s *FileSystemServer) ListFiles(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	files, err := s.storage.ListFiles()
	if err != nil {
		return nil, err
	}
	return &pb.ListResponse{Files: files}, nil
}
