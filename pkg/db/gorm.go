package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
)

// 默认数据库实例
var defaultDB *gorm.DB

func InitSqliteDB(dbfile string) (err error) {
	defaultDB, err = gorm.Open("sqlite3", dbfile)
	initDBCallbacks(defaultDB)
	return err
}

func InitMysqlDB(host, username, password, dbname string, maxOpen, maxIdle int) (err error) {
	defaultDB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, dbname))
	if err != nil {
		return err
	}

	initDBCallbacks(defaultDB)

	// 设置链接数量
	defaultDB.DB().SetMaxOpenConns(maxOpen)
	defaultDB.DB().SetMaxIdleConns(maxIdle)
	//defaultDB.DB().SetConnMaxLifetime(d time.Duration)

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
	//return defaultDB.New()
}
