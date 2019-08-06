// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log15

import (
	"fmt"
	"os"
	"time"

	"github.com/go-stack/stack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Lvl is a type for predefined log levels.
type Lvl int

var (
	DefaultLog *zap.Logger
)

// List of predefined log Levels
const (
	LvlCrit Lvl = iota
	LvlError
	LvlWarn
	LvlInfo
	LvlDebug
)

type ZapLogger struct {
	_log *zap.Logger
}

type Logger interface {
	New(ctx ...interface{}) Logger

	// Log a message at the given level with context key/value pairs
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Crit(msg string, ctx ...interface{})
}

// Returns the name of a Lvl
func (l Lvl) String() string {
	switch l {
	case LvlDebug:
		return "dbug"
	case LvlInfo:
		return "info"
	case LvlWarn:
		return "warn"
	case LvlError:
		return "eror"
	case LvlCrit:
		return "crit"
	default:
		panic("bad level")
	}
}

// LvlFromString returns the appropriate Lvl from a string name.
// Useful for parsing command line args and configuration files.
func LvlFromString(lvlString string) (zapcore.Level, error) {
	switch lvlString {
	case "debug", "dbug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error", "eror":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.DebugLevel, fmt.Errorf("Unknown level: %v", lvlString)
	}
}

// A Record is what a Logger asks its handler to write
type Record struct {
	Time     time.Time
	Lvl      Lvl
	Msg      string
	Ctx      []interface{}
	Call     stack.Call
	KeyNames RecordKeyNames
}

// RecordKeyNames are the predefined names of the log props used by the Logger interface.
type RecordKeyNames struct {
	Time string
	Msg  string
	Lvl  string
}

func getLogger() *zap.Logger {
	if DefaultLog == nil {
		// 日志输出等级
		lvl, _ := LvlFromString("debug")

		// 设置日志级别
		atomicLevel := zap.NewAtomicLevel()
		atomicLevel.SetLevel(lvl)

		encoderConfig := SetLc()
		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.WriteSyncer(zapcore.AddSync(os.Stdout)),atomicLevel)

		DefaultLog = zap.New(core)
	}
	return DefaultLog.WithOptions(zap.AddCallerSkip(1))
}


func SetLc() zapcore.EncoderConfig {
	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	return encoderConfig
}

func New(ctx ...interface{}) *ZapLogger {
	return &ZapLogger{getLogger().With(Any("module", ctx))}
}

func (l ZapLogger) Debug(msg string, ctx ...interface{}) {
	l._log.Debug(fmt.Sprint(msg, ctx))
}

func (l ZapLogger) Info(msg string, ctx ...interface{}) {
	l._log.Info(fmt.Sprint(msg, ctx))
}

func (l ZapLogger) Warn(msg string, ctx ...interface{}) {
	l._log.Warn(fmt.Sprint(msg, ctx))
}

func (l ZapLogger) Error(msg string, ctx ...interface{}) {
	l._log.Error(fmt.Sprint(msg, ctx))
}

func (l ZapLogger) Crit(msg string, ctx ...interface{}) {
	l._log.Error(fmt.Sprint(msg, ctx))
}

type Field = zapcore.Field

func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}