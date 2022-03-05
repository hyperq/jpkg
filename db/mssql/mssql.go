package mssql

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hyperq/jpkg/log"
	"reflect"
)

const timeDefaultFormat = "2006-01-02 15:04:05"

var (
	Tag       = "gorm"
	TableName = "TableName"
	separate  = ";"
	pkindex   = separate + "pk"
)

var empthv = reflect.Value{}

type Config struct {
	DSN     string   // write data source name.
	ReadDSN []string // read data source name.
	Active  int      // pool
	Idle    int      // pool
}

type DB struct {
	writer *conn   // for writer db
	read   []*conn // for many read db
	idx    int64   // for read count
	tx     *conntx
}

type Rows struct {
	*sql.Rows
	cancel func()
}

func New(conf Config) (db *DB, err error) {
	db = new(DB)
	db.writer, err = connect(conf.DSN, conf.Idle, conf.Active)
	if err != nil {
		return
	}
	var read []*conn
	for _, v := range conf.ReadDSN {
		r, err1 := connect(v, conf.Idle, conf.Active)
		if err1 != nil {
			err = err1
			return
		}
		read = append(read, r)
	}
	db.read = read
	return
}
func connect(dsn string, idle, active int) (*conn, error) {
	d, err := sql.Open("mssql", dsn)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	d.SetMaxIdleConns(idle)
	d.SetMaxOpenConns(active)
	c := new(conn)
	c.DB = d
	return c, nil
}
