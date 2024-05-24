package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/traphamxuan/random-game/package/config"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
)

type Redis struct {
	config  *config.Configuration
	rClient *redis.Client
}

var _ servicemanager.IService = (*Redis)(nil)

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

func NewRedis(ctx context.Context, sm *servicemanager.ServiceManager) *Redis {
	return &Redis{
		config: servicemanager.GetServiceOrPanic[*config.Configuration](sm, "Configuration", "Redis"),
	}
}

func (o *Redis) Setup(c context.Context) error {
	var rdbCfg RedisConfig
	if err := o.config.ParseConfig(&rdbCfg); err != nil {
		return err
	}
	o.rClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rdbCfg.Host, rdbCfg.Port), // Redis server address
		Password: rdbCfg.Password,                                // No password by default
		DB:       rdbCfg.DB,                                      // Default DB
	})
	return nil
}

func (o *Redis) Start(c context.Context) error {
	return nil
}

func (o *Redis) Stop(c context.Context) error {
	return nil
}

func (r Redis) Set(c context.Context, key string, value interface{}, ttls time.Duration) error {
	rawData, err := json.Marshal(value)
	if nil != err {
		return err
	}

	return r.rClient.Set(c, key, rawData, ttls).Err()
}

func (r Redis) Get(c context.Context, key string, model interface{}) (interface{}, error) {
	rawData, err := r.rClient.Get(c, key).Result()
	if nil != err {
		return model, err
	}
	if len(rawData) == 0 {
		return model, errors.New("data empty")
	}
	err = json.Unmarshal([]byte(rawData), model)
	return model, err
}

func (r Redis) Wrap(
	c context.Context,
	key string,
	model interface{},
	callback func() (interface{}, error),
	ttls time.Duration,
) (result interface{}, err error) {
	result, err = r.Get(c, key, model)
	if nil != err {
		result, err = callback()
		if nil != err {
			return model, err
		}
		err = r.Set(c, key, result, ttls)
		if nil != err {
			return model, err
		}
	}
	return result, err
}
