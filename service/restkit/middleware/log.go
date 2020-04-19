package middleware

import (
	"github.com/gin-gonic/gin"
	"mizuki/project/core-kit/service/logkit"
	"time"
)

func Log() gin.HandlerFunc  {
	return func(c *gin.Context) {
		t := time.Now()
		sessionIdc ,err := c.Request.Cookie("session")
		sessionId := ""
		if err==nil{
			sessionId = sessionIdc.Value
		}
		// todo params
		//params := make(map[interface{}]interface{})
		//_ = c.Copy().ShouldBind(&params)
		logkit.Info("request",
			logkit.Param{
				Key: "session",
				Val: sessionId,
			}, logkit.Param{
				Key: "url",
				Val: c.Request.URL.String(),
			})

		c.Next()

		latency := time.Since(t).Milliseconds()
		status := c.Writer.Status()
		logkit.Info("response",
			logkit.Param{
				Key: "url",
				Val: c.Request.URL.String(),
			}, logkit.Param{
				Key: "latency",
				Val: latency,
			}, logkit.Param{
				Key: "status",
				Val: status,
			})
	}
}
