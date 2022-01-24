package bmap

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"
)

const (
	NULL = ""
)

type BodyMap map[string]interface{}

var mu = new(sync.RWMutex)

// 设置参数
func (bm BodyMap) Set(key string, value interface{}) BodyMap {
	mu.Lock()
	bm[key] = value
	mu.Unlock()
	return bm
}

func (bm BodyMap) SetBodyMap(key string, value func(bm BodyMap)) BodyMap {
	_bm := make(BodyMap)
	value(_bm)

	mu.Lock()
	bm[key] = _bm
	mu.Unlock()
	return bm
}

// 获取参数，同 GetString()
func (bm BodyMap) Get(key string) string {
	return bm.GetString(key)
}

// 获取参数转换string
func (bm BodyMap) GetString(key string) string {
	if bm == nil {
		return NULL
	}
	mu.RLock()
	defer mu.RUnlock()
	value, ok := bm[key]
	if !ok {
		return NULL
	}
	v, ok := value.(string)
	if !ok {
		return convertToString(value)
	}
	return v
}

// 获取原始参数
func (bm BodyMap) GetInterface(key string) interface{} {
	if bm == nil {
		return nil
	}
	mu.RLock()
	defer mu.RUnlock()
	return bm[key]
}

// 删除参数
func (bm BodyMap) Remove(key string) {
	mu.Lock()
	delete(bm, key)
	mu.Unlock()
}

// 置空BodyMap
func (bm BodyMap) Reset() {
	mu.Lock()
	for k := range bm {
		delete(bm, k)
	}
	mu.Unlock()
}

func (bm BodyMap) JsonBody() (jb string) {
	mu.Lock()
	defer mu.Unlock()
	bs, err := json.Marshal(bm)
	if err != nil {
		return ""
	}
	jb = string(bs)
	return jb
}

// ("bar=baz&foo=quux") sorted by key.
func (bm BodyMap) EncodeSignParams(apiKey string) string {
	var (
		buf     strings.Builder
		keyList []string
	)
	for k := range bm {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	buf.WriteString(apiKey)
	for _, k := range keyList {
		if v := bm.GetString(k); v != NULL {
			buf.WriteString(k)
			buf.WriteString(v)
		}
	}
	buf.WriteString(apiKey)
	return strings.ToUpper(ToMD5(buf.String()))
}

// ("bar=baz&foo=quux") sorted by key.
func (bm BodyMap) EncodeURLParams() string {
	var (
		buf  strings.Builder
		keys []string
	)
	for k := range bm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if v := bm.GetString(k); v != NULL {
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
			buf.WriteByte('&')
		}
	}
	if buf.Len() <= 0 {
		return NULL
	}
	return buf.String()[:buf.Len()-1]
}

func (bm BodyMap) CheckEmptyError(keys ...string) error {
	var emptyKeys []string
	for _, k := range keys {
		if v := bm.GetString(k); v == NULL {
			emptyKeys = append(emptyKeys, k)
		}
	}
	if len(emptyKeys) > 0 {
		return errors.New(strings.Join(emptyKeys, ", ") + " : cannot be empty")
	}
	return nil
}

func convertToString(v interface{}) (str string) {
	if v == nil {
		return NULL
	}
	var (
		bs  []byte
		err error
	)
	if bs, err = json.Marshal(v); err != nil {
		return NULL
	}
	str = string(bs)
	return
}

// ToMD5 MD5加密64位
func ToMD5(old string) (neww string) {
	neww = fmt.Sprintf("%x", md5.Sum([]byte(old)))
	return
}
