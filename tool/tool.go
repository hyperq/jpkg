package tool

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// ToMD5 MD5加密64位
func ToMD5(old string) (neww string) {
	neww = fmt.Sprintf("%x", md5.Sum([]byte(old)))
	return
}

// RandomStr 生成随机字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// GetNows 获取当前时间
func GetNows(formats ...interface{}) (timeres string) {
	timestemp := time.Now().Unix()
	tm := time.Unix(timestemp, 10)
	var format = `2006-01-02 15:04:05`
	if len(formats) == 1 {
		format = formats[0].(string)
	}
	return tm.Format(format)
}

// RandomNum 生成随机数
func RandomNum(length int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func TransSlice2Map(s []string) map[string]struct{} {
	var m = make(map[string]struct{})
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}

func Exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func HiddenPhone(phonenumber string) string {
	phonenumbers := []byte(phonenumber)
	phonenumbers[3] = '*'
	phonenumbers[4] = '*'
	phonenumbers[5] = '*'
	phonenumbers[6] = '*'
	return string(phonenumbers)
}
