package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/twpayne/chezmoi/v2/internal/progress"
)

func run() error {
	i := 0
	cancel := progress.PeriodicWriteIfChanged(os.Stderr, 100*time.Millisecond, func() ([]byte, error) {
		i++
		if i == 25 {
			return nil, errors.New("stopped")
		}
		return []byte("\r" + time.Now().Round(250*time.Millisecond).String() + " "), nil
	})
	time.Sleep(3 * time.Second)
	return cancel()
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
