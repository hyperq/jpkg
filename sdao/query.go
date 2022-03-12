package sdao

import (
	"fmt"
	"github.com/hyperq/jpkg/cache"
	"github.com/hyperq/jpkg/conf"
	"github.com/hyperq/jpkg/db/qs"
	"github.com/hyperq/jpkg/log"
	"reflect"
	"strconv"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var tablemap = sync.Map{}
var RC *cache.RC

var (
	Tag       = "gorm"
	TableName = "TableName"
	separate  = ";"
	pkindex   = separate + "pk"
)

func ClearCache(tablename string, id interface{}) {
	ClearListCache(tablename)
	ClearOneCache(tablename, id)
}

func ClearCacheByIDs(tableName string, ids ...interface{}) {
	for _, id := range ids {
		ClearOneCache(tableName, id)
	}
}

func ClearListCache(tablename string) {
	tablename = conf.Config.AppName + tablename
	keys, err := RC.KEYS(tablename + "-s*")
	if err != nil {
		log.Error(err)
	}
	keyc, err := RC.KEYS(tablename + "-c*")
	if err != nil {
		log.Error(err)
	}
	keyn, err := RC.KEYS(tablename + "-n*")
	if err != nil {
		log.Error(err)
	}
	keys = append(keys, keyc...)
	keys = append(keys, keyn...)
	if len(keys) > 0 {
		err = RC.DEL(keys...)
		if err != nil {
			log.Error(err)
		}
	}
}

func ClearOneCache(tablename string, id interface{}) {
	_ = RC.DEL(conf.Config.AppName + tablename + "-d" + fmt.Sprint(id))
}

var empthv = reflect.Value{}

var pkx map[string]string
var lockp sync.Mutex

func GetTableName(t reflect.Type) (table, pkkey string) {
	structname := t.Name()
	if tablenamecache, ok := tablemap.Load(structname); ok {
		table = tablenamecache.(string)
	} else {
		v := reflect.New(t)
		m := v.MethodByName("TableName")
		if m != empthv {
			tablename := m.Call([]reflect.Value{})[0].String()
			tablemap.Store(structname, tablename)
			table = tablename
		}
	}
	kc, ok := pkx[structname]
	if !ok {
		for k := 0; k < t.NumField(); k++ {
			tagv := t.Field(k).Tag.Get(Tag)
			if strings.Index(tagv, pkindex) > -1 {
				pkkey = strings.Replace(tagv, pkindex, "", -1)
				lockp.Lock()
				pkx[structname] = pkkey
				lockp.Unlock()
				break
			}
		}
	} else {
		pkkey = kc
	}
	return
}

func FindByID(id interface{}, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	t = t.Elem()
	tableName, pkkey := GetTableName(t)
	q := qs.New2().EQ(pkkey, id).LIMIT(-1)
	return QueryByQs(tableName, pkkey, q, obj)
}

func FindByIDCache(id interface{}, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	t = t.Elem()
	tableName, pkkey := GetTableName(t)
	cachekey := conf.Config.AppName + tableName + "-d" + fmt.Sprint(id)
	res, err := RC.GET(cachekey)
	if err != nil {
		q := qs.New2().EQ("id", id).LIMIT(-1)
		err = QueryByQs(tableName, pkkey, q, obj)
		if err != nil {
			return
		}
		res, err = jsoniter.MarshalToString(obj)
		if err != nil {
			return
		}
		err = RC.SET(cachekey, res, 600)
		return
	}
	err = jsoniter.UnmarshalFromString(res, obj)
	return
}

func Finds(q *qs.QuerySet, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	t = t.Elem()
	t = t.Elem()
	tableName, pkkey := GetTableName(t)
	return QueryByQs(tableName, pkkey, q, obj)
}

func FindsCache(q *qs.QuerySet, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	tableName, pkkey := GetTableName(t)
	cachekey := conf.Config.AppName + q.FormatCache(tableName+"-s")
	res, err := RC.GET(cachekey)
	if err != nil {
		err = QueryByQs(tableName, pkkey, q, obj)
		if err != nil {
			return
		}
		res, err = jsoniter.MarshalToString(obj)
		if err != nil {
			return
		}
		err = RC.SET(cachekey, res, 600)
		return
	} else {
		err = jsoniter.UnmarshalFromString(res, obj)
	}
	return
}

func FindByKey(key string, value interface{}, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	t = t.Elem()
	tableName, pkkey := GetTableName(t)
	q := qs.New2().EQ(key, value).LIMIT(-1)
	return QueryByQs(tableName, pkkey, q, obj)
}

func Find(q *qs.QuerySet, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	t = t.Elem()
	tableName, pkkey := GetTableName(t)
	q.LIMIT(-1)
	return QueryByQs(tableName, pkkey, q, obj)
}

// CountCache 将count存入缓存
func CountCache(tablename string, q *qs.QuerySet) int {
	countkey := conf.Config.AppName + q.FormatCache(tablename+"-n")
	counts, err := RC.GET(countkey)
	if err != nil {
		counti := Count(tablename, q)
		err = RC.SET(countkey, strconv.Itoa(counti), 600)
		if err != nil {
			log.Error(err)
		}
		return counti
	}
	count, _ := strconv.Atoi(counts)
	return count
}

func GetCacheKey(table, suffix string, id interface{}) string {
	return conf.Config.AppName + table + "_" + suffix + fmt.Sprint(id)
}
