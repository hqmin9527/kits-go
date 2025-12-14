package oss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	sts "github.com/alibabacloud-go/sts-20150401/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hqmin9527/kits-go/src/logger"
	"github.com/hqmin9527/kits-go/src/utils"
	"github.com/pkg/errors"
)

type aliStorager struct {
	config    *Config
	client    *oss.Client
	stsClient *sts.Client // 安全授权服务客户端
}

func newAliStorager(c *Config) (*aliStorager, error) {
	logger.Info("create aliyun oss client")
	endpoint := utils.If(c.Internal, c.EndpointInner, c.Endpoint)
	client, err := oss.New(endpoint, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return nil, errors.Wrap(err, "create oss client")
	}
	stsClient, err := createAliStsClient(c.StsEndpoint, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return nil, errors.Wrap(err, "create aliyun sts client")
	}
	return &aliStorager{config: c, client: client, stsClient: stsClient}, nil
}

func createAliStsClient(endpoint string, accessKeyId string, accessKeySecret string) (client *sts.Client, err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}
	client, err = sts.NewClient(config)
	return client, err
}

func (a *aliStorager) GetObject(bucket string, key string) ([]byte, error) {
	bucketObj, _ := a.client.Bucket(bucket)

	body, err := bucketObj.GetObject(key)
	if err != nil {
		logger.Error("oss get remote file failed, key: %s, err: %s", key, err)
		return nil, err
	}
	defer func() {
		_ = body.Close()
	}()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, body)
	if err != nil {
		logger.Error("oss copy body failed, key: %s, err: %s", key, err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *aliStorager) GetFile(bucket string, key string, localFile string) error {
	bucketObj, _ := a.client.Bucket(bucket)

	body, err := bucketObj.GetObject(key)
	if err != nil {
		logger.Error("oss get remote file failed, key: %s, err: %s", key, err)
		return err
	}
	defer func() {
		_ = body.Close()
	}()

	fd, err := os.OpenFile(localFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logger.Error("open local file failed, key: %s, filePath: %s err: %s",
			key, localFile, err)
		return err
	}
	defer func() {
		_ = fd.Close()
	}()

	_, err = io.Copy(fd, body)
	if err != nil {
		logger.Error("oss copy body failed, key: %s, err: %s", key, err)
		return err
	}
	return nil
}

func (a *aliStorager) PutObjectWithMeta(bucket string, key string, data []byte, metadata *Metadata) error {
	return a.PutReaderWithMeta(bucket, key, bytes.NewReader(data), metadata)
}

func (a *aliStorager) PutFileWithMeta(bucket string, key string, localFile string, metadata *Metadata) error {
	fi, err := os.Stat(localFile)
	if err != nil {
		return err
	}

	// 分流，大于100M使用分片上传，小于100M直接上传
	if fi.Size() > chunkSize {
		logger.Debug("file size over 100m, use multipart upload, ossPath: %s", key)
		// 向上取整
		chunkNum := (fi.Size() + chunkSize - 1) / chunkSize
		if err := a.putFileByMultipart(bucket, key, localFile, int(chunkNum), metadata); err != nil {
			return err
		}
		if err := a.SetObjectMeta(bucket, key, metadata); err != nil {
			return err
		}
		return nil
	} else {
		return a.putFileByDirect(bucket, key, localFile, metadata)
	}
}

func (a *aliStorager) putFileByDirect(bucket string, key string, localFile string, metadata *Metadata) error {
	bucketObj, err := a.client.Bucket(bucket)
	if err != nil {
		return err
	}

	// TODO 测试PutObjectFromFile
	options := buildOptions(metadata)
	return bucketObj.PutObjectFromFile(key, localFile, options...)
}

func (a *aliStorager) putFileByMultipart(bucket string, key string, localFile string,
	chunkCount int, metadata *Metadata) error {

	// 将本地文件分片，且分片数量指定为chunkCount
	chunks, err := oss.SplitFileByPartNum(localFile, chunkCount)
	if err != nil {
		return err
	}
	bucketObj, err := a.client.Bucket(bucket)
	if err != nil {
		return err
	}
	options := buildOptions(metadata)
	// TODO 怀疑options在这里设置无效，应该设置在后面的 bucketObj.CompleteMultipartUpload
	imr, err := bucketObj.InitiateMultipartUpload(key, options...)
	if err != nil {
		return err
	}
	fd, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = fd.Close()
	}()

	// 上传分片
	var parts []oss.UploadPart
	for _, chunk := range chunks {
		_, _ = fd.Seek(chunk.Offset, io.SeekStart)
		// 调用UploadPart方法上传每个分片。
		part, err := bucketObj.UploadPart(imr, fd, chunk.Size, chunk.Number)
		if err != nil {
			return err
		}
		parts = append(parts, part)
	}

	// 步骤3：完成分片上传。
	cmr, err := bucketObj.CompleteMultipartUpload(imr, parts)
	logger.Debug("upload file, result is %v", cmr)
	return err
}

