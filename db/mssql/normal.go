package mssql

import (
	"database/sql"
	"reflect"
	"sync/atomic"
	"time"
)

// conn base on sql.DB
// in future maybe can add
// more
type conn struct {
	*sql.DB
}

// Query extend sql.DB
func (db *DB) Query(sql string, params ...interface{}) (rows *sql.Rows, err error) {
	var c *conn
	if len(db.read) > 0 {
		c = db.read[db.readIndex()]
	} else {
		c = db.writer
	}
	a := time.Now()
	rows, err = c.Query(sql, params...)
	debugLogQueies(sql, a, err, params...)
	return
}

// Query extend sql.DB
func (db *DB) Begin() (tx *Tx, err error) {
	tx = new(Tx)
	dtx, err := db.writer.Begin()
	if err != nil {
		return
	}
	tx.tx = dtx
	tx.db = db.writer
	sqlLogger.Info("BEGIN")
	return
}

// readIndex
func (db *DB) readIndex() int {
	if len(db.read) == 0 {
		return 0
	}
	v := atomic.AddInt64(&db.idx, 1)
	return int(v) % len(db.read)
}

// Exec sql 更新或者插入使用
func (db *DB) Exec(sql string, params ...interface{}) (ret sql.Result, err error) {
	a := time.Now()
	ret, err = db.writer.Exec(sql, params...)
	debugLogQueies(sql, a, err, params...)
	return
}

func (db *DB) Insert(obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	t = t.Elem()
	v = v.Elem()
	tableName := getTableName(t, v)
	query, params := insert(t, v, tableName)
	_, err = db.Exec(query, params...)
	return
}

// Update
func (db *DB) Update(obj interface{}) (id string, err error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	t = t.Elem()
	v = v.Elem()
	tableName := getTableName(t, v)
	id, query, params := update(t, v, tableName)
	_, err = db.Exec(query, params...)
	return
}

// InsertOrUpdate InsertOrUpdate
func (db *DB) InsertOrUpdate(obj interface{}) (id string, err error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	t = t.Elem()
	v = v.Elem()
	tableName := getTableName(t, v)
	_, pk := getPk(t, v)
	var query string
	var params []interface{}
	if pk == "" {
		query, params = insert(t, v, tableName)
	} else {
		id, query, params = update(t, v, tableName)
	}
	_, err = db.Exec(query, params...)
	return
}
