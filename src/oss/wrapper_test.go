package oss

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"testing"
	"time"
)

var ossHelper = genOssHelperAliyun()

// var ossHelper = genOssHelperMini()

// var ossHelper = genOssHelperHuawei()

var testDir = path.Join(ossHelper.oc.Root, ossHelper.oc.BaseDir, "test")

func TestGetBucketName(t *testing.T) {
	bucketName := ossHelper.GetBucketName()
	fmt.Println("bucketName:", bucketName)
}

func TestPutObjectDefaultPrivate(t *testing.T) {
	data := []byte("hello world")
	ossPath := path.Join(testDir, "hello.txt")
	err := ossHelper.PutObject(ossPath, data)
	fmt.Println("err:", err)
}

func TestDeleteObject(t *testing.T) {
	ossPath := path.Join(testDir, "hello.txt")
	err := ossHelper.DeleteObject(ossPath)
	fmt.Println("err:", err)
}

func TestPutObjectPublic(t *testing.T) {
	data := []byte("hello world")
	ossPath := path.Join(testDir, "hello.txt")
	err := ossHelper.PutObject(ossPath, data, AclPublicRead)
	fmt.Println("err:", err)
}

func TestPutObjectObjectAndFilename(t *testing.T) {
	data := []byte("hello world")
	ossPath := path.Join(testDir, "hello.txt")
	err := ossHelper.PutObject(ossPath, data, AclPublicRead, AttachFileName("哈.txt"))
	fmt.Println("err:", err)
}

func TestGetObject(t *testing.T) {
	ossPath := path.Join(testDir, "hello.txt")
	data, err := ossHelper.GetObject(ossPath)

	fmt.Println("err:", err)
	fmt.Println("data size", len(data))
}

func TestGetFile(t *testing.T) {
	ossPath := path.Join(testDir, "hello.txt")
	localPath := "./test_folder/hello.txt"
	err := ossHelper.GetFile(ossPath, localPath)
	fmt.Println("err:", err)
}

func TestPutFile(t *testing.T) {
	src := "test_folder/hello.txt"
	ossPath := path.Join(testDir, "hello.txt")
	err := ossHelper.PutFile(ossPath, src, AclPublicRead, AttachFileName("hello.txt"))
	fmt.Println("err:", err)
}

func TestPutFile2(t *testing.T) {
	src := "test_folder/7_Wow_Wow.svg"
	ossPath := path.Join(testDir, "7_Wow_Wow.svg")
	err := ossHelper.PutFile(ossPath, src, AclPublicRead)
	fmt.Println("err:", err)
}

func TestPutFolder(t *testing.T) {
	remoteDir := path.Join(testDir, "space")
	localDir := "./test_folder"
	failedFiles, err := ossHelper.PutFolder(remoteDir, localDir)
	fmt.Println("failedFiles:", failedFiles)
	fmt.Println("err:", err)
}

func TestListObjects(t *testing.T) {
	remoteDir := path.Join(testDir, "space")
	objects, err := ossHelper.ListObjects(remoteDir)
	fmt.Println("err:", err)
	fmt.Println("objects:", objects)
}

func TestGetFolder(t *testing.T) {
	remoteDir := path.Join(testDir, "space", "test_folder")
	localDir := "./test_folder"
	failList, err := ossHelper.GetFolder(remoteDir, localDir)
	fmt.Println("err:", err)
	fmt.Println("failList:", failList)
}

func TestGetFolderSize(t *testing.T) {
	remoteDir := path.Join(testDir, "space", "test_folder")
	size, err := ossHelper.GetFolderSize(remoteDir)
	fmt.Println("err:", err)
	fmt.Println("size:", size)
}

func TestCopyFolder(t *testing.T) {
	src := path.Join(testDir, "space", "test_folder")
	dst := path.Join(testDir, "space", "test_folder2")
	err := ossHelper.CopyFolder(src, dst)
	fmt.Println("err:", err)
}

