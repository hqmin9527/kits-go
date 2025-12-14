package oss

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hqmin9527/kits-go/src/collection/_sync"
	"github.com/hqmin9527/kits-go/src/go_limit"
	"github.com/hqmin9527/kits-go/src/logger"
)

const goLimitCount = 20
const (
	downloadType = "download"
	previewType  = "preview"
)

// oss/obs的操作工具类，已经处理过了oss和obs的平台差异
type Wrapper struct {
	st storager
	oc *Config
}

func newOssWrapper(c *Config) (*Wrapper, error) {
	var st storager
	var err error
	switch c.Provider {
	case Aliyun:
		st, err = newAliStorager(c)
	case MinIo:
		st, err = newMinioStorager(c)
	case Huawei:
		st, err = newHwStorager(c)
	default:
		st, err = nil, errors.New("unknown oss provider")
	}

	if err != nil {
		return nil, err
	} else {
		return &Wrapper{oc: c, st: st}, nil
	}
}

func (o *Wrapper) GetBucketName() string {
	return o.oc.Bucket
}

// 如果key包含双斜杠，可能会导致文件无法下载和上传且不报错，所以这里预先检查打印日志
func checkKey(key string) {
	if strings.Contains(key, "//") {
		logger.Error("[OSS] key contains // please check, key: %s", key)
	}
}

func (o *Wrapper) GetObject(key string) ([]byte, error) {
	return o.st.GetObject(o.oc.Bucket, key)
}

func (o *Wrapper) GetFile(key string, localFile string) error {
	return o.st.GetFile(o.oc.Bucket, key, localFile)
}

func (o *Wrapper) PutObject(key string, data []byte, options ...Option) error {
	checkKey(key)
	return o.st.PutObjectWithMeta(o.oc.Bucket, key, data, buildMetadata(options))
}

func (o *Wrapper) PutFile(key string, filePath string, options ...Option) error {
	checkKey(key)
	return o.st.PutFileWithMeta(o.oc.Bucket, key, filePath, buildMetadata(options))
}

// 使用时要注意保护reader不要被其他协程关闭
func (o *Wrapper) PutReader(key string, r io.Reader, options ...Option) error {
	checkKey(key)
	return o.st.PutReaderWithMeta(o.oc.Bucket, key, r, buildMetadata(options))
}

// 上传文件夹, 返回上传失败的文件列表
func (o *Wrapper) PutFolder(remoteDir string, localFolder string, options ...Option) ([]string, error) {
	goLimit := go_limit.New(goLimitCount)
	var oneErr error // 某一个错误
	failedFiles := _sync.NewSlice[string]()
	localFolder = strings.ReplaceAll(localFolder, "\\", "/")
	_ = filepath.Walk(localFolder, func(localPath string, info os.FileInfo, err error) error {
		// 读取子目录没有权限
		if err != nil {
			failedFiles.Append(localPath)
			oneErr = err
			return nil
		}
		// Walk方法会递归遍历目录
		if info.IsDir() {
			return nil
		}

		goLimit.Run(func() {
			localPath = strings.ReplaceAll(localPath, "\\", "/")
			suffix := strings.Replace(localPath, localFolder, "", 1)
			remotePath := path.Join(remoteDir, suffix)
			if err := o.PutFile(remotePath, localPath, options...); err != nil {
				failedFiles.Append(localPath)
				oneErr = err
			}
		})

		return nil
	})
	goLimit.Wait()
	return failedFiles.List(), oneErr
}

// 获取文件概要信息：文件大小，最近修改时间
func (o *Wrapper) ListObjects(prefix string) ([]FileMeta, error) {
	return o.st.ListObjects(o.oc.Bucket, prefix)
}

func (o *Wrapper) DeleteObject(key string) error {
	checkKey(key)
	return o.st.DeleteObject(o.oc.Bucket, key)
}

func (o *Wrapper) DeleteFolder(remoteDir string) error {
	contents, err := o.ListObjects(remoteDir)
	if err != nil {
		return err
	}
	goLimit := go_limit.New(goLimitCount)
	for _, content := range contents {
		contentTmp := content
		goLimit.RunError(func() error {
			return o.DeleteObject(contentTmp.Key)
		})
	}
	goLimit.Wait()
	return goLimit.FirstError()
}

