package dao

import (
	"database/sql"
	"strings"

	"github.com/hyperq/jpkg/db/mysql"
	"github.com/hyperq/jpkg/db/qs"
	"github.com/hyperq/jpkg/log"
	"github.com/hyperq/jpkg/tool"

	"github.com/didi/gendry/scanner"
)

func Query(sql string, params ...interface{}) (*sql.Rows, error) {
	return Db.Query(sql, params...)
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
	return scanner.ScanClose(rows, data)
}

func Exec(sql string, params ...interface{}) (sql.Result, error) {
	return Db.Exec(sql, params...)
}

func Insert(res interface{}) (int64, error) {
	return Db.Insert(res)
}

func Update(res interface{}) (int64, error) {
	return Db.Update(res)
}

func InsertOrUpdate(res interface{}) (int64, error) {
	return Db.InsertOrUpdate(res)
}

func StatusChange(table string, id, version interface{}, key, status string) error {
	_, err := Exec(
		"UPDATE "+table+" SET "+key+" = ?,modify_time = ?,version = version + 1 WHERE id = ? AND version = ?",
		status, tool.GetNows(), id, version,
	)
	return err
}

func StatusChangeShop(table string, id, version, shopid interface{}, key, status string) error {
	_, err := Exec(
		"UPDATE "+table+" SET "+key+" = ?,modify_time = ?,version = version + 1 WHERE id = ? AND version = ? AND shop_id=?",
		status, tool.GetNows(), id, version, shopid,
	)
	return err
}

func Begin() (*mysql.Tx, error) {
	return Db.Begin()
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
	err = scanner.ScanClose(rows, &data)
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
	err = scanner.ScanClose(rows, &data)
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
	err = scanner.ScanClose(rows, &data)
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
	err = scanner.ScanClose(rows, &data)
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
	err = scanner.ScanClose(rows, &data)
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
	err = scanner.ScanClose(rows, &data)
	if err != nil {
		log.Error2(err)
		return
	}
	if len(data) > 0 {
		num = data[0].Num
	}
	return
}
