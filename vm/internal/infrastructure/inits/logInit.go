package inits

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

func LogInit() {
	// 读取配置文件
	type LogConfig struct {
		Level   int8   `yaml:"level"`
		LogPath string `yaml:"logPath"`
		AppName string `yaml:"appName"`
	}
	var config LogConfig
	if err := viper.UnmarshalKey("log", &config); err != nil {
		log.Panicf("无法解码为结构: %s", err)
	}
	level := config.Level
	logPath := config.LogPath
	appName := config.AppName
	
	writeSyncer := GetLogWriter(logPath, appName)
	encoder := GetEncoder()
	
	// 将日志输出到控制台
	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.Level(level))
	// 将日志输出到文件
	fileCore := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(level))
	
	// 合并控制台输出和文件输出
	core := zapcore.NewTee(consoleCore, fileCore)
	// 只输出到文件
	//core := zapcore.NewTee(fileCore)
	
	// 构建logger
	logger := zap.New(core, zap.AddCaller())
	
	// 替换全局zap
	zap.ReplaceGlobals(logger)
	// 替换全局log
	log.SetOutput(zap.NewStdLog(logger).Writer())
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
