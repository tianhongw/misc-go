package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/tianhongw/misc-go/conf"
	"github.com/tianhongw/misc-go/log"
	"go.uber.org/zap/zapcore"
	xormcore "xorm.io/core"
)

const defaultMaxLifeTime = 5 * time.Minute

type Engine struct {
	*xorm.Engine
}

var engine *Engine

func Init(opts *conf.Options) error {
	dbParams := map[string]string{
		"charset":   opts.Database.Charset,
		"collation": opts.Database.Collation,
		"parseTime": strconv.FormatBool(opts.Database.ParseTime),
		"loc":       opts.Database.Loc,
	}

	dbSource := fmt.Sprintf(
		"%s:%s@%s(%s)%s",
		opts.Database.Username,
		opts.Database.Password,
		opts.Database.Network,
		opts.Database.Address,
		opts.Database.Name,
	)

	maxLifeTime := defaultMaxLifeTime
	if opts.Database.MaxLifetime != "" {
		if duration, err := time.ParseDuration(opts.Database.MaxLifetime); err == nil {
			maxLifeTime = duration
		}
	}

	_engine, err := xorm.NewEngineWithParams(opts.Database.Dialect, dbSource, dbParams)
	if err != nil {
		return err
	}

	logger := log.Named("db")
	_engine.SetLogger(newDBLogger(logger))
	_engine.ShowSQL(opts.IsDevMode())
	_engine.ShowExecTime(opts.IsDevMode())
	_engine.SetMaxIdleConns(opts.Database.MaxIdle)
	_engine.SetMaxOpenConns(opts.Database.MaxOpen)
	_engine.SetConnMaxLifetime(maxLifeTime)

	if err := _engine.Ping(); err != nil {
		return err
	}

	engine = &Engine{
		_engine,
	}

	return nil
}

func Close() {
	if engine == nil {
		return
	}

	engine.Close()
	engine = nil
}

type dbLogger struct {
	base    log.ILogger
	showSQL bool
}

func newDBLogger(logger log.ILogger) *dbLogger {
	l := logger.Named("raw")
	l.AddCallerSkip(1)

	return &dbLogger{
		base: l,
	}
}

func (l *dbLogger) Level() xormcore.LogLevel {
	lvl := l.base.Level()
	switch lvl {
	case zapcore.DebugLevel:
		return xormcore.LOG_DEBUG
	case zapcore.InfoLevel:
		return xormcore.LOG_INFO
	case zapcore.WarnLevel:
		return xormcore.LOG_WARNING
	case zapcore.ErrorLevel:
		return xormcore.LOG_ERR
	case zapcore.PanicLevel:
		fallthrough
	case zapcore.DPanicLevel:
		return xormcore.LOG_OFF
	default:
		return xormcore.LOG_UNKNOWN
	}
}

func (l *dbLogger) SetLevel(lvl xormcore.LogLevel) {
	switch lvl {
	case xormcore.LOG_DEBUG:
		l.base.SetLevel(zapcore.DebugLevel)
	case xormcore.LOG_INFO:
		l.base.SetLevel(zapcore.InfoLevel)
	case xormcore.LOG_WARNING:
		l.base.SetLevel(zapcore.WarnLevel)
	case xormcore.LOG_ERR:
		l.base.SetLevel(zapcore.ErrorLevel)
	case xormcore.LOG_OFF:
		l.base.SetLevel(zapcore.PanicLevel)
	}
}

func (l *dbLogger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		l.showSQL = show[0]
	} else {
		l.showSQL = true
	}
}

func (l *dbLogger) IsShowSQL() bool {
	return l.showSQL
}

func (l *dbLogger) Debug(v ...interface{}) {
	l.base.Debug(v...)
}

func (l *dbLogger) Debugf(format string, v ...interface{}) {
	l.base.Debugf(format, v...)
}

func (l *dbLogger) Error(v ...interface{}) {
	l.base.Error(v...)
}

func (l *dbLogger) Errorf(format string, v ...interface{}) {
	l.base.Errorf(format, v...)
}

func (l *dbLogger) Info(v ...interface{}) {
	l.base.Info(v...)
}

func (l *dbLogger) Infof(format string, v ...interface{}) {
	l.base.Infof(format, v...)
}

func (l *dbLogger) Warn(v ...interface{}) {
	l.base.Warn(v...)
}

func (l *dbLogger) Warnf(format string, v ...interface{}) {
	l.base.Warnf(format, v...)
}
