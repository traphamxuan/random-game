package orm

import (
	"context"
	"fmt"

	"github.com/traphamxuan/random-game/package/config"
	"github.com/traphamxuan/random-game/package/logger"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type Orm struct {
	Models
	log    *logger.Logger
	config *config.Configuration
	db     *gorm.DB
}

var _ servicemanager.IService = (*Orm)(nil)

func NewOrm(ctx context.Context, sm *servicemanager.ServiceManager) *Orm {
	return &Orm{
		log:    servicemanager.GetServiceOrPanic[*logger.Logger](sm, "Logger", "Orm"),
		config: servicemanager.GetServiceOrPanic[*config.Configuration](sm, "Configuration", "Orm"),
	}
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	DBName   string `env:"DB_NAME" envDefault:"postgres"`
	Username string `env:"DB_USERNAME" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}

func (o *Orm) Setup(c context.Context) error {
	o.log.Debug("Setting up database")
	var dbConfig DatabaseConfig
	if err := o.config.ParseConfig(&dbConfig); err != nil {
		return err
	}

	o.log.Debug("Connecting to database")
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: zapgorm2.New(o.log.Desugar()),
	})
	if err != nil {
		return err
	}
	o.db = db
	o.Models = NewModels(&BaseORM{
		log: o.log,
		db:  db,
	})
	return nil
}

func (o *Orm) Start(c context.Context) error {
	db, err := o.db.DB()
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

func (o *Orm) Stop(c context.Context) error {
	db, err := o.db.DB()
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

func (o *Orm) StartTransaction(c context.Context, trx_func func(*Orm) error) error {
	return o.db.WithContext(c).
		Transaction(func(tx *gorm.DB) error {
			orm := *o
			orm.db = tx
			orm.Models = NewModels(&BaseORM{
				log: o.log,
				db:  tx,
			})
			return trx_func(&orm)
		})
}
