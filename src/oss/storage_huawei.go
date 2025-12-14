package oss

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hqmin9527/kits-go/src/logger"
	"github.com/hqmin9527/kits-go/src/utils"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/pkg/errors"
)

// SDK接口概览
// https://support.huaweicloud.com/sdk-go-devg-obs/obs_23_0002.html

const (
	obsPutAction = "obs:object:PutObject"
	obsGetAction = "obs:object:GetObject"
)

type hwStorager struct {
	config    *Config
	client    *obs.ObsClient
	iamClient *iam.IamClient
}

func newHwStorager(c *Config) (*hwStorager, error) {
	endpoint := utils.If(c.Internal, c.EndpointInner, c.Endpoint)
	logger.Info("create huawei obs client")
	client, err := obs.New(c.AccessKeyId, c.AccessKeySecret, endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "create client")
	}
	auth, err := global.NewCredentialsBuilder().WithAk(c.AccessKeyId).WithSk(c.AccessKeySecret).SafeBuild()
	if err != nil {
		return nil, errors.Wrap(err, "create iamClient auth")
	}
	builder, err := iam.IamClientBuilder().WithEndpoints([]string{c.StsEndpoint}).WithCredential(auth).SafeBuild()
	if err != nil {
		return nil, errors.Wrap(err, "create iamClient builder")
	}
	iamClient := iam.NewIamClient(builder)
	return &hwStorager{config: c, client: client, iamClient: iamClient}, nil
}

func (h *hwStorager) GetObject(bucket string, key string) ([]byte, error) {
	input := new(obs.GetObjectInput)
	input.Bucket = bucket
	input.Key = key
	output, err := h.client.GetObject(input)
	if err != nil {
		logger.Error("obs get remote file failed, key: %s err: %s", key, err)
		return nil, err
	}
	defer func() {
		_ = output.Body.Close()
	}()

	length := output.ContentLength
	if length > 1000*1000 {
		logger.Warn("obs file is large then 1M ,you should download file then process. %s/%s", bucket, key)
	}

	buf := new(bytes.Buffer)
	_, readErr := io.Copy(buf, output.Body)
	if readErr != nil {
		logger.Error("obs get remote file %s success ,but read fail %s", key, readErr)
		return nil, readErr
	}
	return buf.Bytes(), nil
}

func (h *hwStorager) GetFile(bucket string, key string, localFile string) error {
	input := new(obs.GetObjectInput)
	input.Bucket = bucket
	input.Key = key
	output, err := h.client.GetObject(input)
	if err != nil {
		logger.Error("obs get remote file failed, key: %s, err: %s", key, err)
		return err
	}
	defer func() {
		_ = output.Body.Close()
	}()

	fd, err := os.OpenFile(localFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logger.Error("obs open local file failed, key: %s, filePath: %s, err: %s",
			key, localFile, err)
		return err
	}
	defer func() {
		_ = fd.Close()
	}()

	_, err = io.Copy(fd, output.Body)
	if err != nil {
		logger.Error("obs write to local file failed, key: %s, filePath: %s, err: %s",
			key, localFile, err)
	}
	return nil
}

func (h *hwStorager) PutObjectWithMeta(bucket string, key string, data []byte, metadata *Metadata) error {
	return h.PutReaderWithMeta(bucket, key, bytes.NewReader(data), metadata)
}

func (h *hwStorager) PutFileWithMeta(bucket string, key string, localFile string, metadata *Metadata) error {
	fi, err := os.Stat(localFile)
	if err != nil {
		return err
	}

	// 分流，大于100M使用分片上传，小于100M直接上传
	if fi.Size() > chunkSize {
		logger.Debug("obs file size over 100m, use multipart upload, ossPath: %s", key)
		return h.putFileByMultipart(bucket, key, localFile, metadata)
	} else {
		return h.putFileByDirect(bucket, key, localFile, metadata)
	}
}

func (h *hwStorager) putFileByDirect(bucket string, key string, localFile string, metadata *Metadata) error {
	input := new(obs.PutFileInput)
	input.Bucket = bucket
	input.Key = key
	input.SourceFile = localFile
	setObsInput(&input.ACL, &input.HttpHeader, metadata)

	_, err := h.client.PutFile(input)
	if err != nil {
		logger.Error("obs putFileByDirect failed, key: %s, err: %s", key, err)
		return err
	}
	return nil
}

func (h *hwStorager) putFileByMultipart(bucket string, key string, localFile string, metadata *Metadata) error {
	input := new(obs.UploadFileInput)
	input.Bucket = bucket
	input.Key = key
	input.UploadFile = localFile
	input.EnableCheckpoint = true // 开启断点续传模式
	input.PartSize = chunkSize    // 指定分片大小100MB
	input.TaskNum = 5             // 指定分片上传时最大并发数

	_, err := h.client.UploadFile(input)
	if err != nil {
		logger.Error("obs putFileByMultipart failed, key: %s, err: %s", key, err)
		return err
	}
	return h.SetObjectMeta(bucket, key, metadata)
}