func (o *Wrapper) CopyObject(srcKey string, destKey string, options ...Option) error {
	checkKey(destKey)
	return o.st.CopyObject(o.oc.Bucket, srcKey, destKey, buildMetadata(options))
}

func (o *Wrapper) CopyFolder(remoteDir string, remoteDistDir string) error {
	contents, err := o.ListObjects(remoteDir)

	if err != nil {
		return err
	}
	goLimit := go_limit.New(goLimitCount)
	for _, content := range contents {
		contentTmp := content
		if !strings.HasSuffix(contentTmp.Key, "/") {
			var f = func() {
				subfix := strings.Replace(contentTmp.Key, remoteDir, "", 1)
				_ = o.CopyObject(contentTmp.Key, joinPath(remoteDistDir, subfix))
			}
			goLimit.Run(f)
		}
	}
	goLimit.Wait()
	return nil
}

func (o *Wrapper) SetObjectMeta(ossPath string, options ...Option) error {
	if len(options) == 0 {
		return nil
	}

	return o.st.SetObjectMeta(o.oc.Bucket, ossPath, buildMetadata(options))

}

func (o *Wrapper) SetFolderMeta(remoteDir string, options ...Option) error {
	if len(options) == 0 {
		return nil
	}

	contents, err := o.ListObjects(remoteDir)
	if err != nil {
		return err
	}
	goLimit := go_limit.New(goLimitCount)
	for _, content := range contents {
		contentTmp := content
		goLimit.Run(func() {
			_ = o.SetObjectMeta(contentTmp.Key, options...)
		})
	}

	goLimit.Wait()
	return nil
}

func (o *Wrapper) Move(srcKey string, destKey string) error {
	err := o.CopyObject(srcKey, destKey)
	if err != nil {
		return err
	}
	err = o.DeleteObject(srcKey)
	return err
}

func (o *Wrapper) MoveFolder(remoteDir string, remoteDistDir string) error {
	contents, err := o.ListObjects(remoteDir)
	if err != nil {
		return err
	}
	goLimit := go_limit.New(goLimitCount)
	for _, content := range contents {
		if !strings.HasSuffix(content.Key, "/") {
			contentTmp := content

			goLimit.Run(func() {
				subfix := strings.Replace(contentTmp.Key, remoteDir, "", 1)
				_ = o.CopyObject(contentTmp.Key, joinPath(remoteDistDir, subfix))
				_ = o.DeleteObject(contentTmp.Key)
			})
		}
	}
	goLimit.Wait()
	return nil
}

// 下载oss/obs中指定目录。 err!= nil时，返回失败的列表(nil为未获取到列表)。成功时列表为空
func (o *Wrapper) GetFolder(remoteDir string, localFolder string) ([]FileMeta, error) {
	match := func(key string) bool {
		return true
	}

	return o.GetFolderFilter(remoteDir, localFolder, match)
}

// 下载oss/obs中指定目录, pattern为使用key路径匹配的规则。 err!= nil时，返回失败的列表(nil为未获取到列表)。成功时列表为空
// 非匹配参考：https://www.cnblogs.com/asfeixue/p/lookahead.html
// 但go 不支持正则回溯如(?!，因此非匹配需要通过GetFolderFilter方法的match方法来实现。
func (o *Wrapper) GetFolderRegex(remoteDir string, localFolder string, pattern string) ([]FileMeta, error) {
	match := func(key string) bool {
		ok, _ := regexp.MatchString(pattern, key)
		return ok
	}
	return o.GetFolderFilter(remoteDir, localFolder, match)
}

func (o *Wrapper) GetFolderFilter(remoteDir string, localFolder string, match func(key string) bool) ([]FileMeta, error) {
	resultList, err := o.ListObjects(remoteDir)

	if err != nil {
		return nil, err
	}

	failChan := make(chan FileMeta, goLimitCount)
	var result error
	goLimit := go_limit.New(goLimitCount)
	for _, contentTmp := range resultList {
		content := contentTmp
		if !match(content.Key) {
			continue
		}
		subfix := strings.Replace(content.Key, remoteDir, "", 1)

		destFile := joinPath(localFolder, subfix)
		index := strings.LastIndex(destFile, "/")
		folder := destFile[0 : index+1]
		// 对于对象存储，/的文件目录是可选的，因此先每次都创建
		_ = os.MkdirAll(folder, 0755)
		if !strings.HasSuffix(content.Key, "/") {
			goLimit.Run(func() {
				err = o.GetFile(content.Key, joinPath(localFolder, subfix))
				if err != nil {
					failChan <- content
					result = err
					logger.Error("get storage file from %s to %s error: %s", content.Key, localFolder+subfix, err)
				}
			})

		}
	}
	goLimit.Wait()
	close(failChan)
	var failList []FileMeta
	for fm := range failChan {
		failList = append(failList, fm)
	}
	return failList, result
}

