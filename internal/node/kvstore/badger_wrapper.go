package kvstore

import (
	"bytes"
	"fmt"
	"io"

	badger "github.com/dgraph-io/badger/v4"
)

// BadgerStore 封装badger数据库操作
type BadgerStore struct {
	db *badger.DB
}

// NewBadgerStore 创建一个新的BadgerStore
func NewBadgerStore(path string) (*BadgerStore, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, fmt.Errorf("failed to open badger DB: %w", err)
	}
	return &BadgerStore{db: db}, nil
}

// Close 关闭数据库
func (s *BadgerStore) Close() error {
	return s.db.Close()
}

// GetReader 获取指定ID的io.Reader
// 注意：此方法会将整个值加载到内存中
func (s *BadgerStore) GetReader(id string) (io.Reader, int, error) {
	var valueBytes []byte

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}

		// 从item中复制值
		err = item.Value(func(val []byte) error {
			valueBytes = make([]byte, len(val))
			copy(valueBytes, val)
			return nil
		})
		return err
	})

	if err != nil {
		return nil, -1, err
	}

	// 使用bytes.Reader创建一个io.Reader
	return bytes.NewReader(valueBytes), len(valueBytes), nil
}

// SetValue 存储值到数据库
func (s *BadgerStore) SetValue(id string, value []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(id), value)
	})
}

// SetFromReader 从io.Reader读取数据并存储到指定ID
// 注意：此方法会将整个Reader内容加载到内存中
func (s *BadgerStore) SetFromReader(id string, reader io.Reader) (int64, error) {
	// 读取Reader中的所有数据
	buf := new(bytes.Buffer)
	written, err := io.Copy(buf, reader)
	if err != nil {
		return 0, fmt.Errorf("读取数据失败: %w", err)
	}

	// 存储到数据库
	err = s.SetValue(id, buf.Bytes())
	if err != nil {
		return 0, fmt.Errorf("存储数据失败: %w", err)
	}

	return written, nil
}
