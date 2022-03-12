package mssql

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

func getPk(t reflect.Type, v reflect.Value) (pkkey string, pk string) {
	structname := t.Name()
	if k, ok := pkx[structname]; ok {
		vs := t.Field(k).Tag.Get(Tag)
		pkkey = strings.Split(vs, separate)[0]
		pk = v.Field(k).String()
		return
	}
	for k := 0; k < t.NumField(); k++ {
		vs := t.Field(k).Tag.Get(Tag)
		if vs != "" {
			if strings.Index(vs, pkindex) > -1 {
				vs = strings.Split(vs, separate)[0]
				pkkey = vs
				pk = v.Field(k).String()
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
	is, okf := insertmap[structname]
	var field []string
	for k := 0; k < t.NumField(); k++ {
		tagv := t.Field(k).Tag.Get(Tag)
		tagvs := strings.Split(tagv, separate)
		vs := tagvs[0]
		if vs != "" {
			if !okf {
				field = append(field, vs)
			}
			switch v.Field(k).Interface().(type) {
			case time.Time:
				param = append(param, v.Field(k).Interface().(time.Time).Format(timeDefaultFormat))
			default:
				param = append(param, v.Field(k).Interface())
			}
		}
	}
	if !okf {
		fields := strings.Join(field, ",")
		locki.Lock()
		insertmap[structname] = is
		locki.Unlock()
		is = insertstruct{
			keys:   fields,
			length: len(field),
		}
	}
	query = "INSERT INTO " + tablename + " (" + is.keys + ") VALUES (?" + strings.Repeat(",?", is.length-1) + ")"
	return
}

func insert2(t reflect.Type, v reflect.Value, tablename string) (query string, param []interface{}) {
	var field []string
	for k := 0; k < t.NumField(); k++ {
		tagv := t.Field(k).Tag.Get(Tag)
		tagvs := strings.Split(tagv, separate)
		if strings.Index(tagv, pkindex) > -1 {
			continue
		}
		vs := tagvs[0]
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
	is := insertstruct{
		keys:   strings.Join(field, ","),
		length: len(field),
	}
	query = "INSERT INTO " + tablename + " (" + is.keys + ") VALUES (?" + strings.Repeat(",?", is.length-1) + ")"
	return
}

func update(t reflect.Type, v reflect.Value, tableName string) (pk string, query string, param []interface{}) {
	structname := t.Name()
	is, okf := updatemap[structname]
	var field []string
	for k := 0; k < t.NumField(); k++ {
		vs := t.Field(k).Tag.Get(Tag)
		if vs != "" {
			if strings.Index(vs, pkindex) > -1 {
				pk = v.Field(k).String()
				continue
			}
			if !okf {
				field = append(field, vs+"=?")
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
	if !okf {
		is = insertstruct{
			keys: strings.Join(field, ","),
		}
		locku.Lock()
		updatemap[structname] = is
		locku.Unlock()
	}
	// }
	query = "UPDATE " + tableName + " SET " + is.keys + " WHERE id=?"
	param = append(param, pk)
	return
}
