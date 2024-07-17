package utility

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type ZapGormLogger struct {
	zapLogger *zap.Logger
}

func NewZapGormLogger(zapLogger *zap.Logger) *ZapGormLogger {
	return &ZapGormLogger{
		zapLogger: zapLogger,
	}
}

func (z *ZapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return z // 默认使用同一个实例，可以根据需要创建新的实例
}

func (z *ZapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	z.zapLogger.Sugar().Infof(msg, data...)
}

func (z *ZapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	z.zapLogger.Sugar().Warnf(msg, data...)
}

func (z *ZapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	z.zapLogger.Sugar().Errorf(msg, data...)
}

func (z *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	z.zapLogger.Sugar().Infof("%s [%.3fms] [rows:%v] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
	if err != nil {
		z.zapLogger.Sugar().Errorf("%s [%.3fms] [rows:%v] %s %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql, err)
	}
}
