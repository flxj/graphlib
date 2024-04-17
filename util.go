package graphlib

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrRunTimeout = errors.New("function run timeout")
)

func runWithTimeout(timeout time.Duration, f func() error) error {
	tr := time.NewTimer(timeout)
	defer tr.Stop()

	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- f()
	}()
	select {
	case <-tr.C:
		return ErrRunTimeout
	case err, ok := <-ch:
		if !ok {
			return nil
		}
		return err
	}
}

func runWithRetry(retry int, timeout time.Duration, f func() error) error {
	if retry <= 0 && timeout == time.Duration(0) {
		return f()
	} else if retry <= 0 {
		return runWithTimeout(timeout, f)
	} else if timeout == time.Duration(0) {
		var err error
		for i := 0; i <= retry; i++ {
			if err = f(); err == nil {
				return nil
			}
		}
		return fmt.Errorf("function runs exceeds the retry limit %d, %v", retry, err)
	} else {
		var err error
		for i := 0; i <= retry; i++ {
			if err = runWithTimeout(timeout, f); err == nil {
				return nil
			}
		}
		return fmt.Errorf("function runs exceeds the retry limit %d, %v", retry, err)
	}
}
