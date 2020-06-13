package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
)

// 默认数据库实例
var defaultDB *gorm.DB

// 初始化默认sqlite数据库
func InitSqliteDB(dbfile string) (err error) {
	d, err := NewSqliteDB(dbfile)
	if err != nil {
		return err
	}
	defaultDB = d
	return nil
}

// 初始化默认mysql数据库
func InitMysqlDB(host, username, password, dbname string, maxOpen, maxIdle int) (err error) {
	d, err := NewMysqlDB(host, username, password, dbname, maxOpen, maxIdle)
	if err != nil {
		return err
	}
	defaultDB = d
	return nil
}

// 关闭数据库。
func Close() (err error) {
	if defaultDB != nil {
		return defaultDB.Close()
	}
	return nil
}

// 获取数据库连接对象。
func GetDB() *gorm.DB {
	return defaultDB.Debug()
}

// 初始化sqlite数据库
func NewSqliteDB(dbfile string) (d *gorm.DB, err error) {
	d, err = gorm.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	initDBCallbacks(d)
	return d, nil
}

// 初始化mysql数据库
func NewMysqlDB(host, username, password, dbname string, maxOpen, maxIdle int) (d *gorm.DB, err error) {
	d, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, dbname))
	if err != nil {
		return nil, err
	}

	// 设置链接数量
	d.DB().SetMaxOpenConns(maxOpen)
	d.DB().SetMaxIdleConns(maxIdle)

	initDBCallbacks(d)
	return d, nil
}
