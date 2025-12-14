package oss

import (
	"encoding/json"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/hqmin9527/kits-go/src/utils"
)

// provider 配置选项
const (
	Aliyun = "ali"
	MinIo  = "minio"
	Huawei = "hw"
)

// 给客户端provider
const (
	providerOss   = "OSS"
	providerObs   = "OBS"
	providerMinio = "MINIO"
)

const chunkSize = 100 * 1024 * 1024 // 100MB

type Config struct {
	Provider        string `json:"provider"`        // ali, minio, hw
	AccessKeyId     string `json:"accessKeyId"`     // accessKeyId
	AccessKeySecret string `json:"accessKeySecret"` // accessKeySecret
	RoleArn         string `json:"roleArn"`         // 给客户端临时授权的角色
	Bucket          string `json:"bucket"`          // 存储空间
	Region          string `json:"region"`          // 地域
	Endpoint        string `json:"endpoint"`        // 访问域名
	EndpointInner   string `json:"endpointInner"`   // 内网访问域名
	StsEndpoint     string `json:"stsEndpoint"`     // 临时授权的访问域名
	Internal        bool   `json:"internal"`        // 是否走内网
	Host            string `json:"host"`            // 加速域名（如果没有设置为Bucket.Endpoint）
	Protocol        string `json:"protocol"`        // 协议
	DownloadDomain  string `json:"downloadDomain"`  // 下载域名
	Root            string `json:"root"`            // 正式文件的根目录
	TmpRoot         string `json:"tmpRoot"`         // 临时文件的根目录
	BaseDir         string `json:"baseDir"`         // u3服务文件都在该目录下
	Prefix          string `json:"prefix"`          // 共同前缀
}

func (c *Config) MarshalJSON() ([]byte, error) {
	type alias Config
	return json.Marshal(&struct {
		*alias
		AccessKeyId     string `json:"accessKeyId"`
		AccessKeySecret string `json:"accessKeySecret"`
	}{
		alias:           (*alias)(c),
		AccessKeyId:     utils.MaskString(c.AccessKeyId),
		AccessKeySecret: utils.MaskString(c.AccessKeySecret),
	})
}

func (c *Config) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

// 返回给客户端临时token信息
type StsTokenInfo struct {
	Provider        string `json:"provider"` // OSS, MINIO, HW
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	StsToken        string `json:"stsToken"`
	Bucket          string `json:"bucket"`
	Region          string `json:"region"`
	Expire          int64  `json:"expire"` // 时间戳（毫秒）
	UploadPath      string `json:"uploadPath"`
	Host            string `json:"host"`
	Endpoint        string `json:"endpoint"`
}

func (t *StsTokenInfo) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

type FileMeta struct {
	Key          string
	Size         int64
	ETag         string
	LastModified time.Time
	Metadata
}

type ACL string

const (
	ACL_PRIVATE           ACL = "private"
	ACL_PUBLIC_READ       ACL = "public-read"
	ACL_PUBLIC_READ_WRITE ACL = "public-read-write" // 2021.1.8当前测试华为public-read-write设置未生效
)

type Metadata struct {
	ContentType        string
	ContentEncoding    string
	ContentDisposition string
	Acl                ACL
}

func (m *Metadata) HasAcl() bool {
	return m.Acl != ""
}

func (m *Metadata) HasHeader() bool {
	return m.ContentType != "" || m.ContentEncoding != "" || m.ContentDisposition != ""
}

func buildMetadata(ops []Option) *Metadata {
	m := &Metadata{}
	for _, op := range ops {
		if op != nil {
			op(m)
		}
	}
	return m
}

type Option func(m *Metadata)

func setAcl(acl ACL) Option {
	return func(m *Metadata) {
		m.Acl = acl
	}
}

var AclPublicRead = setAcl(ACL_PUBLIC_READ)
var AclPrivate = setAcl(ACL_PRIVATE)

var AttachFileName = func(fileName string) Option {
	return func(m *Metadata) {
		fileName = url.PathEscape(fileName)
		// oss会忽略开头的“.”，如果只有扩展名，添加文件名为"新文件"
		ext := strings.TrimLeft(filepath.Ext(fileName), ".")
		if ext != "" && strings.TrimLeft(fileName, ".") == ext {
			fileName = url.PathEscape("新文件.") + ext
		}
		m.ContentDisposition = "attachment;" + "filename=\"" + fileName + "\";" +
			"filename*=utf-8''" + fileName
	}
}
