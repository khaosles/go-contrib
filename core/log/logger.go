package glog

/*
   @File: logger.go
   @Author: khaosles
   @Time: 2023/4/11 22:16
   @Desc:
*/

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/khaosles/giz/fileutil"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/khaosles/go-contrib/core/config"
)

const (
	ONE_DAY = time.Hour * 24
	APP     = "logging"
)

var prefix string

var Logger *zap.SugaredLogger

type logging struct {
	LevelConsole string        `mapstructure:"level-console" default:"debug" yaml:"level-console" json:"level-console"`
	LevelFile    string        `mapstructure:"level-file" default:"info" yaml:"level-file" json:"level-file"`
	Prefix       string        `mapstructure:"prefix" default:"" yaml:"prefix" json:"prefix"`
	Path         string        `mapstructure:"path" default:"" yaml:"path" json:"path"`
	MaxHistory   time.Duration `mapstructure:"max-history" default:"7" yaml:"max-history" json:"maxHistory"`
	LogInConsole bool          `mapstructure:"log-in-console" default:"true" yaml:"log-in-console" json:"logInConsole"`
	LogInFile    bool          `mapstructure:"log-in-file" default:"false" yaml:"log-in-file" json:"logInFile"`
	ShowLine     bool          `mapstructure:"show-line" default:"false" yaml:"show-line" json:"show-line"`
}

func init() {
	var logging *logging
	if err := config.Configuration(APP, logging); err != nil {
		log.Fatal(err)
	}

	prefix = logging.Prefix
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     encodeTime,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	})

	levelConsole := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= levelChoice(logging.LevelConsole)
	})

	levelFile := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= levelChoice(logging.LevelFile)
	})

	var cores []zapcore.Core
	if logging.LogInFile {
		path := logPath(logging.Path)
		hook, _ := rotatelogs.New(
			path+"/%Y-%m-%d.log",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(ONE_DAY*logging.MaxHistory),
			rotatelogs.WithRotationTime(ONE_DAY),
			rotatelogs.WithClock(rotatelogs.Local),
		)
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(hook), levelFile))
		if logging.LogInConsole {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), levelConsole))
		}
	} else {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), levelConsole))
	}
	cores = append(cores)
	core := zapcore.NewTee(cores...)
	l := zap.New(core)
	if logging.ShowLine {
		l = l.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}
	Logger = l.Sugar()
}

func logPath(path string) string {
	if !filepath.IsAbs(path) {
		wk, _ := os.Getwd()
		fileutil.Join(wk, path)
	}
	_ = fileutil.Mkdir(path)
	return path
}

func encodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(prefix + t.In(time.FixedZone("CTS", 8*3600)).Format("2006-01-02 15:04:05.000"))
}

func levelChoice(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Debugf(fmt string, args ...interface{}) {
	Logger.Debugf(fmt, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(fmt string, args ...interface{}) {
	Logger.Infof(fmt, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}
func Warnf(fmt string, args ...interface{}) {
	Logger.Warnf(fmt, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Errorf(fmt string, args ...interface{}) {
	Logger.Errorf(fmt, args...)
}

func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

func Panicf(fmt string, args ...interface{}) {
	Logger.Panicf(fmt, args...)
}
