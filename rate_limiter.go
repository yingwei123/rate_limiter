package rate_limiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	RequestMap sync.Map
	Limit      int
	Duration   time.Duration
}

type RequestTimeStamps []time.Time

func NewRateLimiter(limit int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		RequestMap: sync.Map{},
		Limit:      limit,
		Duration:   duration,
	}
}

// AllowRequest checks if the request is allowed based on the rate limit and duration.
// It removes expired requests and checks if the request count is within the limit.
func (r *RateLimiter) AllowRequest(key string) bool {
	r.removeExpiredRequests(key)

	value, ok := r.RequestMap.Load(key)
	if !ok {
		r.RequestMap.Store(key, RequestTimeStamps{time.Now()})
		return true
	}

	timeStamps, ok := value.(RequestTimeStamps)
	if !ok {
		println("error: Type assertion failed for request timestamps for key: ", key)
		return false
	}

	if len(timeStamps) < r.Limit {
		timeStamps = append(timeStamps, time.Now())
		r.RequestMap.Store(key, timeStamps)
		return true
	}

	return false
}

// removeExpiredRequests removes expired requests from the request map.
func (r *RateLimiter) removeExpiredRequests(key string) {
	window := time.Now().Add(-r.Duration)

	value, ok := r.RequestMap.Load(key)
	if !ok {
		return
	}

	timeStamps, ok := value.(RequestTimeStamps)
	if !ok {
		println("error: Type assertion failed for request timestamps")
		return
	}

	if len(timeStamps) == 0 {
		return
	}

	validIndex := 0

	for _, ts := range timeStamps {
		if ts.After(window) {
			timeStamps[validIndex] = ts
			validIndex++
		}
	}

	r.RequestMap.Store(key, timeStamps[:validIndex])
}
