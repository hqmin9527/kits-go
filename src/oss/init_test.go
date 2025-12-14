package oss

import "fmt"

func genOssHelperAliyun() *Wrapper {
	config := &Config{
		Provider:        Aliyun,
		AccessKeyId:     AliAccessKeyId,
		AccessKeySecret: AliAccessKeySecret,
		RoleArn:         "acs:ram::1845455964094256:role/oss-readonly-role",
		Bucket:          "file-plaso",
		Region:          "oss-cn-hangzhou",
		Endpoint:        "oss-cn-hangzhou.aliyuncs.com",
		EndpointInner:   "oss-cn-hangzhou-internal.aliyuncs.com",
		StsEndpoint:     "sts.cn-hangzhou.aliyuncs.com",
		Internal:        false,
		Host:            "file.plaso.cn",
		Root:            "dev-plaso",
		BaseDir:         "infinite_wb/",
		Protocol:        "https",
	}
	o, err := newOssWrapper(config)
	if err != nil {
		fmt.Println("err:", err)
	}
	return o
}

func genOssHelperMini() *Wrapper {
	config := &Config{
		Provider:        MinIo,
		AccessKeyId:     MinioAccessKeyId,
		AccessKeySecret: MinioAccessKeySecret,
		RoleArn:         "",

		Bucket:        "file-plaso",
		Region:        "oss-cn-hangzhou",
		Endpoint:      "120.27.221.27:9000",
		EndpointInner: "",
		StsEndpoint:   "",
		Internal:      false,
		Host:          "",
		// DownloadDomain: "efp.sem.tsinghua.edu.cn/bosuo/minio",
	}
	o, err := newOssWrapper(config)
	if err != nil {
		fmt.Println("err:", err)
	}
	return o
}

func genOssHelperHuawei() *Wrapper {
	config := &Config{
		Provider:        Huawei,
		AccessKeyId:     HuaweiAccessKeyId,
		AccessKeySecret: HuaweiAccessKeySecret,
		RoleArn:         "",
		Bucket:          "plaso-school",
		Region:          "",
		Endpoint:        "https://obs.cidc-rp-12.joint.cmecloud.cn",
		EndpointInner:   "",
		StsEndpoint:     "iam.cidc-rp-12.joint.cmecloud.cn",
		Internal:        false,
		Host:            "",
		Protocol:        "",
		DownloadDomain:  "",
		Root:            "dev",
		BaseDir:         "infinite_wb/",
	}
	o, err := newOssWrapper(config)
	if err != nil {
		fmt.Println("err:", err)
	}
	return o
}
