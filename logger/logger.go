package logger

import (
	"github.com/nik/platform-image-service/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
)

var logger *zap.Logger
var once sync.Once

//GetInstance returns instance of logger that is a singleton
func GetInstance(loggerConfig *config.Logger) *zap.Logger {
	once.Do(func() {
		logger = initLogger(loggerConfig)
	})
	return logger
}

func initLogger(logger *config.Logger) *zap.Logger {
	w := zapcore.AddSync(getCoreLogger(logger))
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)

	Logger := zap.New(core)

	return Logger
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
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
