package middleware

import (
	"eyes/internal/web/common"
	jwt "eyes/utility"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set(common.CtxUserIDKey, "007")
		c.Next()
	}
}

func HHJWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		headerJwt := c.Request.Header.Get("jwt")
		if headerJwt == "" {
			if gin.Mode() == gin.DebugMode {
				zap.L().Error("未能从header中获取到[jwt]")
				c.Next()
			} else {
				common.RespError(c, common.CodeNoJWT)
				c.Abort()
			}
			return
		}

		mc, err := jwt.HHParseToken(headerJwt)
		if err != nil {
			common.RespError(c, common.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(common.CtxUserIDKey, mc.UserID)
		c.Next() // 后续的处理请求的函数中 可以用过c.Get(CtxUserIDKey) 来获取当前请求的用户信息
	}
}
