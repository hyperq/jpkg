package dao

import (
	"fmt"
	"github.com/hyperq/jpkg/conf"
	"reflect"
	"strconv"
	"sync"

	"github.com/hyperq/jpkg/cache"
	"github.com/hyperq/jpkg/db/qs"
	"github.com/hyperq/jpkg/log"

	jsoniter "github.com/json-iterator/go"
)

var tablemap = sync.Map{}
var RC *cache.RC

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

func GetTableName(t reflect.Type) string {
	structname := t.Name()
	if tablenamecache, ok := tablemap.Load(structname); ok {
		return tablenamecache.(string)
	}
	v := reflect.New(t)
	m := v.MethodByName("TableName")
	if m != empthv {
		tablename := m.Call([]reflect.Value{})[0].String()
		tablemap.Store(structname, tablename)
		return tablename
	}
	return structname
}

func FindByID(id interface{}, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	t = t.Elem()
	tableName := GetTableName(t)
	q := qs.New().EQ("id", id)
	sql := "SELECT a.* FROM " + tableName + " a"
	return QueryByQs(sql, q, obj)
}

func FindByIDCache(id interface{}, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	t = t.Elem()
	tableName := GetTableName(t)
	cachekey := conf.Config.AppName + tableName + "-d" + fmt.Sprint(id)
	res, err := RC.GET(cachekey)
	if err != nil {
		q := qs.New().EQ("id", id)
		sql := "SELECT a.* FROM " + tableName + " a"
		err = QueryByQs(sql, q, obj)
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
	tableName := GetTableName(t)
	sql := "SELECT a.* FROM " + tableName + " a"
	return QueryByQs(sql, q, obj)
}

func FindsCache(q *qs.QuerySet, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	tableName := GetTableName(t)
	cachekey := conf.Config.AppName + q.FormatCache(tableName+"-s")
	res, err := RC.GET(cachekey)
	if err != nil {
		sql := "SELECT a.* FROM " + tableName + " a"
		err = QueryByQs(sql, q, obj)
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

func FindsCacheByTableName(q *qs.QuerySet, obj interface{}, tableName string) (err error) {
	cachekey := conf.Config.AppName + q.FormatCache(tableName+"-s")
	res, err := RC.GET(cachekey)
	if err != nil {
		sql := "SELECT a.* FROM " + tableName + " a"
		err = QueryByQs(sql, q, obj)
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

func FindsByTableName(q *qs.QuerySet, obj interface{}, tableName string) (err error) {
	sql := "SELECT a.* FROM " + tableName + " a"
	err = QueryByQs(sql, q, obj)
	return
}

func FindByKey(key string, value interface{}, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	t = t.Elem()
	tableName := GetTableName(t)
	q := qs.New().EQ(key, value)
	sql := "SELECT a.* FROM " + tableName + " a"
	return QueryByQs(sql, q, obj)
}

func Find(q *qs.QuerySet, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	t = t.Elem()
	tableName := GetTableName(t)
	sql := "SELECT a.* FROM " + tableName + " a"
	return QueryByQs(sql, q, obj)
}

func FindCustom(sql string, q *qs.QuerySet, obj interface{}) (err error) {
	log.Debug(sql)
	return QueryByQs(sql, q, obj)
}

func FindsCacheCustom(sql string, q *qs.QuerySet, obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	t = t.Elem()
	t = t.Elem()
	tableName := GetTableName(t)
	cachekey := conf.Config.AppName + q.FormatCache(tableName+"-c")
	res, err := RC.GET(cachekey)
	if err != nil {
		err = QueryByQs(sql, q, obj)
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
