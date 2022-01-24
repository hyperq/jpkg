package ali

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/hyperq/jpkg/conf"
)

var smsclient *dysmsapi.Client

func SmsSend(req *dysmsapi.SendSmsRequest) (err error) {
	_, err = smsclient.SendSms(req)
	return
}

func smsInit() {
	sms, ok := conf.Config.Oss["sms"]
	if ok {
		var err error
		smsclient, err = dysmsapi.NewClientWithAccessKey("cn-hangzhou", sms.AliyunAccessKeyID, sms.AliyunAccessKeySecret)
		if err != nil {
			panic(err)
		}
	}

}