func (h *hwStorager) PutReaderWithMeta(bucket string, key string, r io.Reader, metadata *Metadata) error {
	input := new(obs.PutObjectInput)
	input.Bucket = bucket
	input.Key = key
	input.Body = r
	if metadata != nil && metadata.Acl != "" {
		input.ACL = obs.AclType(metadata.Acl)
	}
	setObsInput(&input.ACL, &input.HttpHeader, metadata)

	_, err := h.client.PutObject(input)
	if err != nil {
		logger.Error("obs PutReaderWithMeta failed, key: %s, err: %s", key, err)
		return err
	}
	return nil
}

func (h *hwStorager) ListObjects(bucket string, prefix string) ([]FileMeta, error) {
	input := new(obs.ListObjectsInput)
	input.Bucket = bucket
	input.Prefix = prefix
	input.MaxKeys = 1000 // 一次调用ListObjects取1000条，相当于分页取

	result := make([]FileMeta, 0, 32)
	for {
		output, err := h.client.ListObjects(input)
		if err != nil {
			return result, err
		}
		for _, val := range output.Contents {
			content := FileMeta{Key: val.Key, Size: val.Size, ETag: val.ETag, LastModified: val.LastModified}
			result = append(result, content)
		}

		if output.IsTruncated {
			input.Marker = output.NextMarker
		} else {
			break
		}

	}
	return result, nil
}

func (h *hwStorager) DeleteObject(bucket string, key string) error {
	input := new(obs.DeleteObjectInput)
	input.Bucket = bucket
	input.Key = key
	_, err := h.client.DeleteObject(input)
	return err
}

func (h *hwStorager) CopyObject(bucket string, srcKey string, destKey string, metadata *Metadata) error {
	input := new(obs.CopyObjectInput)
	input.Bucket = bucket
	input.Key = destKey
	input.CopySourceBucket = bucket
	input.CopySourceKey = srcKey

	_, err := h.client.CopyObject(input)
	if err != nil {
		return err
	}

	// CopyObject时，目标文件的acl默认是私有的
	// 没有指定目标文件acl时，尝试获取源文件acl
	if metadata != nil && !metadata.HasAcl() {
		acl, _ := h.getObjectAcl(bucket, srcKey)
		if acl != "" {
			metadata.Acl = acl
		}
	}

	// copy设置header不成功，再执行一次SetObjectMeta
	return h.SetObjectMeta(bucket, destKey, metadata)
}

func (h *hwStorager) getObjectAcl(bucket string, key string) (ACL, error) {
	input := new(obs.GetObjectAclInput)
	input.Bucket = bucket
	input.Key = key

	output, err := h.client.GetObjectAcl(input)
	if err != nil {
		logger.Error("obs get remote object acl failed, key: %s, err: %s", key, err)
		return "", err
	}
	var acl ACL
	for _, grant := range output.Grants {
		if grant.Grantee.ID == "" {
			acl = mapObsAcl(grant.Permission)
		}
	}
	return acl, nil
}

func mapObsAcl(perm obs.PermissionType) ACL {
	switch perm {
	case obs.PermissionRead:
		return ACL_PUBLIC_READ
	default:
		return ""
	}
}

func (h *hwStorager) SetObjectMeta(bucket string, key string, metadata *Metadata) error {
	if metadata == nil {
		return nil
	}
	if metadata.HasAcl() {
		if err := h.setObjectAcl(bucket, key, metadata); err != nil {
			return err
		}
	}
	if metadata.HasHeader() {
		if err := h.setObjectHeader(bucket, key, metadata); err != nil {
			return err
		}
	}
	return nil
}

func (h *hwStorager) setObjectAcl(bucket string, key string, metadata *Metadata) error {
	input := new(obs.SetObjectAclInput)
	input.Bucket = bucket
	input.Key = key
	setObsInput(&input.ACL, nil, metadata)

	_, err := h.client.SetObjectAcl(input)
	return err
}

func (h *hwStorager) setObjectHeader(bucket string, key string, metadata *Metadata) error {
	input := new(obs.SetObjectMetadataInput)
	input.Bucket = bucket
	input.Key = key
	setObsInput(nil, &input.HttpHeader, metadata)

	_, err := h.client.SetObjectMetadata(input)
	return err
}

