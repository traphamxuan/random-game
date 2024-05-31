package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"game-random-api/package/config"

	"github.com/redis/go-redis/v9"
	"github.com/traphamxuan/gobs"
)

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

type Redis struct {
	config  *RedisConfig
	rClient *redis.Client
}

var _ gobs.IService = (*Redis)(nil)

func (o *Redis) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&config.Configuration{},
	}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		config := dependencies[0].(*config.Configuration)
		var rdbCfg RedisConfig
		if err := config.ParseConfig(&rdbCfg); err != nil {
			return err
		}
		o.config = &rdbCfg
		return o.Setup(ctx)
	}
	sb.OnSetup = &onSetup
	return nil
}

func NewRedis(ctx context.Context) *Redis {
	return &Redis{}
}

func (o *Redis) Setup(c context.Context) error {

	o.rClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", o.config.Host, o.config.Port), // Redis server address
		Password: o.config.Password,                                  // No password by default
		DB:       o.config.DB,                                        // Default DB
	})
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
