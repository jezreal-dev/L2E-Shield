package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/jezrealdev/l2e-shield/internal/cache"
	"github.com/redis/go-redis/v9"
)

func RateLimit(c *cache.Cache, rps float64, burst int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			key := fmt.Sprintf("ratelimit:%s", ip)
			now := time.Now().UnixNano()
			clearBefore := now - int64(time.Second)

			pipe := c.Client().TxPipeline()
			pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(clearBefore, 10))
			cardCmd := pipe.ZCard(ctx, key)
			pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
			pipe.Expire(ctx, key, time.Second)

			_, err = pipe.Exec(ctx)
			if err != nil {
				http.Error(w, "internal server error during rate limit evaluation", http.StatusInternalServerError)
				return
			}

			if int(cardCmd.Val()) >= burst {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
