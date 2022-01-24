package sf

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"github.com/hyperq/jpkg/tool/xhttp"
	"net/url"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"

	uuid "github.com/satori/go.uuid"

	"github.com/hyperq/jpkg/tool/xlog"
)

type CommonReq struct {
	PartnerID   string `json:"partnerID"`
	RequestID   string `json:"requestID"`
	ServiceCode string `json:"serviceCode"`
	Timestamp   string `json:"timestamp"`
	MsgDigest   string `json:"msgDigest"`
	MsgData     string `json:"msgData"`
}

type CommonResp struct {
	ApiErrorMsg   string `json:"apiErrorMsg"`
	ApiResponseID string `json:"apiResponseID"`
	ApiResultCode string `json:"apiResultCode"`
	ApiResultData string `json:"apiResultData"`
}

const requestOk = "A1000"

func (h *SFM) doProdPost(serviceCode, msgData string) (respdata string, err error) {
	httpClient := xhttp.NewClient()

	var req CommonReq
	req.PartnerID = h.PartnerID
	req.RequestID = uuid.NewV4().String()
	req.Timestamp = strconv.Itoa(int(time.Now().Unix() * 1000))
	req.MsgData = msgData
	req.ServiceCode = serviceCode
	req.MsgDigest = MsgDigest(req.MsgData + req.Timestamp + h.CheckWord)
	if h.Debug {
		xlog.Debugf("SF_SERVICE: %s", serviceCode)
		xlog.Debugf("SF_MSGDATA: %s", msgData)
	}
	res, bs, errs := httpClient.Type(xhttp.TypeUrlencoded).Post(h.ApiUrl).SendStruct(req).EndBytes()
	if len(errs) > 0 {
		err = errs[0]
		return
	}
	if h.Debug {
		xlog.Debugf("SF_Response: %d > %s", res.StatusCode, string(bs))
		xlog.Debugf("SF_Headers: %#v", res.Header)
	}
	var resp CommonResp
	if err = jsoniter.Unmarshal(bs, &resp); err != nil {
		return
	}
	if resp.ApiResultCode != requestOk {
		err = errors.New(resp.ApiErrorMsg)
		return
	}
	respdata = resp.ApiResultData
	return
}

func MsgDigest(msgDigests string) string {
	msgDigests = url.QueryEscape(msgDigests)
	m5 := md5.New()
	m5.Write([]byte(msgDigests))
	msgDigestsbyte := m5.Sum(nil)
	msgDigestsbase64 := base64.StdEncoding.EncodeToString(msgDigestsbyte)
	return msgDigestsbase64
}
