package metadata

import (
	"errors"
	"time"
)

type FileChunk struct {
	ChunkID         string   // 分片唯一标识符（如 UUID）
	FileID          string   // 所属文件唯一标识符
	ChunkNumber     int      // 分片编号，从 0 开始
	OriginalName    string   // 原始文件名
	Size            int64    // 分片大小（字节）
	Checksum        string   // 分片校验值，用于数据完整性校验
	StorageLocation string   // 分片存储位置（路径、URI 等）
	Replicas        []string // 副本存储位置列表
}

// 文件元数据
type FileMetadata struct {
	Name             string
	IsDirectory      bool
	Size             int64
	CreationTime     time.Time
	ModificationTime time.Time
	Chunks           []FileChunk // 分片信息列表
}

// 文件树节点
type FileNode struct {
	Metadata *FileMetadata
	Children map[string]*FileNode
	Parent   *FileNode
}

// 文件树结构
type FileTree struct {
	Root    *FileNode
	Current *FileNode
}

// 初始化文件树
func NewFileTree() *FileTree {
	root := &FileNode{
		Metadata: &FileMetadata{
			Name:         "/",
			IsDirectory:  true,
			CreationTime: time.Now(),
		},
		Children: make(map[string]*FileNode),
	}
	return &FileTree{Root: root, Current: root}
}

// 创建目录
func (t *FileTree) Mkdir(name string) error {
	if _, exists := t.Current.Children[name]; exists {
		return errors.New("directory already exists")
	}
	node := &FileNode{
		Metadata: &FileMetadata{
			Name:         name,
			IsDirectory:  true,
			CreationTime: time.Now(),
		},
		Children: make(map[string]*FileNode),
		Parent:   t.Current,
	}
	t.Current.Children[name] = node
	return nil
}

// 切换目录
func (t *FileTree) Cd(name string) error {
	if name == ".." {
		if t.Current.Parent != nil {
			t.Current = t.Current.Parent
			return nil
		}
		return errors.New("already at root")
	}
	if node, exists := t.Current.Children[name]; exists && node.Metadata.IsDirectory {
		t.Current = node
		return nil
	}
	return errors.New("directory not found")
}

// 列出目录
func (t *FileTree) Ls() []string {
	var entries []string
	for name := range t.Current.Children {
		entries = append(entries, name)
	}
	return entries
}

// 添加文件
func (t *FileTree) AddFile(name string, size int64) error {
	if _, exists := t.Current.Children[name]; exists {
		return errors.New("file already exists")
	}
	t.Current.Children[name] = &FileNode{
		Metadata: &FileMetadata{
			Name:             name,
			Size:             size,
			CreationTime:     time.Now(),
			ModificationTime: time.Now(),
			Chunks:           make([]FileChunk, 1024),
		},
		Parent: t.Current,
	}
	return nil
}

// 删除文件
func (t *FileTree) RemoveFile(name string) error {
	if _, exists := t.Current.Children[name]; !exists {
		return errors.New("file not found")
	}
	delete(t.Current.Children, name)
	return nil
}

// 获取文件元数据
func (t *FileTree) GetFileMetadata(name string) (*FileMetadata, error) {
	node, exists := t.Current.Children[name]
	if !exists {
		return nil, errors.New("file not found")
	}
	return node.Metadata, nil
}