func (h *hwStorager) IsObjectExist(bucket string, key string) (bool, error) {
	input := new(obs.GetObjectMetadataInput)
	input.Bucket = bucket
	input.Key = key

	output, err := h.client.GetObjectMetadata(input)
	if err != nil {
		if isObsNotFoundError(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return output != nil, nil
}

func isObsNotFoundError(err error) bool {
	if obsError, ok := errors.Cause(err).(obs.ObsError); ok {
		return obsError.StatusCode == 404
	}
	return false
}

func (h *hwStorager) GetDirToken(bucket string, remoteDir string, expires time.Duration) (*StsTokenInfo, error) {
	return h.createDirToken(bucket, remoteDir, expires, obsPutAction, obsGetAction)
}

func (h *hwStorager) GetDirTokenRead(bucket string, remoteDir string, expires time.Duration) (*StsTokenInfo, error) {
	return h.createDirToken(bucket, remoteDir, expires, obsGetAction)
}

// expires: 15min ~ 24h，默认15min
func (h *hwStorager) createDirToken(bucket string, remoteDir string,
	expires time.Duration, obsActions ...string) (*StsTokenInfo, error) {

	identify := new(model.TokenAuthIdentity)
	identify.Methods = []model.TokenAuthIdentityMethods{model.GetTokenAuthIdentityMethodsEnum().TOKEN}
	identify.Token = new(model.IdentityToken)
	identify.Token.DurationSeconds = utils.Ref(int32(expires.Seconds()))
	policy := new(model.ServicePolicy)
	policy.Version = "1.1"
	resource := []string{fmt.Sprintf("obs:*:*:object:%s/%s*", bucket, remoteDir)}
	policy.Statement = []model.ServiceStatement{{
		Effect:    model.GetServiceStatementEffectEnum().ALLOW,
		Action:    obsActions,
		Resource:  &resource,
		Condition: nil,
	}}
	auth := new(model.TokenAuth)
	auth.Identity = identify
	body := new(model.CreateTemporaryAccessKeyByTokenRequestBody)
	body.Auth = auth
	req := new(model.CreateTemporaryAccessKeyByTokenRequest)
	req.Body = body

	resp, err := h.iamClient.CreateTemporaryAccessKeyByToken(req)
	if err != nil {
		logger.Error("obs create dir token failed, err: %s", err)
		return nil, err
	}
	token := &StsTokenInfo{
		AccessKeyID:     resp.Credential.Access,
		AccessKeySecret: resp.Credential.Secret,
		StsToken:        resp.Credential.Securitytoken,
		Bucket:          bucket,
		Region:          h.config.Region,
		Provider:        providerObs,
		Expire:          time.Now().Add(expires).UnixMilli(),
		UploadPath:      remoteDir,
		Host:            h.config.Host,
		Endpoint:        h.config.Endpoint,
	}
	return token, nil
}

func (h *hwStorager) PresignObject(bucket string, key string, expires time.Duration) (string, error) {
	input := new(obs.CreateSignedUrlInput)
	input.Bucket = bucket
	input.Key = key
	input.Expires = int(expires / time.Second)
	input.Method = obs.HttpMethodPut

	output, err := h.client.CreateSignedUrl(input)
	if err != nil {
		logger.Error("obs createSignedUrl[Put] failed, key: %s, err: %s", key, err)
		return "", err
	}

	return output.SignedUrl, nil
}

func (h *hwStorager) SignFile(bucket string, key string, expires time.Duration) (string, error) {
	input := new(obs.CreateSignedUrlInput)
	input.Bucket = bucket
	input.Key = key
	input.Expires = int(expires / time.Second)
	input.Method = obs.HttpMethodGet

	output, err := h.client.CreateSignedUrl(input)
	if err != nil {
		logger.Error("obs createSignedUrl[Get] failed, key: %s, err: %s", key, err)
		return "", err
	}
	return output.SignedUrl, nil
}

func (h *hwStorager) GetObjectMeta(bucket string, key string) (*FileMeta, error) {
	input := new(obs.GetObjectMetadataInput)
	input.Bucket = bucket
	input.Key = key

	output, err := h.client.GetObjectMetadata(input)
	if err != nil {
		return nil, err
	}
	res := &FileMeta{
		Key:          key,
		Size:         output.ContentLength,
		ETag:         output.ETag,
		LastModified: output.LastModified,
	}

	acl, err := h.getObjectAcl(bucket, key)
	if err != nil {
		return nil, err
	}
	res.Metadata = Metadata{
		ContentDisposition: output.ContentDisposition,
		ContentType:        output.ContentType,
		ContentEncoding:    output.ContentEncoding,
		Acl:                acl,
	}
	return res, nil
}

func setObsInput(acl *obs.AclType, header *obs.HttpHeader, metadata *Metadata) {
	if metadata == nil {
		return
	}
	// 设置acl
	if acl != nil {
		if metadata.Acl != "" {
			*acl = obs.AclType(metadata.Acl)
		}
	}
	// 设置header
	if header != nil {
		if metadata.ContentType != "" {
			header.ContentType = metadata.ContentType
		}
		if metadata.ContentEncoding != "" {
			header.ContentEncoding = metadata.ContentEncoding
		}
		if metadata.ContentDisposition != "" {
			header.ContentDisposition = metadata.ContentDisposition
		}
	}
}
