package rate_limiter

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const windowTime = time.Second * 10

func TestNewRateLimiterSuite(t *testing.T) {
	t.Run("test rate limiter with exact limit", testRateLimiter_ExactLimit)
	t.Run("test rate limiter remove expired requests", testRateLimit_RemoveExpiredRequests)
	t.Run("test rate limiter remove expired request on add", testRateLimit_RemoveExpiredRequestOnAdd)
	t.Run("test rate limiter concurrent allow request", testRateLimit_ConcurrentAllowRequest)
	t.Run("test rate limiter time spread", testRateLimiter_TimeSpread)
	t.Run("test rate limiter background cleanup", testRateLimiter_BackgroundCleanup)
	t.Run("test rate limiter with different users", testRateLimiter_WithDifferentUsers)
}

func testRateLimiter_ExactLimit(t *testing.T) {
	limiter := NewRateLimiter(3, windowTime)

	assert.Equal(t, true, limiter.AllowRequest("user1"))
	assert.Equal(t, true, limiter.AllowRequest("user1"))
	assert.Equal(t, true, limiter.AllowRequest("user1"))
	assert.Equal(t, false, limiter.AllowRequest("user1"))
}

func testRateLimiter_WithDifferentUsers(t *testing.T) {
	limiter := NewRateLimiter(1, windowTime)

	assert.Equal(t, true, limiter.AllowRequest("user1"))
	assert.Equal(t, false, limiter.AllowRequest("user1"))
	assert.Equal(t, true, limiter.AllowRequest("user2"))
	assert.Equal(t, false, limiter.AllowRequest("user2"))
}

func testRateLimit_RemoveExpiredRequests(t *testing.T) {
	limiter := NewRateLimiter(1, windowTime)

	expiredTime := time.Now().Add(-windowTime)

	limiter.RequestMap.Store("user1", RequestTimeStamps{expiredTime})

	limiter.removeExpiredRequests("user1")

	value, ok := limiter.RequestMap.Load("user1")

	assert.Equal(t, true, ok)
	assert.Equal(t, 0, len(value.(RequestTimeStamps)))
}

func testRateLimit_RemoveExpiredRequestOnAdd(t *testing.T) {
	limiter := NewRateLimiter(1, windowTime)

	limiter.RequestMap.Store("user1", RequestTimeStamps{time.Now().Add(-windowTime)})

	assert.Equal(t, limiter.AllowRequest("user1"), true)

	value, ok := limiter.RequestMap.Load("user1")
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(value.(RequestTimeStamps)))
}

func testRateLimit_ConcurrentAllowRequest(t *testing.T) {
	concurrentUser := 5
	numRequest := 5

	limiter := NewRateLimiter(numRequest, windowTime)

	wg := sync.WaitGroup{}

	for i := 0; i < concurrentUser; i++ {
		wg.Add(1)

		go func(user string) {
			defer wg.Done()

			for j := 0; j < numRequest; j++ {
				add := limiter.AllowRequest(user)
				assert.Equal(t, true, add)
			}
		}(fmt.Sprintf("user%d", i))
	}

	wg.Wait()

	limiter.RequestMap.Range(func(key, value interface{}) bool {
		assert.Equal(t, numRequest, len(value.(RequestTimeStamps)))
		return true
	})
}

func testRateLimiter_TimeSpread(t *testing.T) {
	limiter := NewRateLimiter(2, time.Millisecond*10)

	assert.Equal(t, true, limiter.AllowRequest("user1"))
	time.Sleep(time.Millisecond * 5)
	assert.Equal(t, true, limiter.AllowRequest("user1"))
	time.Sleep(time.Millisecond * 5)
	assert.Equal(t, true, limiter.AllowRequest("user1"))  //first request is expired and removed
	assert.Equal(t, false, limiter.AllowRequest("user1")) //second request is not expired so should fail
}

func testRateLimiter_BackgroundCleanup(t *testing.T) {
	limiter := NewRateLimiter(3, windowTime)

	limiter.RequestMap.Store("user1", RequestTimeStamps{time.Now().Add(-windowTime)})
	limiter.RequestMap.Store("user2", RequestTimeStamps{time.Now().Add(-windowTime)})

	// Simulate background cleanup
	go limiter.removeExpiredRequests("user1")
	go limiter.removeExpiredRequests("user2")

	time.Sleep(time.Millisecond * 100) // Allow cleanup to complete
}
