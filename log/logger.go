package log

import (
	"fmt"
	"os"
	"time"

	"github.com/tianhongw/misc-go/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ ILogger = new(zapLogger)

var logger *zapLogger

type ILogger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(format string, args ...interface{})

	Named(name string) ILogger
	Level() zapcore.Level
	SetLevel(l zapcore.Level)
	ZapLogger() *zap.Logger
	AddCallerSkip(skip int)
	Flush()
}

type zapLogger struct {
	base     *zap.Logger
	zapLevel zap.AtomicLevel
	children []*zapLogger
}

func Init(opts *conf.Options) error {
	var (
		zapLevel   zap.AtomicLevel
		stackLevel zapcore.Level
		zapEncoder zapcore.Encoder
		encoderCfg zapcore.EncoderConfig
	)

	if err := zapLevel.UnmarshalText([]byte(opts.Log.Level)); err != nil {
		return fmt.Errorf("failed to set log level: %v", err)
	}

	if opts.IsDevMode() {
		stackLevel = zap.WarnLevel
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		stackLevel = zap.ErrorLevel
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
	}

	switch opts.Log.Format {
	case "json":
		zapEncoder = zapcore.NewJSONEncoder(encoderCfg)
	default:
		zapEncoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return zapLevel.Enabled(lvl) && lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return zapLevel.Enabled(lvl) && lvl < zapcore.ErrorLevel
	})

	var cores []zapcore.Core

	for _, infoOut := range opts.Log.Output {
		cores = append(cores, zapcore.NewCore(zapEncoder, newWriteSyncer(&lumberjack.Logger{
			Filename:   infoOut,
			MaxSize:    opts.Log.MaxSize,
			MaxAge:     opts.Log.MaxAge,
			MaxBackups: opts.Log.MaxBackups,
			LocalTime:  true,
			Compress:   true,
		}), lowPriority))
	}

	for _, errOut := range opts.Log.ErrOutput {
		cores = append(cores, zapcore.NewCore(zapEncoder, newWriteSyncer(&lumberjack.Logger{
			Filename:   errOut,
			MaxSize:    opts.Log.MaxSize,
			MaxAge:     opts.Log.MaxAge,
			MaxBackups: opts.Log.MaxBackups,
			LocalTime:  true,
			Compress:   true,
		}), highPriority))
	}

	logger = &zapLogger{
		base: zap.New(
			zapcore.NewTee(cores...),
			zap.AddStacktrace(stackLevel),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				if opts.IsProdMode() {
					return zapcore.NewSampler(core, time.Second, 100, 100)
				}
				return core
			}),
		),
		zapLevel: zapLevel,
	}

	return nil
}

func newWriteSyncer(logger *lumberjack.Logger) (w zapcore.WriteSyncer) {
	switch logger.Filename {
	case "stdout":
		w = zapcore.Lock(os.Stdout)
	case "stderr":
		w = zapcore.Lock(os.Stderr)
	default:
		w = zapcore.AddSync(logger)
	}
	return
}

func (l *zapLogger) Named(name string) ILogger {
	child := &zapLogger{
		base:     l.base.Named(name),
		zapLevel: l.zapLevel,
		children: []*zapLogger{},
	}

	l.children = append(l.children, child)
	return child
}

func (l *zapLogger) AddCallerSkip(skip int) {
	l.base = l.base.WithOptions(zap.AddCallerSkip(skip))
}

func (l *zapLogger) SetLevel(v zapcore.Level) {
	l.zapLevel.SetLevel(zapcore.Level(v))
}

func (l *zapLogger) Level() zapcore.Level {
	return l.zapLevel.Level()
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.base.Debug(fmt.Sprint(args...))
}

func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.base.Debug(fmt.Sprintf(template, args...))
}

func (l *zapLogger) Info(args ...interface{}) {
	l.base.Info(fmt.Sprint(args...))
}

func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.base.Info(fmt.Sprintf(template, args...))
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.base.Warn(fmt.Sprint(args...))
}

func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.base.Warn(fmt.Sprintf(template, args...))
}

func (l *zapLogger) Error(args ...interface{}) {
	l.base.Error(fmt.Sprint(args...))
}

func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.base.Error(fmt.Sprintf(template, args...))
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.base.Fatal(fmt.Sprint(args...))
}

func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.base.Fatal(fmt.Sprintf(template, args...))
}

func (l *zapLogger) Panic(args ...interface{}) {
	l.base.Panic(fmt.Sprint(args...))
}

func (l *zapLogger) Panicf(template string, args ...interface{}) {
	l.base.Panic(fmt.Sprintf(template, args...))
}

func (l *zapLogger) DPanic(args ...interface{}) {
	l.base.DPanic(fmt.Sprint(args...))
}

func (l *zapLogger) DPanicf(template string, args ...interface{}) {
	l.base.DPanic(fmt.Sprintf(template, args...))
}

func (l *zapLogger) ZapLogger() *zap.Logger {
	return l.base
}

func (l *zapLogger) Flush() {
	_ = l.base.Sync()

	for _, c := range l.children {
		c.Flush()
	}
}

func Named(name string) ILogger {
	if logger == nil {
		panic("log is not inited.")
	}
	return logger.Named(name)
}

func Debug(args ...interface{}) {
	logger.base.Debug(fmt.Sprint(args...))
}

func Debugf(template string, args ...interface{}) {
	logger.base.Debug(fmt.Sprintf(template, args...))
}

func Info(args ...interface{}) {
	logger.base.Info(fmt.Sprint(args...))
}

func Infof(template string, args ...interface{}) {
	logger.base.Info(fmt.Sprintf(template, args...))
}

func Warn(args ...interface{}) {
	logger.base.Warn(fmt.Sprint(args...))
}

func Warnf(template string, args ...interface{}) {
	logger.base.Warn(fmt.Sprintf(template, args...))
}

func Error(args ...interface{}) {
	logger.base.Error(fmt.Sprint(args...))
}

func Errorf(template string, args ...interface{}) {
	logger.base.Error(fmt.Sprintf(template, args...))
}

func Fatal(args ...interface{}) {
	logger.base.Fatal(fmt.Sprint(args...))
}

func Fatalf(template string, args ...interface{}) {
	logger.base.Fatal(fmt.Sprintf(template, args...))
}

func Panic(args ...interface{}) {
	logger.base.Panic(fmt.Sprint(args...))
}

func Panicf(template string, args ...interface{}) {
	logger.base.Panic(fmt.Sprintf(template, args...))
}

func DPanic(args ...interface{}) {
	logger.base.DPanic(fmt.Sprint(args...))
}

func DPanicf(template string, args ...interface{}) {
	logger.base.DPanic(fmt.Sprintf(template, args...))
}

func ZapLogger() *zap.Logger {
	return logger.base
}

func Flush() {
	_ = logger.base.Sync()
}
