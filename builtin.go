package starter

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kardianos/osext"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

func (s *Starter) cmdServe(cmd *cobra.Command, args []string) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	for {
		sig, err := s.run(cmd, signals)
		if err != nil {
			return err
		}

		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			return nil
		}
	}
}

func (s *Starter) cmdStart(cmd *cobra.Command, args []string) error {
	flags := s.getFlags()

	dmx := newDaemon(flags.Pid(cmd), flags.Log(cmd))
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

	for {
		sig, err := s.run(cmd, signals)
		if err != nil {
			return err
		}

		if sig == syscall.SIGTERM {
			return nil
		}
	}
}

func (s *Starter) cmdStop(cmd *cobra.Command, args []string) error {
	return newDaemon(s.getFlags().Pid(cmd), "").shutdown()
}

func (s *Starter) cmdReload(cmd *cobra.Command, args []string) error {
	return newDaemon(s.getFlags().Pid(cmd), "").reload()
}

func (s *Starter) cmdInstall(cmd *cobra.Command, args []string) error {
	wd, err := osext.ExecutableFolder()
	if err != nil {
		return fmt.Errorf("get executable folder: %s", err)
	}

	name := s.name
	if f := cmd.Flag("name"); f != nil {
		name = cmd.Flag("name").Value.String()
	}

	prg := &program{}
	svcConfig := &service.Config{
		Name:             name,
		DisplayName:      name,
		Description:      name,
		WorkingDirectory: wd,
	}

	if len(args) > 0 {
		svcConfig.Arguments = args
	} else {
		svcConfig.Arguments = []string{"serve"}
	}

	svc, err := service.New(prg, svcConfig)
	if err != nil {
		return fmt.Errorf("create service: %s", err)
	}

	if err := svc.Install(); err != nil {
		return fmt.Errorf("install service: %s", err)
	}

	if err := svc.Start(); err != nil {
		return fmt.Errorf("start service: %s", err)
	}

	return nil
}

func (s *Starter) cmdUnInstall(cmd *cobra.Command, args []string) error {
	name := s.name
	if f := cmd.Flag("name"); f != nil {
		name = cmd.Flag("name").Value.String()
	}

	prg := &program{}
	svcConfig := &service.Config{
		Name:        name,
		DisplayName: s.name,
		Description: s.name,
	}

	svc, err := service.New(prg, svcConfig)
	if err != nil {
		return fmt.Errorf("create service: %s", err)
	}

	if st, err := svc.Status(); err == nil && st == service.StatusRunning {
		if err := svc.Stop(); err != nil {
			log.Printf("stop service:%s", err)
		}
	}

	if err := svc.Uninstall(); err != nil {
		return fmt.Errorf("uninstall service:%s", err)
	}
	return nil
}
