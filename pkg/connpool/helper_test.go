package connpool

import (
	"testing"
	"time"
)

func Test_isExpiredTimeout(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		pastTime time.Duration
		timeout  time.Duration
		want     bool
	}{
		{
			name:     "Test timeout is expired",
			pastTime: time.Second * 7,
			timeout:  time.Second * 5,
			want:     true,
		},
		{
			name:     "Test timeout is not expired",
			pastTime: time.Second * 3,
			timeout:  time.Second * 5,
			want:     false,
		},
		{
			name:     "Test timeout and past time is equal",
			pastTime: time.Second * 5,
			timeout:  time.Second * 5,
			want:     false,
		},
		{
			name:     "Test timeout is null",
			pastTime: time.Second,
			timeout:  0,
			want:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isTimeoutExpired(tc.pastTime, tc.timeout); got != tc.want {
				t.Fatalf("got %t, want %t", got, tc.want)
			}
		})
	}
}
