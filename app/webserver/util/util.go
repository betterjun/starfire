package util

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
)

// 用户token默认过期时间，单位秒
const TokenExpireTime = 86400 * 7

// 用户token信息
type UserToken struct {
	Username   string `json:"username"`
	UID        uint   `json:"uid"` // 用户id
	Token      string `json:"token"`
	ExpireTime int64  `json:"expire_time"` // 过期时间戳
}

// 生成token
func GenToken() string {
	return uuid.NewV4().String()
}

// uid找token
func FormatUserTokenKey(uid uint) string {
	return fmt.Sprintf("user2token_%v", uid)
}

// token找uid
func FormatTokenUserKey(token string) string {
	return fmt.Sprintf("token2user_%v", token)
}