func (o *Wrapper) GetFolderSize(remoteDir string) (int64, error) {
	contents, err := o.ListObjects(remoteDir)
	if err != nil {
		return 0, err
	}
	var sum int64 = 0
	for _, content := range contents {
		sum += content.Size
	}
	return sum, nil
}

func (o *Wrapper) IsObjectExist(key string) (bool, error) {
	return o.st.IsObjectExist(o.oc.Bucket, key)
}

func (o *Wrapper) GetDirToken(remoteDir string, expires time.Duration) (*StsTokenInfo, error) {
	return o.st.GetDirToken(o.oc.Bucket, remoteDir, expires)
}

func (o *Wrapper) GetDirTokenRead(remoteDir string, expires time.Duration) (*StsTokenInfo, error) {
	return o.st.GetDirTokenRead(o.oc.Bucket, remoteDir, expires)
}

func (o *Wrapper) PresignObject(key string, expires time.Duration) (string, error) {
	return o.st.PresignObject(o.oc.Bucket, key, expires)
}

func (o *Wrapper) SignFile(key string, expires time.Duration) (string, error) {
	return o.st.SignFile(o.oc.Bucket, key, expires)
}

func (o *Wrapper) GetObjectMeta(key string) (*FileMeta, error) {
	return o.st.GetObjectMeta(o.oc.Bucket, key)
}

func (o *Wrapper) EscapeDownloadUrl(key string) string {
	return o.getUrlByType(key, downloadType)
}

func (o *Wrapper) EscapePreviewUrl(key string) string {
	return o.getUrlByType(key, previewType)
}
func (o *Wrapper) EscapeRawUrl(key string) string {
	if o.oc.Internal {
		return fmt.Sprintf("http://%s.%s/%s", o.oc.Bucket, o.oc.EndpointInner, key)
	}
	return o.getUrlByType(key, previewType)
}

func (o *Wrapper) GetConfig() *Config {
	return o.oc
}

func (o *Wrapper) getUrlByType(key, ty string) string {
	var domain string
	switch ty {
	case downloadType:
		domain = o.oc.Bucket + "." + o.oc.Endpoint
	case previewType:
		domain = o.oc.Host
	default:
		domain = o.oc.Host
	}
	ss := strings.Split(key, "/")
	for i, s := range ss {
		ss[i] = url.QueryEscape(s)
		if o.oc.Provider == MinIo {
			ss[i] = strings.ReplaceAll(ss[i], "+", "%20")
		}
	}
	encodedOssPath := strings.Join(ss, "/")
	switch o.oc.Provider {
	case MinIo:
		res := ""
		if o.oc.DownloadDomain != "" {
			res = fmt.Sprintf("https://%s/%s/%s", o.oc.DownloadDomain, o.oc.Bucket, encodedOssPath)
		} else {
			res = fmt.Sprintf("%s://%s/%s/%s", o.oc.Protocol, o.oc.Host, o.oc.Bucket, encodedOssPath)
		}
		return res
	default:
		return fmt.Sprintf("%s://%s/%s", o.oc.Protocol, domain, encodedOssPath)
	}
}

func (o *Wrapper) PutObjectWithMeta(key string, data []byte, md *Metadata) error {
	checkKey(key)
	return o.st.PutObjectWithMeta(o.oc.Bucket, key, data, md)
}

func joinPath(path1 string, path2 string) string {
	// 防止path1和path2之间的出现双斜杆
	if !strings.HasSuffix(path1, "/") {
		path1 = path1 + "/"
	}
	if strings.HasPrefix(path2, "/") {
		path2 = strings.TrimPrefix(path2, "/")
	}
	return path1 + path2
}
