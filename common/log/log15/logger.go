// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log15

import (
	"fmt"
	"os"
	"time"
	"bytes"

	"go.uber.org/multierr"
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

func init() {
	// 日志输出等级
	lvl, _ := LvlFromString("debug")

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(lvl)

	encoderConfig := SetLc()
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.WriteSyncer(zapcore.AddSync(os.Stdout)), atomicLevel)

	DefaultLog = zap.New(core)
}

type ZapLogger struct {
	_log   *zap.Logger
	fields []Field
}

type Handler struct {}

func (l ZapLogger) SetHandler(h Handler) {}

func DiscardHandler() Handler {
	return Handler{}
}

type Logger interface {
	New(ctx ...interface{}) ZapLogger

	// Log a message at the given level with context key/value pairs
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Crit(msg string, ctx ...interface{})
}

type Format interface {
	Format(r *Record) []byte
}

func FormatFunc(f func(*Record) []byte) Format {
	return formatFunc(f)
}

type formatFunc func(*Record) []byte

func (f formatFunc) Format(r *Record) []byte {
	return f(r)
}

func LogfmtFormat() Format {
	return FormatFunc(func(r *Record) []byte {
		//common := []interface{}{r.KeyNames.Time, r.Time, r.KeyNames.Lvl, r.Lvl, r.KeyNames.Msg, r.Msg}
		buf := &bytes.Buffer{}
		return buf.Bytes()
	})
}

type Lazy struct {
	Fn interface{}
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
	case "error", "eror", "crit":
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
		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.WriteSyncer(zapcore.AddSync(os.Stdout)), atomicLevel)

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
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 大写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	return encoderConfig
}

func New(ctx ...interface{}) ZapLogger {
	return ZapLogger{DefaultLog, ParseFields(ctx)}
}

func Debug(msg string, ctx ...interface{}) {
	DefaultLog.Debug(msg, ParseFields(ctx)...)
}

func  Info(msg string, ctx ...interface{}) {
	DefaultLog.Info(msg, ParseFields(ctx)...)
}

func Warn(msg string, ctx ...interface{}) {
	DefaultLog.Warn(msg, ParseFields(ctx)...)
}

func  Error(msg string, ctx ...interface{}) {
	DefaultLog.Error(msg, ParseFields(ctx)...)
}

func  Crit(msg string, ctx ...interface{}) {
	DefaultLog.Error(msg, ParseFields(ctx)...)
}

func (l ZapLogger) New(ctx ...interface{}) ZapLogger {
	fields := AppendFields(l.fields, ParseFields(ctx))
	return ZapLogger{DefaultLog, fields}
}

func (l ZapLogger) Debug(msg string, ctx ...interface{}) {
	fields := AppendFields(l.fields, ParseFields(ctx))
	DefaultLog.Debug(msg, fields...)
}

func (l ZapLogger) Info(msg string, ctx ...interface{}) {
	fields := AppendFields(l.fields, ParseFields(ctx))
	DefaultLog.Info(msg, fields...)
}

func (l ZapLogger) Warn(msg string, ctx ...interface{}) {
	fields := AppendFields(l.fields, ParseFields(ctx))
	DefaultLog.Warn(msg, fields...)
}

func (l ZapLogger) Error(msg string, ctx ...interface{}) {
	fields := AppendFields(l.fields, ParseFields(ctx))
	DefaultLog.Error(msg, fields...)
}

func (l ZapLogger) Crit(msg string, ctx ...interface{}) {
	fields := AppendFields(l.fields, ParseFields(ctx))
	DefaultLog.Error(msg, fields...)
}



type Field = zapcore.Field

func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}

func AppendFields(prefix []Field, suffix []Field) []Field {
	newFields := make([]Field, len(prefix)+len(suffix))
	n := copy(newFields, prefix)
	copy(newFields[n:], suffix)
	return newFields
}

func ParseFields(args []interface{}) []Field {
	if len(args) == 0 {
		return nil
	}

	// Allocate enough space for the worst case; if users pass only structured
	// fields, we shouldn't penalize them with extra allocations.
	fields := make([]Field, 0, len(args))
	var invalid invalidPairs

	for i := 0; i < len(args); {
		// This is a strongly-typed field. Consume it and move on.
		if f, ok := args[i].(Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		// Make sure this element isn't a dangling key.
		if i == len(args)-1 {
			DefaultLog.DPanic("Ignored key without a value", zap.Any("ignored", args[i]))
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string, add this pair to the slice of invalid pairs.
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			// Subsequent errors are likely, so allocate once up front.
			if cap(invalid) == 0 {
				invalid = make(invalidPairs, 0, len(args)/2)
			}
			invalid = append(invalid, invalidPair{i, key, val})
		} else {
			fields = append(fields, Any(keyStr, val))
		}
		i += 2
	}

	// If we encountered any invalid key-value pairs, log an error.
	if len(invalid) > 0 {
		DefaultLog.DPanic("Ignored key-value pairs with non-string keys", zap.Array("invalid", invalid))
	}
	return fields
}

type invalidPair struct {
	position   int
	key, value interface{}
}

func (p invalidPair) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("position", int64(p.position))
	Any("key", p.key).AddTo(enc)
	Any("value", p.value).AddTo(enc)
	return nil
}

type invalidPairs []invalidPair

func (ps invalidPairs) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var err error
	for i := range ps {
		err = multierr.Append(err, enc.AppendObject(ps[i]))
	}
	return err
}