package oss

import (
	"io"
	"time"
)

// 屏蔽不同平台的oss接口
type storager interface {
	GetObject(bucket string, key string) ([]byte, error)
	GetFile(bucket string, key string, localFile string) error
	PutObjectWithMeta(bucket string, key string, data []byte, metadata *Metadata) error
	PutFileWithMeta(bucket string, key string, filePath string, metadata *Metadata) error
	PutReaderWithMeta(bucket string, key string, reader io.Reader, metadata *Metadata) error
	ListObjects(bucket string, prefix string) ([]FileMeta, error)
	DeleteObject(bucket string, key string) error
	CopyObject(bucket string, srcKey string, destKey string, metadata *Metadata) error
	SetObjectMeta(bucket string, key string, metadata *Metadata) error
	IsObjectExist(bucket string, key string) (bool, error)
	GetDirToken(bucket string, remoteDir string, expires time.Duration) (*StsTokenInfo, error)
	GetDirTokenRead(bucket string, remoteDir string, expires time.Duration) (*StsTokenInfo, error)
	PresignObject(bucket string, key string, expired time.Duration) (string, error)
	SignFile(bucket string, key string, expired time.Duration) (string, error)
	GetObjectMeta(bucket string, key string) (*FileMeta, error)
}
