package ali

import (
	"fmt"
	"github.com/hyperq/jpkg/conf"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var bucket *oss.Bucket
var bucketbak *oss.Bucket

func ossInit() {
	img, ok := conf.Config.Oss["img"]
	if ok {
		client, err := oss.New(img.AliyunEndPoint, img.AliyunAccessKeyID, img.AliyunAccessKeySecret)
		if err != nil {
			fmt.Println(err)
		}
		bucket, err = client.Bucket(img.AliyunBucket)
		if err != nil {
			fmt.Println(err)
		}
	}
	bak, ok := conf.Config.Oss["bak"]
	if ok {
		clientbak, err := oss.New(bak.AliyunEndPoint, bak.AliyunAccessKeyID, bak.AliyunAccessKeySecret)
		if err != nil {
			fmt.Println(err)

		}
		bucketbak, err = clientbak.Bucket(bak.AliyunBucket)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// OssUpfile 上传图片
func OssUpfile(osspath, filepath string) (rtpath string, err error) {
	img, ok := conf.Config.Oss["img"]
	if ok && img.AliyunBucket != "" {
		err = bucket.PutObjectFromFile(osspath, filepath)
		if err != nil {
			return
		}
		rtpath = img.AliyunBucketDomin + "/" + osspath
		if strings.Contains(filepath, "png") || strings.Contains(filepath, "jpg") || strings.Contains(filepath, "jpeg") {
			rtpath += "?x-ali-process=style/" + img.AliyunBucketStyle
		}
	}
	return
}

// OssUpfileBak 上传备份文件
func OssUpfileBak(osspath, filepath string) (rtpath string, err error) {
	bak, ok := conf.Config.Oss["img"]
	if bak.AliyunBucket != "" && ok {
		err = bucketbak.PutObjectFromFile(osspath, filepath)
		if err != nil {
			return
		}
		rtpath = bak.AliyunBucketDomin + "/" + osspath
	}
	return
}
