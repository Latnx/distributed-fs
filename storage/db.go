package storage

import (
	"log"

	"github.com/dgraph-io/badger/v3"
)

type FileDB struct {
	db *badger.DB
}

// 初始化数据库
func NewFileDB(dbPath string) *FileDB {
	opts := badger.DefaultOptions(dbPath).WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("Failed to open BadgerDB: %v", err)
	}
	return &FileDB{db: db}
}

// 写入文件，使用完整路径作为唯一标识符
func (fdb *FileDB) WriteFile(filename, parentPath string, data []byte) error {
	// 生成完整路径作为键
	filePath := parentPath + "/" + filename
	return fdb.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(filePath), data)
	})
}

// 读取文件，使用完整路径作为键
func (fdb *FileDB) ReadFile(filename, parentPath string) ([]byte, error) {
	filePath := parentPath + "/" + filename
	var data []byte
	err := fdb.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(filePath))
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		return err
	})
	return data, err
}

// 删除文件，使用完整路径作为键
func (fdb *FileDB) DeleteFile(filename, parentPath string) error {
	filePath := parentPath + "/" + filename
	return fdb.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(filePath))
	})
}

// 关闭数据库
func (fdb *FileDB) Close() {
	fdb.db.Close()
}
