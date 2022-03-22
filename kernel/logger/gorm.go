package logger

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	//"os"
	"time"
)

type GormLogrus struct {
	logger                logrus.Logger
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
}

func NewGormLogger() *GormLogrus {
	return &GormLogrus{
		logger:                *Logger.SQL,
		SkipErrRecordNotFound: false,
	}
}

func (l *GormLogrus) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.logger.SetLevel(logrus.Level(level))
	return l
}

func (l *GormLogrus) Info(ctx context.Context, s string, args ...any) {
	fields := logrus.Fields{
		"sql": s,
	}
	if len(args) > 0 {
		fields["latency"] = args[0]
	}
	l.logger.WithContext(ctx).Info(fields)
}

func (l *GormLogrus) Warn(ctx context.Context, s string, args ...any) {
	fields := logrus.Fields{
		"SQL": s,
	}
	if len(args) > 0 {
		fields["latency"] = args[0]
	}
	l.logger.WithContext(ctx).Warn(fields)
}

func (l *GormLogrus) Error(ctx context.Context, s string, args ...any) {
	fields := logrus.Fields{
		"SQL": s,
	}
	if len(args) > 0 {
		fields["latency"] = args[0]
	}
	l.logger.WithContext(ctx).Error(fields)
}

func (l *GormLogrus) Debug(ctx context.Context, s string, args ...any) {
	fields := logrus.Fields{
		"SQL": s,
	}
	if len(args) > 0 {
		fields["latency"] = args[0]
	}
	l.logger.WithContext(ctx).Debug(fields)
}

func (l *GormLogrus) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}

	traceLevel := logrus.InfoLevel

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		traceLevel = logrus.ErrorLevel
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		traceLevel = logrus.WarnLevel
	}

	switch traceLevel {
	case logrus.ErrorLevel:
		l.Error(ctx, sql, elapsed)
	case logrus.InfoLevel:
		l.Info(ctx, sql, elapsed)
	case logrus.WarnLevel:
		l.Warn(ctx, sql, elapsed)
	default:
		l.Debug(ctx, sql, elapsed)
	}
}

func (l *GormLogrus) SetSkipErrRecordNotFound(SkipErrRecordNotFound bool) gormLogger.Interface {
	l.SkipErrRecordNotFound = SkipErrRecordNotFound
	return l
}

func (l *GormLogrus) SetLogMode(level logrus.Level) gormLogger.Interface {
	l.logger.SetLevel(level)
	return l
}
