package model

import "github.com/betterjun/starfire/pkg/db"

type User struct {
	db.Model

	// 用户名
	Username string `json:"username"`
	// 密码
	Password string `json:"password"`
	// 用户状态
	State uint `json:"state"`
	// 用户头像
	Avatar string `json:"avatar"`
}

// 用户的状态
const (
	// 正常
	UserStateNormal = 1
	// 禁用
	UserStateDisable = 2
)
