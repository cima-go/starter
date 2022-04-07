package starter

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func (s *Starter) cmdServe(cmd *cobra.Command, args []string) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	s.initC(cmd, signals)

	for {
		sig, err := s.run(signals)
		if err != nil {
			return err
		}

		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			return nil
		}
	}
}

func (s *Starter) cmdStart(cmd *cobra.Command, args []string) error {
	dmx := newDaemon(cmd.Flag("pid").Value.String(), cmd.Flag("log").Value.String())

	if pid := dmx.runs(); pid > 0 {
		return errors.New(fmt.Sprintf("process is running, pid = %d", pid))
	}

	if done, err := dmx.start(); err != nil {
		return err
	} else if done {
		return nil
	}

	defer dmx.free()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGTERM)

	s.initC(cmd, signals)

	for {
		sig, err := s.run(signals)
		if err != nil {
			return err
		}

		if sig == syscall.SIGTERM {
			return nil
		}
	}
}

func (s *Starter) cmdStop(cmd *cobra.Command, args []string) error {
	return newDaemon(cmd.Flag("pid").Value.String(), "").shutdown()
}

func (s *Starter) cmdReload(cmd *cobra.Command, args []string) error {
	return newDaemon(cmd.Flag("pid").Value.String(), "").reload()
}
