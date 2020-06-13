package svc

import (
	"encoding/json"
	"github.com/betterjun/starfire/app/webserver/errcode"
	"github.com/betterjun/starfire/app/webserver/model"
	"github.com/betterjun/starfire/app/webserver/params"
	"github.com/betterjun/starfire/app/webserver/util"
	"github.com/betterjun/starfire/pkg/db"
	"github.com/betterjun/starfire/pkg/kvdb"
	"github.com/betterjun/starfire/pkg/logs"
	"github.com/jinzhu/gorm"
	"time"
)

// 用户注册
func SignUp(param *params.SignupReq) errcode.APIError {
	c := 0
	err := db.GetDB().Model(&model.User{}).Where("username=?", param.Username).Count(&c).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return errcode.MysqlFailed
	}

	if c != 0 {
		logs.Sugar().Error("用户已注册", "username", param.Username)
		return errcode.UserAlreadyRegistered
	}

	// 验证通过,保存用户信息
	var user model.User
	user.Username = param.Username
	user.Password = param.Password
	user.State = model.UserStateNormal
	user.Avatar = "https://lh3.googleusercontent.com/proxy/8n9pUPUbZj_h7VUYTJWNrKkiWAmcjTGPPgXxQBBqUdNsWcvLmwTwKihB5yEl7BYbj7NgJlIlfxxa09oOxCbjsd4ZWoV9-GIv"
	err = db.GetDB().Create(&user).Error
	if err != nil {
		logs.Sugar().Error("保存用户信息出错", "error", err.Error())
		return errcode.MysqlFailed
	}

	return errcode.Success
}

// 用户登陆
func SignIn(param *params.SigninReq) (rsp *params.SigninRsp, ae errcode.APIError) {
	var user model.User
	err := db.GetDB().Model(&model.User{}).Where("username=?", param.Username).First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errcode.UserNotExists
		}
		return nil, errcode.MysqlFailed
	}

	if user.Password != param.Password {
		return nil, errcode.UsernameOrPasswordError
	}

	if user.State != model.UserStateNormal {
		return nil, errcode.UserDisabled
	}

	// 生成token
	rsp = &params.SigninRsp{}
	rsp.UID = user.ID
	rsp.Token = util.GenToken()

	ut := util.UserToken{
		Username:   user.Username,
		UID:        user.ID,
		Token:      rsp.Token,
		ExpireTime: time.Now().Unix() + util.TokenExpireTime,
	}

	tokenData, err := json.Marshal(ut)
	if err != nil {
		return nil, errcode.InternalError
	}

	// 记录uid和token的互相关联关系
	kvdb.Set(util.FormatUserTokenKey(rsp.UID), rsp.Token)
	kvdb.Set(util.FormatTokenUserKey(rsp.Token), tokenData)

	return rsp, errcode.Success
}

// 用户登出
func SignOut(token string) errcode.APIError {
	tokenData, err := kvdb.Get(util.FormatTokenUserKey(token))
	if err != nil {
		return errcode.RedisError
	}
	kvdb.Delete(util.FormatTokenUserKey(token))

	ut := util.UserToken{}
	err = json.Unmarshal(tokenData, &ut)
	if err != nil {
		return errcode.InternalError
	}
	kvdb.Delete(util.FormatUserTokenKey(ut.UID))

	return errcode.Success
}

// 获取用户信息
func UserInfo(uid uint) (rsp *params.UserInfoRsp, ae errcode.APIError) {
	var user model.User
	err := db.GetDB().Model(&model.User{}).Where("id=?", uid).First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errcode.UserNotExists
		}
		return nil, errcode.MysqlFailed
	}

	if user.State != model.UserStateNormal {
		return nil, errcode.UserDisabled
	}

	rsp = &params.UserInfoRsp{

		Name:   user.Username,
		Avatar: user.Avatar,
	}

	return rsp, errcode.Success
}