func (a *aliStorager) PutReaderWithMeta(bucket string, key string, r io.Reader, metadata *Metadata) error {
	bucketObj, err := a.client.Bucket(bucket)
	if err != nil {
		return err
	}

	options := buildOptions(metadata)
	return bucketObj.PutObject(key, r, options...)
}

func (a *aliStorager) ListObjects(bucket string, prefix string) ([]FileMeta, error) {
	bucketObj, _ := a.client.Bucket(bucket)

	result := make([]FileMeta, 0, 32)

	prefixOption := oss.Prefix(prefix)
	maxKeys := oss.MaxKeys(1000)
	continueToken := ""
	for {

		lsRes, err := bucketObj.ListObjectsV2(prefixOption, maxKeys, oss.ContinuationToken(continueToken))
		if err != nil {
			return result, err
		}
		for _, val := range lsRes.Objects {
			content := FileMeta{Key: val.Key, Size: val.Size, ETag: val.ETag, LastModified: val.LastModified}
			result = append(result, content)
		}

		if lsRes.IsTruncated {
			continueToken = lsRes.NextContinuationToken
		} else {
			break
		}
	}
	return result, nil
}

func (a *aliStorager) DeleteObject(bucket string, key string) error {
	bucketObj, err := a.client.Bucket(bucket)

	err = bucketObj.DeleteObject(key)
	return err
}

func (a *aliStorager) CopyObject(bucket string, srcKey string, destKey string, metadata *Metadata) error {
	bucketObj, err := a.client.Bucket(bucket)
	_, err = bucketObj.CopyObject(srcKey, destKey)
	if err != nil {
		return err
	}
	if err = a.SetObjectMeta(bucket, destKey, metadata); err != nil {
		return err
	}

	return nil
}

func (a *aliStorager) SetObjectAcl(bucket string, key string, acl ACL) error {
	bucketObj, err := a.client.Bucket(bucket)
	err = bucketObj.SetObjectACL(key, oss.ACLType(acl))

	return err
}

func (a *aliStorager) SetObjectMeta(bucket string, key string, metadata *Metadata) error {
	bucketObj, err := a.client.Bucket(bucket)
	if err != nil {
		return err
	}

	options := buildOptions(metadata)
	return bucketObj.SetObjectMeta(key, options...)
}

func (a *aliStorager) IsObjectExist(bucket string, key string) (bool, error) {
	bucketObj, _ := a.client.Bucket(bucket)
	isExists, err := bucketObj.IsObjectExist(key)

	return isExists, err
}

type Policy struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}
type Statement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource string   `json:"Resource"`
}
type TokenInfo struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	ExpireTimeStamp int64 // 时间戳毫秒
}

func (a *aliStorager) AssumeRole(policy Policy, expires time.Duration) (*TokenInfo, error) {

	policyByte, err := json.Marshal(&policy)
	if err != nil {
		return nil, err
	}
	var assumeRoleFunc = func([]byte) (res *sts.AssumeRoleResponse, err error) {
		request := sts.AssumeRoleRequest{}
		request.SetPolicy(string(policyByte))
		request.SetRoleSessionName("test")
		request.SetRoleArn(a.config.RoleArn)
		request.SetDurationSeconds(int64(expires / time.Second))
		response, err := a.stsClient.AssumeRole(&request)
		return response, err
	}

	var response *sts.AssumeRoleResponse
	var aErr error
	var retryCount = 3
	for i := 1; i <= retryCount; i++ {
		response, aErr = assumeRoleFunc(policyByte)
		if aErr != nil {
			// 暂停i秒 ，再去尝试
			if i == retryCount {
				break
			}
			time.Sleep(time.Duration(i) * time.Second)
		} else {
			break
		}
	}

	if aErr != nil {
		return nil, aErr
	}

	expireTimeStamp := time.Now().Add(expires).UnixNano() / 1e6
	res := &TokenInfo{
		AccessKeyId:     *response.Body.Credentials.AccessKeyId,
		AccessKeySecret: *response.Body.Credentials.AccessKeySecret,
		SecurityToken:   *response.Body.Credentials.SecurityToken,
		ExpireTimeStamp: expireTimeStamp,
	}
	return res, nil

}

