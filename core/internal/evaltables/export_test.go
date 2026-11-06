package evaltables

import (
	"testing"
	"time"
)

// SetRetryWaitMinForTest sets the minimum wait before retrying a CES
// request until the test ends.
func SetRetryWaitMinForTest(t *testing.T, d time.Duration) {
	t.Helper()
	original := retryWaitMin
	retryWaitMin = d
	t.Cleanup(func() { retryWaitMin = original })
}
