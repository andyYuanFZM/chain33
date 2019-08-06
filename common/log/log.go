// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log 日志相关接口以及函数
package log

import (
	"os"
	"github.com/33cn/chain33/types"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/33cn/chain33/common/log/log15"
)

func init() {
	//resetWithLogLevel("error")
}

//SetLogLevel 设置控制台日志输出级别
func SetLogLevel(logLevel string) {

	// 日志输出等级
	lvl, err := log15.LvlFromString(logLevel)
	if err != nil {
		// 日志级别配置不正确时默认为error级别
		lvl = zap.ErrorLevel
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(lvl)

	encoderConfig := log15.SetLc()
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.WriteSyncer(zapcore.AddSync(os.Stdout)),atomicLevel)

	logger := zap.New(core)
	log15.DefaultLog = logger
}

//SetFileLog 设置文件日志和控制台日志信息
func SetFileLog(log *types.Log) {
	if log == nil {
		log = &types.Log{LogFile: "logs/chain33.log"}
	}
	if log.LogFile == "" {
		SetLogLevel(log.LogConsoleLevel)
	} else {
		resetLog(log)
	}
}

// 清空原来所有的日志Handler，根据配置文件信息重置文件和控制台日志
func resetLog(log *types.Log) {
	fillDefaultValue(log)

	hook := lumberjack.Logger{
		Filename:   log.LogFile,
		MaxSize:    int(log.MaxFileSize),
		MaxBackups: int(log.MaxBackups),
		MaxAge:     int(log.MaxAge),
		LocalTime:  log.LocalTime,
		Compress:   log.Compress,
	}

	// 日志输出等级
	lvl, err := log15.LvlFromString(log.Loglevel)
	if err != nil {
		// 日志级别配置不正确时默认为error级别
		lvl = zap.ErrorLevel
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(lvl)

	encoderConfig := log15.SetLc()

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), atomicLevel)

	logger := zap.New(core)
	log15.DefaultLog = logger
}

// 保证默认性况下为error级别，防止打印太多日志
func fillDefaultValue(log *types.Log) {
	if log.Loglevel == "" {
		log.Loglevel = zapcore.ErrorLevel.String()
	}
	if log.LogConsoleLevel == "" {
		log.LogConsoleLevel = zapcore.ErrorLevel.String()
	}
}


