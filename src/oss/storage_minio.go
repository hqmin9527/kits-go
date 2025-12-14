package oss

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/hqmin9527/kits-go/src/logger"
	"github.com/hqmin9527/kits-go/src/utils"
	"github.com/minio/minio-go/v6"
	"github.com/minio/minio-go/v6/pkg/credentials"
	"github.com/pkg/errors"
)

type minStorager struct {
	config *Config
	client *minio.Client
}

func newMinioStorager(c *Config) (*minStorager, error) {
	logger.Info("create minIo oss client")
	endpoint := utils.If(c.Internal, c.EndpointInner, c.Endpoint)
	client, err := minio.New(endpoint, c.AccessKeyId, c.AccessKeySecret, false)
	if err != nil {
		return nil, errors.Wrap(err, "create minio client")
	}
	return &minStorager{config: c, client: client}, nil
}

func (m *minStorager) GetObject(bucket string, key string) (byte []byte, er error) {
	obj, err := m.client.GetObject(bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, obj)
	return buf.Bytes(), nil
}

func (m *minStorager) GetFile(bucket string, key string, localFile string) error {
	return m.client.FGetObject(bucket, key, localFile, minio.GetObjectOptions{})
}

func (m *minStorager) PutObject(bucket string, key string, data []byte, metadata map[string]string) error {
	_, err := m.client.PutObject(bucket, key, bytes.NewBuffer(data), int64(len(data)), mapToPutObjOptions(metadata))
	return err
}

func (m *minStorager) PutObjectWithMeta(bucket string, key string, data []byte, metadata *Metadata) error {
	_, err := m.client.PutObject(bucket, key, bytes.NewBuffer(data), int64(len(data)), metadataToPutObjOptions(metadata))
	return err
}

func (m *minStorager) PutFile(bucket string, key string, srcFile string, metadata map[string]string) error {
	_, err := m.client.FPutObject(bucket, key, srcFile, mapToPutObjOptions(metadata))
	return err
}

func (m *minStorager) PutFileWithMeta(bucket string, key string, srcFile string, metadata *Metadata) error {
	_, err := m.client.FPutObject(bucket, key, srcFile, metadataToPutObjOptions(metadata))
	return err
}

func (m *minStorager) PutReaderWithMeta(bucket string, key string, r io.Reader, metadata *Metadata) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return m.PutObjectWithMeta(bucket, key, data, metadata)
}

func (m *minStorager) ListObjects(bucket string, prefix string) ([]FileMeta, error) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	res := make([]FileMeta, 0, 32)
	// List all objects from a bucket-name with a matching prefix.
	for object := range m.client.ListObjectsV2(bucket, prefix, true, doneCh) {
		if object.Err != nil {
			return res, object.Err
		}
		// protect
		if len(res) >= 100000 {
			return res, errors.New("too much data, break")
		}
		res = append(res, *objectInfoToContent(&object))
	}
	return res, nil
}

func objectInfoToContent(obj *minio.ObjectInfo) *FileMeta {
	res := &FileMeta{
		Key:          obj.Key,
		Size:         obj.Size,
		ETag:         obj.ETag,
		LastModified: obj.LastModified,
	}
	res.Metadata = Metadata{
		ContentType:        obj.ContentType,
		ContentEncoding:    obj.Metadata.Get("content-encoding"),
		ContentDisposition: obj.Metadata.Get("content-disposition"),
	}
	return res
}

func (m *minStorager) DeleteObject(bucket string, key string) error {
	return m.client.RemoveObject(bucket, key)
}

func (m *minStorager) CopyObject(bucket string, srcKey string, destKey string, metadata *Metadata) error {
	// Source object
	src := minio.NewSourceInfo(bucket, srcKey, nil)
	// _ = src.SetModifiedSinceCond(time.Now())
	// Destination object
	userMeta := metadataToUserMeta(metadata)
	dst, err := minio.NewDestinationInfo(bucket, destKey, nil, userMeta)
	if err != nil {
		return err
	}
	if err = m.client.CopyObject(dst, src); err != nil {
		return err
	}
	return nil
}

func (m *minStorager) SetObjectAcl(bucket string, key string, acl ACL) error {
	// return errors.New("minio not support SetObjectAcl")
	return nil
}

func (m *minStorager) SetObjectMeta(bucket string, key string, metadata *Metadata) error {
	// return errors.New("minio not support SetObjectMeta")
	return nil
}

