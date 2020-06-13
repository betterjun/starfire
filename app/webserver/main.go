package main

import (
	"flag"
	"fmt"
	"github.com/betterjun/starfire/app/webserver/model"
	"github.com/betterjun/starfire/app/webserver/router"
	"github.com/betterjun/starfire/pkg/boltdb"
	"github.com/betterjun/starfire/pkg/cfg"
	"github.com/betterjun/starfire/pkg/db"
	"github.com/betterjun/starfire/pkg/kvdb"
	"github.com/betterjun/starfire/pkg/logs"
	"github.com/betterjun/starfire/pkg/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
)

// @title starfire量化系统 API
// @version 1.0
// @description pub开头的api不需要登录就可访问，pri开头的需要登录才能访问，访问pri开头的路由时，需要把登录返回的token放到header X-Token中，服务器要做鉴权。
// @termsOfService http://swagger.io/terms/

// @contact.name StarFire API Support
// @contact.url http://www.starfire.io/support
// @contact.email starfire@starfire.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:16888
// @BasePath /
func main() {
	// 1 命令行参数解析
	mode := flag.String("m", "dev", "指定执行模式,只支持[dev|test|prod],默认是dev")
	flag.Parse()
	dev := true
	if *mode != "dev" {
		dev = false
	}
	if dev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 2 载入配置文件
	configFile := fmt.Sprintf("conf/%s.toml", *mode)
	fmt.Printf("使用的配置文件:%s..........", configFile)
	err := cfg.Initialize(configFile)
	if err != nil {
		fmt.Printf("Error reading configuration[%s]: %s\n", configFile, err.Error())
		os.Exit(1)
	}

	// 3 初始化日志
	logs.Init(cfg.GetString("log.file"), cfg.GetString("log.level"), cfg.GetInt("log.maxSize"),
		cfg.GetInt("log.maxAge"), cfg.GetInt("log.maxBackup"), dev)

	// 4 sql数据库初始化
	err = InitDB()
	if err != nil {
		fmt.Printf("数据库初始化失败: %s\n", err.Error())
		logs.Info("数据库初始化失败", zap.Error(err))
		os.Exit(1)
	}
	model.Setup()

	// 5 nosql数据库初始化
	err = InitNosqlDB()
	if err != nil {
		fmt.Printf("redis初始化失败: %s\n", err.Error())
		os.Exit(1)
	}

	// 6 初始化路由
	err = router.Init()
	if err != nil {
		fmt.Printf("启动出错:%+v", err)
		return
	}

	// 7 启动webserver

	logs.Info("程序已启动")

	// 阻塞
	select {}
}

// 连接数据库
func InitDB() error {
	const SectionDB = "database"

	dbType := cfg.GetString(SectionDB + ".db_type")
	switch dbType {
	case "mysql":
		host := cfg.GetString(SectionDB + ".mysql_host")
		username := cfg.GetString(SectionDB + ".mysql_username")
		password := cfg.GetString(SectionDB + ".mysql_password")
		dbname := cfg.GetString(SectionDB + ".mysql_dbname")
		maxOpen := cfg.GetInt(SectionDB + ".mysql_max_open")
		maxIdle := cfg.GetInt(SectionDB + ".mysql_max_idle")
		return db.InitMysqlDB(host, username, password, dbname, maxOpen, maxIdle)
	case "sqlite3":
		dbfile := cfg.GetString(SectionDB + ".sqlite3_file")
		return db.InitSqliteDB(dbfile)
	default:
		return fmt.Errorf("数据库配置出错:不支持的数据库类型%q\n", dbType)
	}
}

// 连接nosql数据库
func InitNosqlDB() error {
	const SectionDB = "nosql"

	dbType := cfg.GetString(SectionDB + ".db_type")
	switch dbType {
	case "redis":
		host := cfg.GetString(SectionDB + ".redis_host")
		auth := cfg.GetString(SectionDB + ".redis_auth")
		db := cfg.GetInt(SectionDB + ".redis_db")
		maxActive := cfg.GetInt(SectionDB + ".redis_max_active")
		maxIdle := cfg.GetInt(SectionDB + ".redis_max_idle")
		idleTimeout := cfg.GetInt(SectionDB + ".redis_idle_timeout")

		nosqldb, err := redis.NewRedis(host, auth, db, maxActive, maxIdle, idleTimeout)
		if err != nil {
			return err
		}
		kvdb.Init(nosqldb)
	case "bolt":
		dbfile := cfg.GetString(SectionDB + ".bolt_file")
		nosqldb, err := boltdb.NewBoltdb(dbfile)
		if err != nil {
			return err
		}
		kvdb.Init(nosqldb)
	default:
		return fmt.Errorf("数据库配置出错:不支持的数据库类型%q\n", dbType)
	}

	return nil
}
