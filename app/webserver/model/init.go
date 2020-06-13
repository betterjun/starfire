package model

import "github.com/betterjun/starfire/pkg/db"

func Setup() {
	db.GetDB().AutoMigrate(&User{})

	// 设置链接数量
	db.GetDB().DB().SetMaxIdleConns(10)
	db.GetDB().DB().SetMaxOpenConns(100)
}
