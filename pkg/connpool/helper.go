package connpool

import "time"

// isTimeoutExpired functoin returns timeout is expired or not.
func isTimeoutExpired(pastTime, timeout time.Duration) bool {
	return pastTime > timeout
}
