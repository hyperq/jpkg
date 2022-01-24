package tool

import jsoniter "github.com/json-iterator/go"

func UnmarshalFromString(s string, o interface{}) (err error) {
	err = jsoniter.UnmarshalFromString(s, o)
	SetDCM(o)
	return
}