func TestDeleteFolder(t *testing.T) {
	remoteDir := path.Join(testDir, "space")
	err := ossHelper.DeleteFolder(remoteDir)
	fmt.Println("err:", err)
}

func TestCopyObject(t *testing.T) {
	src := path.Join(testDir, "hello.txt")
	dst := path.Join(testDir, "^?测试.txt")
	// err := ossHelper.CopyObject(src, dst, AttachFileName("^?测试.txt"), AclPublicRead)
	err := ossHelper.CopyObject(src, dst)
	fmt.Println("err:", err)
}

func TestSetObjectMeta(t *testing.T) {
	ossPath := path.Join(testDir, "hello.txt")
	err := ossHelper.SetObjectMeta(ossPath, AclPrivate, AttachFileName("哈 哈 哈.txt"))
	fmt.Println("err:", err)
}

func TestIsObjectExist(t *testing.T) {
	ossPath := path.Join(testDir, "hello.txt")
	exist, err := ossHelper.IsObjectExist(ossPath)
	fmt.Println("err:", err)
	fmt.Println("exist:", exist)
}

func TestIsObjectExist2(t *testing.T) {
	ossPath := path.Join(testDir, "hello.txt111111")
	exist, err := ossHelper.IsObjectExist(ossPath)
	fmt.Println("err:", err)
	fmt.Println("exist:", exist)
}

func TestGetDirToken(t *testing.T) {
	remoteDir := path.Join(testDir)
	token, err := ossHelper.GetDirToken(remoteDir, time.Hour)
	fmt.Println("err:", err)
	d, _ := json.Marshal(token)
	fmt.Println("token:", string(d))
}

func TestGetDirTokenRead(t *testing.T) {
	remoteDir := path.Join(testDir)
	token, err := ossHelper.GetDirTokenRead(remoteDir, time.Hour)
	fmt.Println("err:", err)
	d, _ := json.Marshal(token)
	fmt.Println("token:", string(d))
}

func TestPresignObject(t *testing.T) {
	ossPath := path.Join(testDir, "hello.txt")
	url, err := ossHelper.PresignObject(ossPath, time.Hour)
	if err != nil {
		t.Fatalf("sign url failed, err: %s", err)
	}
	t.Logf("sign url success, url: %s", url)

	payload := strings.NewReader("你好")
	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		t.Fatalf("build request failed, err: %s", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request failed, err: %s", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body failed, err: %s", err)
	}
	t.Logf("put success, response body: %s", string(body))
}

func TestSignFile(t *testing.T) {
	ossPath := path.Join(testDir, "hello.txt")
	sign, err := ossHelper.SignFile(ossPath, time.Hour)
	fmt.Println("err:", err)
	fmt.Println("sign:", sign)
}

func TestGetObjectMeta(t *testing.T) {
	ossPath := path.Join(testDir, "hello111.txt")
	meta, err := ossHelper.GetObjectMeta(ossPath)
	fmt.Println("err:", err)
	fmt.Println("meta:", meta)
}

func TestSetFolderMeta(t *testing.T) {
	remoteDir := path.Join(testDir)
	err := ossHelper.SetFolderMeta(remoteDir, AclPrivate)
	fmt.Println("err:", err)
}

func TestPutMultiFile(t *testing.T) {
	// src := "E:\\test\\old\\2G.mp4"
	src := "/Users/hqmin/Desktop/测试文件/音视频/200M.mp4"
	dst := path.Join(testDir, "test.mp4")
	err := ossHelper.PutFile(dst, src, AclPublicRead, AttachFileName("哈.mp4"))
	fmt.Println("err:", err)
}

func TestEscapeRawUrl(t *testing.T) {
	ossPath := path.Join(testDir, "0920.mp4")
	url := ossHelper.EscapeRawUrl(ossPath)
	fmt.Println("url:", url)
}
