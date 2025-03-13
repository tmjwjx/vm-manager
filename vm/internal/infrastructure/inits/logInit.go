package inits

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
	"vm/internal/infrastructure/globals"
)

func LogInit() {
	type LogConfig struct {
		Level   string `yaml:"level"`
		LogPath string `yaml:"logPath"`
		AppName string `yaml:"appName"`
	}
	var config LogConfig
	if err := viper.UnmarshalKey("log", &config); err != nil {
		globals.Log.Panicf("无法解码为结构: %s", err)
	}

	//level := viper.GetString("level")
	logPath := config.LogPath
	appName := config.AppName

	writeSyncer := GetLogWriter(logPath, appName)
	encoder := GetEncoder()

	// 将日志输出到控制台
	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	// 将日志输出到文件
	fileCore := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	// 合并控制台输出和文件输出
	core := zapcore.NewTee(consoleCore, fileCore)
	//// 只输出到文件
	//core := zapcore.NewTee(fileCore)

	logger := zap.New(core, zap.AddCaller())

	// 配置默认的log
	log.SetOutput(zap.NewStdLog(logger).Writer())

	globals.Log = logger.Sugar()
}

func GetEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func GetLogWriter(logPath, appName string) zapcore.WriteSyncer {
	// 确保日志目录存在
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		fmt.Printf("failed to create log directory: %v\n", err)
		return nil
	}

	currentDate := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("./%s/%s-%s.log", logPath, appName, currentDate)
	file, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return zapcore.AddSync(file)
}
