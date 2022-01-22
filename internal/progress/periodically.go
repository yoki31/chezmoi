package progress

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

// PeriodicWrite calls f every d and writes the returned bytes to w.
func PeriodicWrite(w io.Writer, d time.Duration, f func() ([]byte, error)) func() error {
	doneCh := make(chan struct{})
	errCh := make(chan error)

	go func(doneCh <-chan struct{}) {
		timer := time.NewTimer(d)
		defer timer.Stop()
		for {
			select {
			case <-doneCh:
				errCh <- nil
				return
			case <-timer.C:
				data, err := f()
				if err != nil {
					errCh <- err
					return
				}
				if _, err := w.Write(data); err != nil {
					errCh <- err
					return
				}
				timer.Reset(d)
			}
		}
	}(doneCh)

	cancelFunc := func() error {
		close(doneCh)
		return <-errCh
	}
	return cancelFunc
}

func PeriodicWriteIfChanged(w io.Writer, d time.Duration, f func() ([]byte, error)) func() error {
	var lastData []byte
	return PeriodicWrite(w, d, func() ([]byte, error) {
		switch data, err := f(); {
		case err != nil:
			return nil, err
		case bytes.Equal(lastData, data):
			lastData = data
			return nil, nil
		default:
			lastData = data
			return data, nil
		}
	})
}

func PeriodicWriteStringerIfChanged(w io.Writer, d time.Duration, final string, stringer fmt.Stringer) func() {
	var anyString bool
	var lastString string
	cancel := PeriodicWrite(w, d, func() ([]byte, error) {
		switch s := stringer.String(); {
		case s == lastString:
			if s != "" {
				anyString = true
			}
			return nil, nil
		default:
			if s != "" {
				anyString = true
			}
			return []byte(s), nil
		}
	})
	return func() {
		_ = cancel()
		if anyString {
			_, _ = w.Write([]byte(final))
		}
	}
}
