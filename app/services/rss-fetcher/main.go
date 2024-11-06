package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/Zanda256/rss-tool/app/services/rss-fetcher/v1/handlers"
	"github.com/Zanda256/rss-tool/business/data/searchIndex/es"
	v1 "github.com/Zanda256/rss-tool/business/web/v1"
	"github.com/Zanda256/rss-tool/foundation/logger"
)

const (
	MaxIdleConnsEnvKey         = "MAX_IDLE_CONNS_PER_HOST"
	HttpClientReqTimeoutEnvKey = "CLIENT_REQ_TIMEOUT"
	APIPortEnvKey              = "API_PORT"
	APIHostEnvKey              = "API_HOST"

	EsUrlEnvKey      = "ES_URL"
	EsPasswordEnvKey = "ES_PASSWORD"
	EsUserEnvKey     = "ES_USER"
)

// if we are building locally, build will default to `develop`
// But in the docker container, a build reference will be passed
var build = "develop"

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* ERROR ENCOUNTERED ******")
		},
	}

	traceIDFunc := func(ctx context.Context) string {
		return "test run"
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "RSS-PERCEVAL", traceIDFunc, events)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		return
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "build", build)

	// -------------------------------------------------------------------------
	// Configuration
	cfg := struct {
		Web struct {
			API struct {
				ReadTimeout     time.Duration
				WriteTimeout    time.Duration
				IdleTimeout     time.Duration
				ShutdownTimeout time.Duration
				Host            string // `conf:"default:0.0.0.0:3000"`
				//	DebugHost       string       // `conf:"default:0.0.0.0:4000"`
				Port string
			}
			Client struct {
				MaxIdleConnsPerHost string
				Timeout             string
			}
		}
		ES struct {
			User     string
			Password string
			URL      string
		}
	}{}

	// Get web configs
	cfg.Web.API.Port = getEnvValue(APIPortEnvKey, "6000")
	cfg.Web.API.Host = getEnvValue(APIHostEnvKey, "0.0.0.0")
	cfg.Web.Client.MaxIdleConnsPerHost = getEnvValue(MaxIdleConnsEnvKey, "20")
	cfg.Web.Client.Timeout = getEnvValue(HttpClientReqTimeoutEnvKey, "10")

	mxConns, err := strconv.Atoi(cfg.Web.Client.MaxIdleConnsPerHost)
	if err != nil {
		return err
	}

	clientTimeout, err := strconv.Atoi(cfg.Web.Client.Timeout)
	if err != nil {
		return err
	}

	log.Info(ctx, "waste", "timeout", clientTimeout)

	var c = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: mxConns,
		},
		// Timeout specifies a time limit for requests made by this Client.
		// The timeout includes connection time, any redirects, and reading the response body.
		Timeout: 10 * time.Second,
	}

	// Get ES configs
	cfg.ES.URL = mustGet(EsUrlEnvKey)
	cfg.ES.Password = mustGet(EsPasswordEnvKey)
	cfg.ES.User = mustGet(EsUserEnvKey)

	// Bootstrap es db
	esClient, err := es.NewESClient(es.Config{
		URL:      cfg.ES.URL,
		User:     cfg.ES.User,
		Password: cfg.ES.Password,
	})
	if err != nil {
		return err
	}

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := v1.APIMuxConfig{
		WebClient: c,
		Build:     build,
		Shutdown:  shutdown,
		Log:       log,
		Db:        esClient,
	}

	apiMux := v1.APIMux(cfgMux, handlers.Routes{})

	api := http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Web.API.Host, cfg.Web.API.Port),
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.API.ReadTimeout,
		WriteTimeout: cfg.Web.API.WriteTimeout,
		IdleTimeout:  cfg.Web.API.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.API.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func getEnvValue(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func mustGet(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("env var %s is required", key))
	}
	return value
}
