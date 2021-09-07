package utils

import (
	"errors"
	"time"
)

// WaitChannelWithTimeout read exactly number of nbWaitFor values out from channel then return
// if this cannot achieve within timeout duration, non-nil error returned
func WaitChannelWithTimeout(c <-chan bool, nbWaitFor int, d time.Duration) (err error) {
	if nbWaitFor <= 0 {
		return
	}
	ready := 0
	timeoutC := time.After(d)
	for ready < nbWaitFor {
		select {
		case <-c:
			ready++
		case <-timeoutC:
			err = errors.New("waiting channel times out")
			return
		}
	}
	return
}