func (m *minStorager) IsObjectExist(bucket string, key string) (bool, error) {
	res, err := m.GetObjectMeta(bucket, key)
	switch err := err.(type) {
	case minio.ErrorResponse:
		if err.Code == "NoSuchKey" {
			return false, nil
		} else {
			return false, err
		}
	default:
		return res != nil, err
	}
}

func (m *minStorager) GetDirToken(bucket string, remoteDir string, expires time.Duration) (*StsTokenInfo, error) {
	li, err := credentials.NewSTSAssumeRole("http://"+m.config.Endpoint, credentials.STSAssumeRoleOptions{
		AccessKey:       "rw_client",
		SecretKey:       "#$infi0831",
		DurationSeconds: int(expires.Seconds()),
	})
	if err != nil {
		return nil, err
	}
	to, err := li.Get()
	if err != nil {
		return nil, err
	}
	// url to StsTokenInfo
	return &StsTokenInfo{
		Provider:        providerMinio,
		AccessKeyID:     to.AccessKeyID,
		AccessKeySecret: to.SecretAccessKey,
		StsToken:        to.SessionToken,
		Bucket:          bucket,
		Region:          m.config.Region,
		Expire:          time.Now().Add(expires).UnixMilli(),
		UploadPath:      remoteDir,
		Host:            fmt.Sprintf("%s://%s", m.config.Protocol, m.config.Host),
		Endpoint:        m.config.Endpoint,
	}, nil
}

func (m *minStorager) GetDirTokenRead(bucket string, remoteDir string, expires time.Duration) (*StsTokenInfo, error) {
	return m.GetDirToken(bucket, remoteDir, expires)
}

func (m *minStorager) PresignObject(bucket string, key string, expires time.Duration) (string, error) {
	return fmt.Sprintf("%s://%s/%s/%s", m.config.Protocol, m.config.Host, bucket, key), nil
	// url, err := m.client.PresignedPutObject(bucket, key, expires)
	// if err != nil {
	// 	return "", err
	// }
	//
	// query := url.Query().Encode()
	// query = strings.Replace(query, "\u0026", "&", -1)
	// uploadUrl := fmt.Sprintf("http://%s%s?%s", m.config.Endpoint, url.EscapedPath(), query)
	// return uploadUrl, nil
}

func (m *minStorager) SignFile(bucket string, key string, expires time.Duration) (string, error) {
	res := ""
	if m.config.DownloadDomain != "" {
		res = fmt.Sprintf("https://%s/%s/%s", m.config.DownloadDomain, bucket, key)
	} else {
		res = fmt.Sprintf("%s://%s/%s/%s", m.config.Protocol, m.config.Host, bucket, key)
	}
	return res, nil
}

func (m *minStorager) GetObjectMeta(bucket string, key string) (*FileMeta, error) {
	objInfo, err := m.client.GetObjectACL(bucket, key)
	if objInfo == nil {
		return nil, err
	}
	return objectInfoToContent(objInfo), err
}

func metadataToPutObjOptions(metadata *Metadata) minio.PutObjectOptions {
	ops := minio.PutObjectOptions{}
	if metadata != nil {
		ops.ContentType = metadata.ContentType
		ops.ContentEncoding = metadata.ContentEncoding
		ops.ContentDisposition = metadata.ContentDisposition
	}
	return ops
}

func mapToPutObjOptions(metadata map[string]string) minio.PutObjectOptions {
	ops := minio.PutObjectOptions{}
	if metadata == nil || len(metadata) == 0 {
		return ops
	}
	for key, value := range metadata {
		switch key {
		case "ContentType":
			ops.ContentType = value
		case "ContentEncoding":
			ops.ContentEncoding = value
		case "ContentDisposition":
			ops.ContentDisposition = value
		case "ContentLanguage":
			ops.ContentLanguage = value
		case "CacheControl":
			ops.CacheControl = value
		default:
			// do nothing
		}
	}
	return ops
}

func metadataToUserMeta(metadata *Metadata) map[string]string {
	if metadata == nil {
		return nil
	}
	res := make(map[string]string)
	if metadata.ContentType != "" {
		res["Content-Type"] = metadata.ContentType
	}
	if metadata.ContentEncoding != "" {
		res["Content-Encoding"] = metadata.ContentEncoding
	}
	if metadata.ContentDisposition != "" {
		res["Content-Disposition"] = metadata.ContentDisposition
	}
	return res
}
