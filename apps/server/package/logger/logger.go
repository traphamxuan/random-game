package logger

import (
	"context"

	"github.com/traphamxuan/random-game/app/utils"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
	config *zap.Config
}

var _ servicemanager.IService = (*Logger)(nil)

func NewLogger(ctx context.Context) *Logger {
	var config zap.Config
	env := utils.GetAppMode(ctx)
	if env == "prod" {
		config = zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.WarnLevel)
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Level.SetLevel(zapcore.DebugLevel)
	}

	return &Logger{
		config: &config,
	}
}

func (l *Logger) Setup(c context.Context) error {
	logger, err := l.config.Build()
	if err != nil {
		return err
	}
	l.SugaredLogger = logger.Sugar()
	return nil
}

func (l *Logger) Start(c context.Context) error {
	l.Infof("Start logger at mode %s ", utils.GetAppMode(c))
	return nil
}

func (l *Logger) Stop(c context.Context) error {
	l.Info("End of logger. Flush all logs")
	l.Sync()
	return nil
}
