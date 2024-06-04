package middlewares

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/flavioamaral-dev/go-experts-desafio-rate-limiter/internal/infra/storage"
)

type RateLimiterRateConfig struct {
	MaxRequestsPerSecond  int64 `json:"maxRequestsPerSecond"`
	BlockTimeMilliseconds int64 `json:"blockTimeMilliseconds"`
}

type RateLimiterConfig struct {
	LimitByIP      *RateLimiterRateConfig
	LimitByToken   *RateLimiterRateConfig
	StorageAdapter storage.Storage
	CustomTokens   *map[string]*RateLimiterRateConfig
}

func NewRateLimiter(config *RateLimiterConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return rateLimiter(config, next)
	}
}

func rateLimiter(config *RateLimiterConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var block *time.Time
		var err error

		token := r.Header.Get("API_KEY")
		if token != "" {
			var tokenConfig *RateLimiterRateConfig
			customTokenConfig, ok := (*config.CustomTokens)[token]
			if ok {
				tokenConfig = customTokenConfig
			} else {
				tokenConfig = config.LimitByToken
			}

			block, err = checkRateLimit(r.Context(), "TOKEN", token, config, tokenConfig)
		} else {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			block, err = checkRateLimit(r.Context(), "IP", host, config, config.LimitByIP)
		}

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
			return
		}

		if block != nil {
			w.WriteHeader(429)
			w.Write([]byte("You have reached the maximum number of requests or actions allowed within a certain time frame."))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func checkRateLimit(ctx context.Context, keyType string, key string, config *RateLimiterConfig, rateConfig *RateLimiterRateConfig) (*time.Time, error) {
	if key == "" {
		return nil, nil
	}

	block, err := config.StorageAdapter.GetBlock(ctx, keyType, key)
	if err != nil {
		return nil, err
	}

	if block == nil {
		success, count, err := config.StorageAdapter.AddAccess(ctx, keyType, key, rateConfig.MaxRequestsPerSecond)
		if err != nil {
			return nil, err
		}

		if success {
			fmt.Printf("Access count within this window: %d\n", count)
		} else {
			fmt.Println("Access Denied")
			block, err = config.StorageAdapter.AddBlock(ctx, keyType, key, rateConfig.BlockTimeMilliseconds)
			if err != nil {
				return nil, err
			}
		}
	}

	if block != nil {
		fmt.Println("Blocked for", rateConfig.BlockTimeMilliseconds/1000, "seconds")
		return block, nil
	}

	return nil, nil
}
