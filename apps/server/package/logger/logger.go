package logger

import (
	"context"

	"game-random-api/utils"

	"github.com/traphamxuan/gobs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

var _ gobs.IService = (*Logger)(nil)

// Init implements gobs.IService.
func (l *Logger) Init(ctx context.Context, sb *gobs.Component) error {
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		return l.Setup(ctx)
	}
	onStart := func(ctx context.Context) error {
		return l.Start(ctx)
	}
	onStop := func(ctx context.Context) error {
		return l.Stop(ctx)
	}
	sb.OnSetup = &onSetup
	sb.OnStart = &onStart
	sb.OnStop = &onStop

	return nil
}

func (l *Logger) Setup(c context.Context) error {
	var config zap.Config
	env := utils.GetAppMode(c)
	if env == "prod" {
		config = zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.WarnLevel)
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Level.SetLevel(zapcore.DebugLevel)
	}

	logger, err := config.Build()
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
