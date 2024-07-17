package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {
	// 初始化 Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "YOUR_SENTRY_DSN",
	})
	if err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	// defer 关闭 Sentry
	defer sentry.Flush(2 * time.Second)

	// 初始化 Zap logger
	logger := initLogger()

	// 创建一个 Gin 路由
	r := gin.Default()

	// 设置一个中间件，用于将请求日志发送到 Sentry
	r.Use(func(c *gin.Context) {
		hub := sentry.GetHubFromContext(c)
		hub.Scope().SetRequest(c.Request)
		c.Next()
	})

	// 设置一个路由处理程序
	r.GET("/", func(c *gin.Context) {
		logger.Info("Hello, Gin!")
		c.String(200, "Hello, Gin!")
	})

	// 启动 Gin 服务器
	r.Run(":8080")
}

func initLogger() *zap.Logger {
	// 配置 Zap logger
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{"logfile.log"}, // 输出到文件
		ErrorOutputPaths: []string{"error.log"},   // 错误日志输出到文件
	}

	// 创建 Zap logger
	logger, err := cfg.Build()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
	}

	return logger
}
