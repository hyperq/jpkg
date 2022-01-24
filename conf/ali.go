package conf

type alipay struct {
	AppID  string
	Secret string
}

type oss struct {
	AliyunEndPoint        string
	AliyunBucket          string
	AliyunBucketDomin     string
	AliyunBucketStyle     string
	AliyunAccessKeyID     string
	AliyunAccessKeySecret string
}
