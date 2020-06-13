package params

// 注册请求
type SignupReq struct {
	// 用户名称
	Username string `json:"username" form:"username" binding:"required"`
	// 用户密码
	Password string `json:"password" form:"password" binding:"required"`
}

// 登陆请求
type SigninReq struct {
	// 用户名称
	Username string `json:"username" form:"username" binding:"required"`
	// 用户密码
	Password string `json:"password" form:"password" binding:"required"`
}

// 登陆响应
type SigninRsp struct {
	// 用户id
	UID uint `json:"uid" form:"uid" binding:"required"`
	// 用户令牌，后续请求都需要在header X-Token带上此token
	Token string `json:"token" form:"token" binding:"required"`
}

// 获取用户信息响应
type UserInfoRsp struct {
	// 用户id
	UID uint `json:"uid"`
	// 用户名
	Name string `json:"name"`
	// 用户头像
	Avatar string `json:"avatar"`
}
