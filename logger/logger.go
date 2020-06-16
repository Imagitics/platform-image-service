package logger

import (
	"github.com/nik/platform-image-service/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger

func InitLogger(config *config.ConfigModel) {
	w := zapcore.AddSync(getCoreLogger(config.Logger))
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)

	Logger = logger
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter(filePath string) zapcore.WriteSyncer {
	file, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	return zapcore.AddSync(file)
}

func getCoreLogger(logger *config.Logger) *lumberjack.Logger {
	var lumlog = &lumberjack.Logger{
		Filename:   logger.LoggerFileName,
		MaxSize:    10, // megabytes
		MaxBackups: 3,  // number of log files
		MaxAge:     3,  // days
	}

	return lumlog
}
