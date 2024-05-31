package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	eMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/traphamxuan/gobs"

	"game-random-api/internal/api/handler"
	"game-random-api/internal/api/middleware"
	"game-random-api/internal/api/validator"
	"game-random-api/package/config"
	"game-random-api/package/logger"
	v "game-random-api/package/validator"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type APIConfig struct {
	AllowOrigins     string `env:"ALLOW_ORIGINS" envDefault:"*"`
	IdleTimeout      int    `env:"IDLE_TIMEOUT" envDefault:"10"`
	MaxConcurrent    int    `env:"MAX_CONCURRENT" envDefault:"1000"`
	MaxReadFrameSize int    `env:"MAX_READ_FRAME_SIZE" envDefault:"1048576"`
	Port             int    `env:"PORT" envDefault:"8080"`
}
type Router struct {
	log        *logger.Logger
	config     *APIConfig
	v          *v.CustomValidator
	validator  *validator.Validator
	middleware *middleware.Middleware
	handler    *handler.Handler
	httpServer *http.Server
}

var _ gobs.IService = (*Router)(nil)

func (r *Router) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&config.Configuration{},
		&validator.Validator{},
		&middleware.Middleware{},
		&handler.Handler{},
		&v.CustomValidator{},
		&logger.Logger{},
	}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		config := dependencies[0].(*config.Configuration)
		var cfg APIConfig
		if err := config.ParseConfig(&cfg); err != nil {
			return err
		}
		r.config = &cfg
		r.validator = dependencies[1].(*validator.Validator)
		r.middleware = dependencies[2].(*middleware.Middleware)
		r.handler = dependencies[3].(*handler.Handler)
		r.v = dependencies[4].(*v.CustomValidator)
		r.log = dependencies[5].(*logger.Logger)
		return r.Setup(ctx)
	}
	onStart := func(ctx context.Context) error {
		return r.Start(ctx)
	}
	onStop := func(ctx context.Context) error {
		return r.Stop(ctx)
	}
	sb.OnSetup = &onSetup
	sb.OnStart = &onStart
	sb.OnStop = &onStop
	return nil
}

// Setup implements IAppInstance.
func (r *Router) Setup(ctx context.Context) error {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(r.middleware.Log.Handler())
	e.Use(eMiddleware.CORSWithConfig(eMiddleware.CORSConfig{
		AllowOrigins: strings.Split(r.config.AllowOrigins, ","),
		AllowHeaders: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
	}))

	e.Use(eMiddleware.Recover())

	e.HTTPErrorHandler = r.middleware.ErrorHandling.Handler()
	e.Validator = r.v
	r.handler.SetRoute(ctx, e.Group("/api/v1"))

	h2s := &http2.Server{
		MaxConcurrentStreams: uint32(r.config.MaxConcurrent),
		MaxReadFrameSize:     uint32(r.config.MaxReadFrameSize),
		IdleTimeout:          time.Duration(r.config.IdleTimeout) * time.Second,
	}
	r.httpServer = &http.Server{
		Addr:    ":" + strconv.Itoa(r.config.Port),
		Handler: h2c.NewHandler(e, h2s),
	}
	return nil
}

// Start implements IAppInstance.
func (r *Router) Start(ctx context.Context) error {
	r.log.Infof("API server is running on port %d", r.config.Port)
	if err := r.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop implements IAppInstance.
func (r *Router) Stop(ctx context.Context) error {
	return r.httpServer.Shutdown(ctx)
}
