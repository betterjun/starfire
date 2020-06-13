package middleware

import (
	"bytes"
	"github.com/betterjun/starfire/pkg/logs"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

// 记录请求信息
func LogRequestMiddleware(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	logs.Sugar().Infow("received request", "proto", c.Request.Proto, "method", c.Request.Method, "uri", c.Request.RequestURI, "body", string(body))

	// 把读过的字节流重新放到body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	c.Next()
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// 记录响应消息
func LogResponseMiddleware(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()

	logs.Sugar().Infow("send response", "proto", c.Request.Proto, "method", c.Request.Method, "uri", c.Request.RequestURI, "body", blw.body.String())
}
