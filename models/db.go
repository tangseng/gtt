package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)

type DB struct {
	Host string
	Port string
	User string
	Pass string
	Protocol string
}

func NewDB(host, user, pass string) *DB {
	return &DB{
		Host : host,
		Port : "3306",
		User : user,
		Pass : pass,
		Protocol : "tcp",
	}
}

type DBDo struct {
	db *DB
	database string
	orm orm.Ormer
}

func NewDBDo(db *DB, database string) *DBDo {
	return &DBDo{
		db : db,
		database : database,
		orm : orm.NewOrm(),
	}
}

func (this *DBDo) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", this.db.User, this.db.Pass, this.db.Host, this.db.Port)
}

func (this *DBDo) dsnd(database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", this.db.User, this.db.Pass, this.db.Host, this.db.Port, database)
}

func (this *DBDo) Check() error {
	conn, err := sql.Open("mysql", this.dsn())
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Ping()
}

func (this *DBDo) RDB(host, database string) error {
	alias := fmt.Sprintf("%s", host)
	if database != "" {
		alias = fmt.Sprintf("%s_%s", host, database)
	}
	_, err := orm.GetDB(alias)
	if err != nil {
		dsn := ""
		if database != "" {
			dsn = this.dsnd(database)
		} else {
			dsn = this.dsn()
		}
		err = orm.RegisterDataBase(alias, "mysql", dsn)
		if err != nil {
			return err
		}
	}
	this.orm.Using(alias)
	return nil
}

func (this *DBDo) ShowData(table string, number int) ([]orm.Params, error) {
	err := this.RDB(this.db.Host, this.database)
	if err != nil {
		return nil, err
	}
	sql := ""
	if number > 0 {
		sql = fmt.Sprintf(`SELECT * FROM %s ORDER BY id DESC LIMIT %d`, table, number)
	} else {
		sql = fmt.Sprintf(`SELECT * FROM %s ORDER BY id DESC`, table)
	}
	var values []orm.Params
	_, err = this.orm.Raw(sql).Values(&values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
