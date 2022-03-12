package mssql

import (
	"database/sql"
	"reflect"
	"time"
)

// Tx transaction.
type Tx struct {
	db *conn
	tx *sql.Tx
}
type conntx struct {
	*sql.Tx
}

func (tx *Tx) Rollback() (err error) {
	sqlLogger.Info("ROLLBACK")
	return tx.tx.Rollback()
}

func (tx *Tx) Commit() (err error) {
	sqlLogger.Info("COMMIT")
	return tx.tx.Commit()
}

// Exec sql 更新或者插入使用
func (tx *Tx) Exec(sql string, params ...interface{}) (ret sql.Result, err error) {
	a := time.Now()
	ret, err = tx.tx.Exec(sql, params...)
	debugLogQueies(sql, a, err, params...)
	return
}

// Query extend sql.DB
func (tx *Tx) Query(sql string, params ...interface{}) (rows *sql.Rows, err error) {
	a := time.Now()
	rows, err = tx.tx.Query(sql, params...)
	debugLogQueies(sql, a, err, params...)
	return
}

// Insert insert
func (tx *Tx) Insert(obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	t = t.Elem()
	v = v.Elem()
	tableName := getTableName(t, v)
	query, params := insert(t, v, tableName)
	_, err = tx.Exec(query, params...)
	return
}

// Insert2 insert
func (tx *Tx) Insert2(obj interface{}) (err error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	t = t.Elem()
	v = v.Elem()
	tableName := getTableName(t, v)
	query, params := insert2(t, v, tableName)
	_, err = tx.Exec(query, params...)
	return
}

// Insert insert
func (tx *Tx) Update(obj interface{}) (id string, err error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	t = t.Elem()
	v = v.Elem()
	tableName := getTableName(t, v)
	id, query, params := update(t, v, tableName)
	_, err = tx.Exec(query, params...)
	return
}

// InsertOrUpdate InsertOrUpdate
func (tx *Tx) InsertOrUpdate(obj interface{}) (id string, err error) {
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
	_, err = tx.Exec(query, params...)
	return
}
