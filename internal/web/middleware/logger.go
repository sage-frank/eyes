package middleware

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"

	"eyes/internal/web/common"
	"eyes/setting"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var dateLayout = "2006-01-02 15:04:05.999"

// NewLogger 初始化lg
func NewLogger(cfg *setting.LogConfig, mode string) (*zap.Logger, error) {
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	errSyncer := getLogWriter(cfg.ErrFilename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	encoder := getEncoder()

	l := new(zapcore.Level)
	err := l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return nil, err
	}

	var core zapcore.Core
	var logger *zap.Logger

	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	if mode == "dev" {
		// 进入开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(MyDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
		logger = zap.New(core, zap.AddCaller())

	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, errSyncer, highPriority),
			zapcore.NewCore(encoder, writeSyncer, lowPriority),
		)
		logger = zap.New(core, zap.AddCaller())
	}

	// sugarLogger := logger.Sugar()
	// zap.ReplaceGlobals(LG)
	return logger, nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(dateLayout))
	}
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func MyDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format(dateLayout))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// getLogWriter 日志
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(common.CtxUserIDKey)
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		params := make(map[string]any)

		switch c.Request.Method {

		case http.MethodGet:
			for key, values := range c.Request.URL.Query() {
				if len(values) > 0 {
					params[key] = values[0]
				}
			}
		case http.MethodPost:
			contentType := c.ContentType()
			if contentType == "application/x-www-form-urlencoded" {
				_ = c.Request.ParseForm()
				for key, values := range c.Request.Form {
					if len(values) > 0 {
						params[key] = values[0]
					}
				}
			} else if contentType == "application/json" {

				body, err := io.ReadAll(c.Request.Body)
				if err != nil {
					logger.Error("c.BindJSON(&jsonParams)", zap.Error(err))
				} else {
					params["body"] = body
				}

				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			}
		}

		c.Set("_params", c.Params)
		c.Next()
		cost := time.Since(start)
		session := sessions.Default(c)
		sessionId := session.Get("sessionID")
		if sessionId == nil {
			sessionId = session.ID()
			session.Set("sessionID", sessionId)
		}

		logger.Info(path,
			zap.String("sessionID", sessionId.(string)),
			zap.String("x-request-id", c.GetHeader("X-Request-ID")),
			zap.String("userID", userID),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Any("params", params),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Any("errors", c.Errors.ByType(gin.ErrorTypePrivate).JSON()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				fields := []zap.Field{
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("x-request-id", c.GetHeader("X-Request-ID")),
				}

				if brokenPipe {
					logger.Error(c.Request.URL.Path, fields...)
					// If the connection is dead, we can't write a status to it.
					err = c.Error(err.(error)) // nolint: err check
					c.Abort()
					return
				}

				if stack {
					fields = append(fields, zap.String("stack", string(debug.Stack())))
				}

				logger.Error("[Recovery from panic]", fields...)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
