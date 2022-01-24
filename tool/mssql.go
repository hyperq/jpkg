package tool

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	formatTime      = "15:04:05"
	formatDate      = "2006-01-02"
	formatDateTime  = "2006-01-02 15:04:05"
	formatDateTimeT = "2006-01-02T15:04:05"
)

// SetStructDefault post params
func SetStructDefault(obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if !isStructPtr(t) {
		return fmt.Errorf("%v must be  a struct pointer", obj)
	}
	t = t.Elem()
	v = v.Elem()
	return ParseDefultToStruct(t, v)
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func ParseDefultToStruct(objT reflect.Type, objV reflect.Value) error {
	for i := 0; i < objT.NumField(); i++ {
		fieldV := objV.Field(i)
		if !fieldV.CanSet() {
			continue
		}
		fieldT := objT.Field(i)
		if fieldT.Anonymous && fieldT.Type.Kind() == reflect.Struct {
			continue
		}
		if !fieldV.IsZero() {
			continue
		}
		var value string
		defaultValue := fieldT.Tag.Get("default")
		if defaultValue != "" {
			value = defaultValue
		} else {
			continue
		}

		switch fieldT.Type.Kind() {
		case reflect.Bool:
			if strings.ToLower(value) == "on" || strings.ToLower(value) == "1" || strings.ToLower(value) == "yes" {
				fieldV.SetBool(true)
				continue
			}
			if strings.ToLower(value) == "off" || strings.ToLower(value) == "0" || strings.ToLower(value) == "no" {
				fieldV.SetBool(false)
				continue
			}
			b, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			fieldV.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			fieldV.SetFloat(x)
		case reflect.Interface:
			fieldV.Set(reflect.ValueOf(value))
		case reflect.String:
			fieldV.SetString(value)
		case reflect.Struct:
			switch fieldT.Type.String() {
			case "time.Time":
				var (
					t   time.Time
					err error
				)
				if len(value) >= 25 {
					value = value[:25]
					t, err = time.ParseInLocation(time.RFC3339, value, time.Local)
				} else if strings.HasSuffix(strings.ToUpper(value), "Z") {
					t, err = time.ParseInLocation(time.RFC3339, value, time.Local)
				} else if len(value) >= 19 {
					if strings.Contains(value, "T") {
						value = value[:19]
						t, err = time.ParseInLocation(formatDateTimeT, value, time.Local)
					} else {
						value = value[:19]
						t, err = time.ParseInLocation(formatDateTime, value, time.Local)
					}
				} else if len(value) >= 10 {
					if len(value) > 10 {
						value = value[:10]
					}
					t, err = time.ParseInLocation(formatDate, value, time.Local)
				} else if len(value) >= 8 {
					if len(value) > 8 {
						value = value[:8]
					}
					t, err = time.ParseInLocation(formatTime, value, time.Local)
				}
				if err != nil {
					return err
				}
				fieldV.Set(reflect.ValueOf(t))
			}
		}
	}
	return nil
}
