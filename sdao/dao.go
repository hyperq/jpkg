package dao

import (
	"database/sql"
	"github.com/hyperq/jpkg/db"
	"github.com/hyperq/jpkg/db/mssql"
	"strings"

	"github.com/hyperq/jpkg/db/qs"
	"github.com/hyperq/jpkg/log"
)

func Query(sql string, params ...interface{}) (*sql.Rows, error) {
	return db.Msdb.Query(sql, params...)
}

func QueryByQs(sql string, q *qs.QuerySet, data interface{}) error {
	wh, params, others := q.Format()
	s := "a.*"
	if len(q.Select) > 0 {
		s = strings.Join(q.Select, ",")
	}
	sql = strings.Replace(sql, "{{select}}", s, 1)
	rows, err := Query(sql+wh+others, params...)
	if err != nil {
		return err
	}
	return mssql.ScanClose(rows, data)
}

func Exec(sql string, params ...interface{}) (sql.Result, error) {
	return db.Msdb.Exec(sql, params...)
}

func Insert(res interface{}) error {
	return db.Msdb.Insert(res)
}

func Update(res interface{}) (string, error) {
	return db.Msdb.Update(res)
}

func InsertOrUpdate(res interface{}) (string, error) {
	return db.Msdb.InsertOrUpdate(res)
}

func Begin() (*mssql.Tx, error) {
	return db.Msdb.Begin()
}

type count struct {
	Num int `gorm:"num" json:"num"`
}

func Count(table string, q *qs.QuerySet) (num int) {
	wh, params := q.FormatWhere()
	rows, err := Query(
		`
		SELECT count(*) as num
		FROM `+table+` a
		`+wh, params...,
	)
	if err != nil {
		log.Error2(err)
		return
	}
	var data []count
	err = mssql.ScanClose(rows, &data)
	if err != nil {
		log.Error2(err)
		return
	}
	if len(data) > 0 {
		num = data[0].Num
	}
	return
}

type sum struct {
	Num float64 `gorm:"sumnum" json:"sumnum"`
}

func Sum(table, key string, q *qs.QuerySet) (sumnum float64) {
	wh, params := q.FormatWhere()
	rows, err := Query(
		`
		SELECT sum(`+key+`) as sumnum
		FROM `+table+` a
		`+wh, params...,
	)
	if err != nil {
		log.Error2(err)
		return
	}
	var data []sum
	err = mssql.ScanClose(rows, &data)
	if err != nil {
		log.Error2(err)
		return
	}
	if len(data) > 0 {
		sumnum = data[0].Num
	}
	return
}

func SumSql(sql string, q *qs.QuerySet) (sumnum float64) {
	wh, params := q.FormatWhere()
	rows, err := Query(sql+wh, params...)
	if err != nil {
		log.Error2(err)
		return
	}
	var data []sum
	err = mssql.ScanClose(rows, &data)
	if err != nil {
		log.Error2(err)
		return
	}
	if len(data) > 0 {
		sumnum = data[0].Num
	}
	return
}

func SumSqlC(sql string, params ...interface{}) (sumnum float64) {
	rows, err := Query(sql, params...)
	if err != nil {
		log.Error2(err)
		return
	}
	var data []sum
	err = mssql.ScanClose(rows, &data)
	if err != nil {
		log.Error2(err)
		return
	}
	if len(data) > 0 {
		sumnum = data[0].Num
	}
	return
}

func Countsql(sql string, q *qs.QuerySet) (num int) {
	wh, params := q.FormatWhere()
	rows, err := Query(sql+wh, params...)
	if err != nil {
		log.Error2(err)
		return
	}
	var data []count
	err = mssql.ScanClose(rows, &data)
	if err != nil {
		log.Error2(err)
		return
	}
	if len(data) > 0 {
		num = data[0].Num
	}
	return
}
func CountsqlC(sql string, params ...interface{}) (num int) {
	rows, err := Query(sql, params...)
	if err != nil {
		log.Error2(err)
		return
	}
	var data []count
	err = mssql.ScanClose(rows, &data)
	if err != nil {
		log.Error2(err)
		return
	}
	if len(data) > 0 {
		num = data[0].Num
	}
	return
}
