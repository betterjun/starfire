package handler

import (
	"encoding/json"
	"fmt"
	"github.com/betterjun/starfire/app/webserver/errcode"
	"github.com/betterjun/starfire/app/webserver/params"
	"github.com/betterjun/starfire/app/webserver/svc"
	"github.com/betterjun/starfire/app/webserver/util"
	"github.com/betterjun/starfire/pkg/kvdb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 鉴权处理器，pri的接口都需要鉴权才能访问
func Authorize(c *gin.Context) {
	token := c.GetHeader("X-Token")
	if token == "" {
		// 没有登陆过
		c.AbortWithStatusJSON(http.StatusOK, errcode.Resp(errcode.NoToken))
		return
	}

	tokenData, err := kvdb.Get(util.FormatTokenUserKey(token))
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			c.AbortWithStatusJSON(http.StatusOK, errcode.Resp(errcode.InvalidToken))
		} else {
			c.AbortWithStatusJSON(http.StatusOK, errcode.Resp(errcode.InternalError))
		}
		return
	}

	ut := util.UserToken{}
	err = json.Unmarshal(tokenData, &ut)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.Resp(errcode.InvalidToken))
		return
	}

	if ut.UID == 0 || ut.ExpireTime < time.Now().Unix() {
		c.AbortWithStatusJSON(http.StatusOK, errcode.Resp(errcode.TokenExpired))
		return
	}

	// 在header中追加用户id
	c.Request.Header.Set("X-UID", fmt.Sprint(ut.UID))
	c.Next()
}

// @Summary 用户注册
// @Description 用户注册
// @Tags 用户
// @Accept  json
// @Param	signUpParam			body	params.SignupReq			true	"注册参数"
// @Produce  json
// @Success 200 {string} string "{"code":200,"data":"ok"}"
// @Router /v1/pub/signup [POST]
// @ID SignUp
func SignUp(c *gin.Context) {
	var param params.SignupReq
	err := c.ShouldBind(&param)
	if err != nil {
		c.JSON(http.StatusOK, errcode.Resp(errcode.InvalidParams, err.Error()))
		return
	}

	ae := svc.SignUp(&param)
	c.JSON(http.StatusOK, errcode.Resp(ae))
}

// @Summary 用户登录
// @Description 用户登录
// @Tags 用户
// @Accept  json
// @Param	signInParam			body	params.SigninReq			true	"登录参数"
// @Produce  json
// @Success 200 {object} params.SigninRsp
// @Router /v1/pub/signin [POST]
// @ID SignIn
func SignIn(c *gin.Context) {
	var param params.SigninReq
	err := c.ShouldBind(&param)
	if err != nil {
		c.JSON(http.StatusOK, errcode.Resp(errcode.InvalidParams, err.Error()))
		return
	}

	rsp, ae := svc.SignIn(&param)
	c.JSON(http.StatusOK, errcode.Resp(ae, rsp))
}

// @Summary 用户登出
// @Description 用户登出
// @Tags 用户
// @Accept  json
// @Produce  json
// @Success 200 {string} string "{"code":200,"data":"ok"}"
// @Router /v1/pri/signout [POST]
// @ID SignOut
func SignOut(c *gin.Context) {
	ae := svc.SignOut(getToken(c))
	c.JSON(http.StatusOK, errcode.Resp(ae))
}

// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Success 200 {object} params.UserInfoRsp
// @Router /v1/pri/userinfo [GET]
// @ID UserInfo
func UserInfo(c *gin.Context) {
	uid := getUID(c)
	if uid == 0 {
		c.JSON(http.StatusOK, errcode.Resp(errcode.InvalidParams))
	}

	rsp, ae := svc.UserInfo(uid)
	c.JSON(http.StatusOK, errcode.Resp(ae, rsp))
}

func getUID(c *gin.Context) uint {
	id, _ := strconv.ParseInt(c.Request.Header.Get("X-UID"), 10, 64)
	return uint(id)
}

func getToken(c *gin.Context) string {
	return c.Request.Header.Get("X-Token")
}
