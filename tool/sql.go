package tool

import (
	"reflect"
	"strings"
	"time"
)

func SetDCM(obj interface{}) {
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()
	ParseDefultToStruct(t, v)
	var iszero bool
	for i := 0; i < t.NumField(); i++ {
		fieldV := v.Field(i)
		fieldT := t.Field(i)
		tags := fieldT.Tag.Get("gorm")
		if strings.Index(tags, ";pk") > -1 {
			iszero = fieldV.IsZero()
			break
		}
	}
	now := time.Now()
	var index = 0
	for i := 0; i < t.NumField(); i++ {
		fieldV := v.Field(i)
		fieldT := t.Field(i)
		tags := fieldT.Tag.Get("gorm")
		if tags == "is_delete" && iszero {
			fieldV.SetInt(0)
			index++
		}
		if tags == "create_time" && iszero {
			fieldV.Set(reflect.ValueOf(now))
			index++
		}
		if tags == "modify_time" {
			fieldV.Set(reflect.ValueOf(now))
			index++
		}
		if index > 2 {
			break
		}
	}
}
