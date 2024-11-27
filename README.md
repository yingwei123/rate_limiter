
# Rate Limiter Package

This package provides a lightweight implementation of a rate-limiting mechanism for Go applications. It limits the number of requests allowed per user (or key) within a specified time window(sliding window).

## Features

- Rate limiting based on the number of requests per key.
- Configurable request limit and time window duration.
- Efficient handling using `sync.Map` for thread-safe access in concurrent environments.
- Automatic cleanup of expired requests for each key.

---

## Installation

```bash
go get github.com/yingwei123/rate_limiter
```

---

## Usage

### Create a Rate Limiter

To create a new rate limiter:

```go
import (
    "time"
    "github.com/yingwei123/rate_limiter"
)

limiter := rate_limiter.NewRateLimiter(5, time.Minute) // Allow 5 requests per key within 1 minute
```

### Allow Requests

To check if a request is allowed for a specific key, use the `AllowRequest` method:

```go
if limiter.AllowRequest("user1") {
    fmt.Println("Request allowed")
} else {
    fmt.Println("Rate limit exceeded")
}
```

### Example

Here is a complete example demonstrating how to use the `rate_limiter` package:

```go
package main

import (
    "fmt"
    "time"
    "github.com/yingwei123/rate_limiter"
)

func main() {
    limiter := rate_limiter.NewRateLimiter(3, 10*time.Second) // 3 requests allowed per 10 seconds

    // Simulate requests
    fmt.Println("Request 1:", limiter.AllowRequest("user1")) // true
    fmt.Println("Request 2:", limiter.AllowRequest("user1")) // true
    fmt.Println("Request 3:", limiter.AllowRequest("user1")) // true
    fmt.Println("Request 4:", limiter.AllowRequest("user1")) // false (rate limit exceeded)

    // Wait for requests to expire
    time.Sleep(11 * time.Second)

    fmt.Println("Request 5:", limiter.AllowRequest("user1")) // true (after waiting)
}
```

---

## Methods

### `NewRateLimiter(limit int, duration time.Duration) *RateLimiter`
Creates a new rate limiter instance.

- `limit`: Maximum number of requests allowed per key.
- `duration`: Time window for rate limiting.

### `AllowRequest(key string) bool`
Checks if a request is allowed for the specified key.

- Returns `true` if the request is allowed, otherwise `false`.

### `removeExpiredRequests(key string)` (private)
Removes expired timestamps from the request history of the specified key.

---

## How It Works

1. **Initialization**: Create a `RateLimiter` with a specified limit and duration.
2. **Request Handling**: For each request:
   - Expired timestamps are removed.
   - If the total requests within the time window are below the limit, the request is allowed.
3. **Thread-Safety**: Uses `sync.Map` to handle concurrent access in a safe manner.

---

## Testing

### Unit Tests

To test the `rate_limiter` package, you can write unit tests using the Go testing framework:

```go
package rate_limiter_test

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/yingwei123/rate_limiter"
)

func TestRateLimiter(t *testing.T) {
    limiter := rate_limiter.NewRateLimiter(3, 10*time.Second)

    assert.True(t, limiter.AllowRequest("user1"))
    assert.True(t, limiter.AllowRequest("user1"))
    assert.True(t, limiter.AllowRequest("user1"))
    assert.False(t, limiter.AllowRequest("user1"))
}
```

Run tests using:
```bash
go test -v ./...
```

---

## License

This project is open-source and available under the [MIT License](LICENSE.md).

---

## Contributions

Feel free to submit issues or pull requests to enhance the package. Contributions are always welcome!