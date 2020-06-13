package router

import (
	"fmt"
	_ "github.com/betterjun/starfire/app/webserver/docs"
	"github.com/betterjun/starfire/app/webserver/handler"
	"github.com/betterjun/starfire/pkg/cfg"
	"github.com/betterjun/starfire/pkg/logs"
	"github.com/betterjun/starfire/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
	"net/http"
)

func Init() error {
	r := gin.Default()

	// 配置了路径，才启用文件服务，部署时可以选择将前端的代码一起部署
	webPath := cfg.GetString("common.webpath")
	if len(webPath) > 0 {
		r.Static("/web", webPath)
	}

	// swagger文档
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 跨域
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowMethods:    []string{http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete, http.MethodOptions},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "X-Token"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowCredentials:       true,
		ExposeHeaders:          nil,
		MaxAge:                 0,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        true,
		AllowFiles:             true,
	}))

	v1 := r.Group("/v1")
	v1.Use(middleware.LogRequestMiddleware)
	v1.Use(middleware.LogResponseMiddleware)
	pub := v1.Group("/pub")
	{
		// 不需要登录即可访问的api
		// 登录
		pub.POST("/signin", handler.SignIn)
		// 注册
		pub.POST("/signup", handler.SignUp)
	}

	pri := v1.Group("/pri", handler.Authorize)
	{
		// 需要登录才可访问的api
		// 登出
		pri.POST("/signout", handler.SignOut)

		// 获取用户信息
		pri.GET("/userinfo", handler.UserInfo)
	}

	if err := r.Run(fmt.Sprintf(":%s", cfg.GetString("common.port"))); err != nil {
		logs.Error("运行出错:%+v", zap.Error(err))
		return err
	}
	return nil
}
