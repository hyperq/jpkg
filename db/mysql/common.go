package mysql

import (
	"reflect"
	"strings"
	"sync"
	"time"
)

var tablemap map[string]string
var lockt sync.Mutex

var pkx map[string]int
var lockp sync.Mutex

var insertmap map[string]insertstruct
var locki sync.Mutex

var updatemap map[string]insertstruct
var locku sync.Mutex

type insertstruct struct {
	keys   string
	length int
}

func init() {
	tablemap = make(map[string]string)
	pkx = make(map[string]int)
	insertmap = make(map[string]insertstruct)
	updatemap = make(map[string]insertstruct)
}

func getTableName(t reflect.Type, v reflect.Value) string {
	structname := t.Name()
	if tablenamecache, ok := tablemap[structname]; ok {
		return tablenamecache
	}
	m := v.MethodByName(TableName)
	if m != empthv {
		tablename := m.Call([]reflect.Value{})[0].String()
		lockt.Lock()
		tablemap[structname] = tablename
		lockt.Unlock()
		return tablename
	}
	return structname
}

func getPk(t reflect.Type, v reflect.Value) (pkkey string, pk int64) {
	structname := t.Name()
	if k, ok := pkx[structname]; ok {
		vs := t.Field(k).Tag.Get(Tag)
		pkkey = strings.Split(vs, separate)[0]
		pk = v.Field(k).Int()
		return
	}
	for k := 0; k < t.NumField(); k++ {
		vs := t.Field(k).Tag.Get(Tag)
		if vs != "" {
			if strings.Index(vs, pkindex) > -1 {
				vs = strings.Split(vs, separate)[0]
				pkkey = vs
				pk = v.Field(k).Int()
				lockp.Lock()
				pkx[structname] = k
				lockp.Unlock()
				return
			}
		}
	}
	return
}

func insert(t reflect.Type, v reflect.Value, tablename string) (query string, param []interface{}) {
	structname := t.Name()
	kc, ok := pkx[structname]
	is, okf := insertmap[structname]
	if okf {
		for k := 0; k < t.NumField(); k++ {
			tagv := t.Field(k).Tag.Get(Tag)
			tagvs := strings.Split(tagv, separate)
			vs := tagvs[0]
			if !ok {
				if strings.Index(tagv, pkindex) > -1 {
					lockp.Lock()
					pkx[structname] = k
					lockp.Unlock()
					continue
				}
			} else {
				if k == kc {
					continue
				}
			}
			if vs != "" {
				switch v.Field(k).Interface().(type) {
				case time.Time:
					param = append(param, v.Field(k).Interface().(time.Time).Format(timeDefaultFormat))
				default:
					param = append(param, v.Field(k).Interface())
				}
			}
		}
	} else {
		var field []string
		for k := 0; k < t.NumField(); k++ {
			tagv := t.Field(k).Tag.Get(Tag)
			tagvs := strings.Split(tagv, separate)
			vs := tagvs[0]
			if !ok {
				if strings.Index(tagv, pkindex) > -1 {
					lockp.Lock()
					pkx[structname] = k
					lockp.Unlock()
					continue
				}
			} else {
				if k == kc {
					continue
				}
			}
			if vs != "" {
				field = append(field, vs)
				switch v.Field(k).Interface().(type) {
				case time.Time:
					param = append(param, v.Field(k).Interface().(time.Time).Format(timeDefaultFormat))
				default:
					param = append(param, v.Field(k).Interface())
				}
			}
		}
		fields := strings.Join(field, ",")
		is = insertstruct{
			keys:   fields,
			length: len(field),
		}
		locki.Lock()
		insertmap[structname] = is
		locki.Unlock()
	}
	query = "INSERT INTO " + tablename + " (" + is.keys + ") VALUES (?" + strings.Repeat(",?", is.length-1) + ")"
	return
}

func update(t reflect.Type, v reflect.Value, tableName string) (pk int64, query string, param []interface{}) {
	structname := t.Name()
	is, okf := updatemap[structname]
	if okf {
		for k := 0; k < t.NumField(); k++ {
			vs := t.Field(k).Tag.Get(Tag)
			if vs != "" {
				if strings.Index(vs, pkindex) > -1 {
					pk = v.Field(k).Int()
					continue
				}
				switch v.Field(k).Interface().(type) {
				case time.Time:
					param = append(param, v.Field(k).Interface().(time.Time).Format(timeDefaultFormat))
				case bool:
					if v.Field(k).Interface().(bool) {
						param = append(param, 1)
					} else {
						param = append(param, 0)
					}
				default:
					param = append(param, v.Field(k).Interface())
				}
			}
		}
	} else {
		var field []string
		for k := 0; k < t.NumField(); k++ {
			vs := t.Field(k).Tag.Get(Tag)
			if vs != "" {
				if strings.Index(vs, pkindex) > -1 {
					pk = v.Field(k).Int()
					continue
				}
				field = append(field, vs+"=?")
				switch v.Field(k).Interface().(type) {
				case time.Time:
					param = append(param, v.Field(k).Interface().(time.Time).Format(timeDefaultFormat))
				case bool:
					if v.Field(k).Interface().(bool) {
						param = append(param, 1)
					} else {
						param = append(param, 0)
					}
				default:
					param = append(param, v.Field(k).Interface())
				}
			}
		}
		is = insertstruct{
			keys: strings.Join(field, ","),
		}
		locku.Lock()
		updatemap[structname] = is
		locku.Unlock()
	}
	query = "UPDATE " + tableName + " SET " + is.keys + " WHERE id=?"
	param = append(param, pk)
	return
}
func updateold(t reflect.Type, v reflect.Value, tableName string) (pk int64, query string, param []interface{}) {
	var field []string
	for k := 0; k < t.NumField(); k++ {
		vs := t.Field(k).Tag.Get(Tag)
		if vs != "" {
			if strings.Index(vs, pkindex) > -1 {
				pk = v.Field(k).Int()
				continue
			}
			field = append(field, vs+"=?")
			switch v.Field(k).Interface().(type) {
			case time.Time:
				param = append(param, v.Field(k).Interface().(time.Time).Format(timeDefaultFormat))
			default:
				param = append(param, v.Field(k).Interface())
			}
		}
	}
	query = "UPDATE " + tableName + " SET " + strings.Join(field, ",") + " WHERE id=?"
	param = append(param, pk)
	return
}
