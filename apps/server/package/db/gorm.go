package db

import (
	"context"
	"fmt"

	"game-random-api/package/config"
	"game-random-api/package/logger"

	"github.com/traphamxuan/gobs"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	DBName   string `env:"DB_NAME" envDefault:"postgres"`
	Username string `env:"DB_USERNAME" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}
type Gorm struct {
	config *DatabaseConfig
	log    *logger.Logger
	Db     *gorm.DB
}

var _ gobs.IService = (*Gorm)(nil)

// Init implements gobs.IService.
func (o *Gorm) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&logger.Logger{},
		&config.Configuration{},
	}
	onSetup := func(ctx context.Context, deps []gobs.IService, extraDeps []gobs.CustomService) error {
		o.log = deps[0].(*logger.Logger)
		config := deps[1].(*config.Configuration)
		var dbConfig DatabaseConfig
		if err := config.ParseConfig(&dbConfig); err != nil {
			return err
		}
		o.config = &dbConfig
		return o.Setup(ctx)
	}
	onStart := func(ctx context.Context) error {
		return o.Start(ctx)
	}
	onStop := func(ctx context.Context) error {
		return o.Stop(ctx)
	}
	sb.OnSetup = &onSetup
	sb.OnStart = &onStart
	sb.OnStop = &onStop
	return nil
}

func (o *Gorm) Setup(ctx context.Context) error {
	o.log.Debug("Connecting to database")
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		o.config.Username,
		o.config.Password,
		o.config.Host,
		o.config.Port,
		o.config.DBName,
		o.config.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: zapgorm2.New(o.log.Desugar()),
	})
	if err != nil {
		return err
	}
	o.Db = db
	return nil
}

func (o *Gorm) Start(c context.Context) error {
	db, err := o.Db.DB()
	if err != nil {
		o.log.Error("Failed to get database connection", zap.Error(err))
		return err
	}
	err = db.Ping()
	if err != nil {
		o.log.Error("Database connection failed", zap.Error(err))
		return err
	}
	o.log.Info("Database connected")
	return nil
}

func (o *Gorm) Stop(c context.Context) error {
	db, err := o.Db.DB()
	if err != nil {
		o.log.Error("Failed to get database connection", zap.Error(err))
		return err
	}
	if err := db.Close(); err != nil {
		o.log.Error("Failed to close database connection", zap.Error(err))
		return err
	}
	o.log.Info("Database connection closed")
	return nil
}
