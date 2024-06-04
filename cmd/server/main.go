package main

import (
	"net/http"

	"github.com/flavioamaral-dev/go-experts-desafio-rate-limiter/configs"
	"github.com/flavioamaral-dev/go-experts-desafio-rate-limiter/internal/infra"
	"github.com/flavioamaral-dev/go-experts-desafio-rate-limiter/internal/infra/middlewares"
	"github.com/flavioamaral-dev/go-experts-desafio-rate-limiter/internal/infra/storage"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	webserver := infra.NewWebServer(configs.WebServerPort)
	storage_adapter, err := storage.InitRedis(configs.RedisAddr)
	if err != nil {
		panic(err)
	}

	customTokens := populateCustomTokens()
	rateLimiterConfig := &middlewares.RateLimiterConfig{
		LimitByIP: &middlewares.RateLimiterRateConfig{
			MaxRequestsPerSecond:  configs.LimitByIPMaxRPS,
			BlockTimeMilliseconds: configs.LimitByIPBlockTimeMs,
		},
		LimitByToken: &middlewares.RateLimiterRateConfig{
			MaxRequestsPerSecond:  configs.LimitByTokenMaxRPS,
			BlockTimeMilliseconds: configs.LimitByTokenBlockTimeMs,
		},
		StorageAdapter: storage_adapter,
		CustomTokens:   &customTokens,
	}

	rateLimiter := middlewares.NewRateLimiter(rateLimiterConfig)
	webserver.Use(rateLimiter)

	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!"))
	}
	webserver.AddHandler("/", rootHandler, "GET")
	webserver.Start()
}

func populateCustomTokens() map[string]*middlewares.RateLimiterRateConfig {
	return map[string]*middlewares.RateLimiterRateConfig{
		"ABC": {
			MaxRequestsPerSecond:  20,
			BlockTimeMilliseconds: 3000,
		},
		"DEF": {
			MaxRequestsPerSecond:  20,
			BlockTimeMilliseconds: 3000,
		},
	}
}