func (a *aliStorager) SignFile(bucket string, key string, expires time.Duration) (string, error) {
	// 获取存储空间。
	bucketObj, err := a.client.Bucket(bucket)
	if err != nil {
		return "", err
	}
	// 使用签名URL将OSS文件下载到流。
	signedUrl, err := bucketObj.SignURL(key, oss.HTTPGet, int64(expires/time.Second))
	if err != nil {
		return "", err
	}

	signedUrl = strings.Replace(signedUrl, "http", "https", 1)
	signedUrl = strings.Replace(signedUrl, a.config.EndpointInner, a.config.Endpoint, 1)
	return signedUrl, nil
}

const (
	ossActionPut    = "oss:PutObject"
	ossActionGet    = "oss:GetObject"
	ossActionPutAcl = "oss:PutObjectAcl"
)

func (a *aliStorager) GetDirToken(bucket string, remoteDir string, expires time.Duration) (*StsTokenInfo, error) {
	return a.createDirToken(bucket, remoteDir, expires, ossActionPut, ossActionGet)
}

func (a *aliStorager) GetDirTokenRead(bucket string, remoteDir string, expires time.Duration) (*StsTokenInfo, error) {
	return a.createDirToken(bucket, remoteDir, expires, ossActionGet)
}

func (a *aliStorager) createDirToken(bucket string, remoteDir string,
	expires time.Duration, actions ...string) (*StsTokenInfo, error) {

	resource := fmt.Sprintf("acs:oss:*:*:%s/%s*", bucket, remoteDir)
	policy := Policy{Version: "1", Statement: []Statement{{
		Effect:   "Allow",
		Action:   actions,
		Resource: resource,
	}}}
	ossToken, err := a.AssumeRole(policy, expires)
	if err != nil {
		return nil, err
	}
	res := &StsTokenInfo{
		Provider:        providerOss,
		AccessKeyID:     ossToken.AccessKeyId,
		AccessKeySecret: ossToken.AccessKeySecret,
		StsToken:        ossToken.SecurityToken,
		Bucket:          bucket,
		Region:          a.config.Region,
		Expire:          ossToken.ExpireTimeStamp,
		UploadPath:      remoteDir,
		Host:            a.config.Host,
		Endpoint:        a.config.Endpoint,
	}
	return res, nil
}

func (a *aliStorager) PresignObject(bucket string, key string, expires time.Duration) (string, error) {
	obj, _ := a.client.Bucket(bucket)
	signedUrl, err := obj.SignURL(key, oss.HTTPPut, int64(expires/time.Second))
	if err != nil {
		return "", err
	}

	signedUrl = strings.Replace(signedUrl, "http", "https", 1)
	signedUrl = strings.Replace(signedUrl, a.config.EndpointInner, a.config.Endpoint, 1)
	return signedUrl, nil
}

func (a *aliStorager) GetObjectMeta(bucket string, key string) (*FileMeta, error) {
	bucketObj, _ := a.client.Bucket(bucket)

	props, err := bucketObj.GetObjectDetailedMeta(key)
	if err != nil {
		return nil, err
	}

	res := &FileMeta{Key: key}
	res.ETag = props.Get("Etag")
	res.Size, _ = strconv.ParseInt(props.Get("Content-Length"), 10, 64)
	res.LastModified, err = time.Parse(http.TimeFormat, props.Get("Last-Modified"))

	// 填充 Metadata 信息
	res.Metadata = Metadata{
		ContentType:        props.Get("Content-Type"),
		ContentEncoding:    props.Get("Content-Encoding"),
		ContentDisposition: props.Get("Content-Disposition"),
	}

	// 获取对象的 ACL
	aclResult, err := bucketObj.GetObjectACL(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get object ACL: %v", err)
	}
	res.Metadata.Acl = ACL(aclResult.ACL)

	return res, nil
}

func buildOptions(metadata *Metadata) []oss.Option {
	if metadata == nil {
		return nil
	}

	var options []oss.Option
	if metadata.ContentType != "" {
		options = append(options, oss.ContentType(metadata.ContentType))
	}
	if metadata.ContentEncoding != "" {
		options = append(options, oss.ContentEncoding(metadata.ContentEncoding))
	}
	if metadata.ContentDisposition != "" {
		options = append(options, oss.ContentDisposition(metadata.ContentDisposition))
	}
	if metadata.Acl != "" {
		options = append(options, oss.ObjectACL(oss.ACLType(metadata.Acl)))
	}
	return options
}
