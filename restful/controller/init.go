package controller

import (
	"crypto-follower/core/config"
	"crypto-follower/core/log"
	"github.com/gin-gonic/gin"
	"time"
)

func Router(engine *gin.Engine) {
	index := engine.Group("/")
	{
		indexController := NewIndexController()
		index.GET("/welcome", indexController.Welcome)
		index.GET("/", indexController.Welcome)
	}
}

func Validator() {

}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志格式
		log.Infof("| %3d | %13v | %15s | %s | %s",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

func SetMode(config *config.ApplicationConfig) {
	gin.SetMode(gin.ReleaseMode)
}