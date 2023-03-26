package mainimpl

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
)

func ThrottleHttp() (*throttled.HTTPRateLimiter, error) {
	// HTTP速率限制（一个轻量级的限流工具）
	throttleStore, err := memstore.New(65536) // 64k different combinations of limitByPathAndAddr
	if err != nil {
		return nil, fmt.Errorf("failed to init HTTP rate limiter store: %w", err)
	}
	quota := throttled.RateQuota{ //限流规则
		MaxRate:  throttled.PerSec(5), //一个path每秒的请求数量
		MaxBurst: 25,                  //最多支持的并发请求量
	}
	limiter, err := throttled.NewGCRARateLimiter(throttleStore, quota)
	if err != nil {
		return nil, fmt.Errorf("failed to init HTTP rate limiter: %w", err)
	}

	httpRateLimiter := &throttled.HTTPRateLimiter{
		RateLimiter: limiter,
		//VaryBy: &throttled.VaryBy{Path: true},
		VaryBy: limitByPathAndAddr{}, //根据path和addr进行限流
	}
	return httpRateLimiter, nil
}

type limitByPathAndAddr struct{}

func (limitByPathAndAddr) Key(r *http.Request) string {
	var k strings.Builder

	k.WriteString(r.URL.Path)
	k.WriteString("\n")

	remoteIP := r.Header.Get("X-Forwarded-For")
	if remoteIP == "" {
		remoteIP = r.RemoteAddr
	}
	k.WriteString(remoteIP)

	return k.String()
}
