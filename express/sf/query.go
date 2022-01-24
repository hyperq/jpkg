package sf

import (
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type RouteReq struct { // | 必填 | 默认值 | 描述 |
	Language        string   `json:"language,omitempty"`        // | 否 | 0 | 返回描述语语言 0：中文 1：英文 2：繁体 |
	TrackingType    string   `json:"trackingType"`              // | 是 | 1 | 查询号类别: 1:根据顺丰运单号查询,trackingNumber将被当作顺丰运单号处理 2:根据客户订单号查询,trackingNumber将被当作客户订单号处理 |
	TrackingNumber  []string `json:"trackingNumber"`            // | 是 |   | 查询号: trackingType=1,则此值为顺丰运单号 如果trackingType=2,则此值为客户订单号 |
	MethodType      string   `json:"methodType,omitempty"`      // | 否 | 1 | 路由查询类别: 1:标准路由查询 2:定制路由查询 |
	ReferenceNumber string   `json:"referenceNumber,omitempty"` // | 否 |   | 参考编码(目前针对亚马逊客户,由客户传) |
	CheckPhoneNo    string   `json:"checkPhoneNo,omitempty"`    // | 否 |   | 电话号码验证 |
}

type Response struct {
	Success   bool         `json:"success"`
	ErrorCode string       `json:"errorCode"`
	ErrorMsg  interface{}  `json:"errorMsg"`
	MsgData   RouteMsgData `json:"msgData"`
}

type RouteMsgData struct {
	RouteResps []RouteResps `json:"routeResps"`
}

type RouteResps struct {
	MailNo string   `json:"mailNo"`
	Routes []Routes `json:"routes"`
}

type Routes struct {
	AcceptTime    string `json:"acceptTime"`
	AcceptAddress string `json:"acceptAddress"`
	Remark        string `json:"remark"`
	OpCode        string `json:"opCode"`
}

const queryRouteServiceCode = "EXP_RECE_SEARCH_ROUTES"

func (h *SFM) QueryRoute(req RouteReq) (resq Response, err error) {
	msgdata, _ := jsoniter.MarshalToString(req)
	respdata, err := h.doProdPost(queryRouteServiceCode, msgdata)
	if err != nil {
		return
	}
	if err = jsoniter.UnmarshalFromString(respdata, &resq); err != nil {
		return
	}
	if !resq.Success {
		err = errors.New(fmt.Sprint(resq.ErrorMsg))
	}
	return
}
