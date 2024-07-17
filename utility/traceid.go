package utility

import (
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

const ginContextTraceId = "__gin_context_trace_id__"

var hostName = ""

func init() {
	var err error
	hostName, err = os.Hostname()
	if err != nil {
		hostName = ""
	}
}

func SetTraceId(c *gin.Context, id string) {
	c.Set(ginContextTraceId, id)
}

func GetTraceId(c *gin.Context) string {
	return c.GetString(ginContextTraceId)
}

func getFields(c *gin.Context, fields []zapcore.Field) []zapcore.Field {
	logId := c.GetString(ginContextTraceId)
	var method, path, raw, ip string
	if c.Request != nil {
		method = c.Request.Method
		path = c.Request.URL.Path
		raw = c.Request.URL.RawQuery
		ip = c.ClientIP()
	}
	if raw != "" {
		path = path + "?" + raw
	}

	fields = append(fields,
		zap.String("logId", logId),
		zap.String("hostName", hostName),
		zap.String("path", path),
		zap.String("method", method),
		zap.String("clientIp", ip),
		zap.String("dateTime", time.Now().Format(time.DateTime)),
	)

	return fields
}

type LogContext struct {
	Ctx    *gin.Context
	Msg    string
	Fields []zapcore.Field
}

func (l *LogContext) Formatter() (string, []zapcore.Field) {
	return l.Msg, getFields(l.Ctx, l.Fields)
}

func GetLogContext(c *gin.Context, msg string, fields ...zapcore.Field) *LogContext {
	return &LogContext{
		Ctx:    c,
		Msg:    msg,
		Fields: fields,
	}
}
