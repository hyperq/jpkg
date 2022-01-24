package mysql

import (
	"database/sql"
	"github.com/hyperq/jpkg/log"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const timeDefaultFormat = "2006-01-02 15:04:05"

var (
	Tag       = "gorm"
	TableName = "TableName"
	separate  = ";"
	pkindex   = separate + "pk"
)

var empthv = reflect.Value{}

// DB sql db struct
type DB struct {
	writer *conn   // for writer db
	read   []*conn // for many read db
	idx    int64   // for read count
	tx     *conntx
}

// Config mysql config struct
type Config struct {
	DSN     string   // write data source name.
	ReadDSN []string // read data source name.
	Active  int      // pool
	Idle    int      // pool
	mongodb string
}

// Rows struct
type Rows struct {
	*sql.Rows
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

// connect mysql
func connect(dsn string, idle, active int) (*conn, error) {
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	d.SetMaxIdleConns(idle)
	d.SetMaxOpenConns(active)
	d.SetConnMaxLifetime(14400 * time.Second)
	c := new(conn)
	c.DB = d
	return c, nil
}
