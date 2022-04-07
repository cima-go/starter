package starter

import (
	"context"
	"os"

	"github.com/pkg/errors"
)

func (s *Starter) run(signal <-chan os.Signal) (os.Signal, error) {
	proc := s.app(s.config)

	startCtx, cancel := context.WithTimeout(context.Background(), proc.StartTimeout())
	defer cancel()

	if err := proc.Start(startCtx); err != nil {
		return nil, errors.Errorf("ERROR\t\tFailed to start: %v", err)
	}

	done := <-signal

	stopCtx, cancel := context.WithTimeout(context.Background(), proc.StopTimeout())
	defer cancel()

	if err := proc.Stop(stopCtx); err != nil {
		return nil, errors.Errorf("ERROR\t\tFailed to stop cleanly: %v", err)
	}

	return done, nil
}
