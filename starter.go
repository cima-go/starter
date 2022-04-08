package starter

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// App is builder for create fxApp instance
type App func(conf interface{}) *fx.App

// Conf is config decoder for provide actual type
type Conf func(vip *viper.Viper) (interface{}, error)

type Starter struct {
	name string
	app  App
	conf Conf

	flags     Flags
	confName  string
	confValue interface{}
}

func NewStarter(name string, app App, conf Conf, opts ...Option) *Starter {
	s := &Starter{
		name: name,
		app:  app,
		conf: conf,
	}

	if len(opts) > 0 {
		s.SetOptions(opts...)
	}

	return s
}

func (s *Starter) AppName() string {
	return s.name
}

func (s *Starter) SetOptions(opts ...Option) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *Starter) getFlags() Flags {
	if s.flags == nil || len(s.flags.Flags()) == 0 {
		s.flags = NewDefaultFlags()
	}
	return s.flags
}

func (s *Starter) initC(confFile string, watchConf bool, signals chan<- os.Signal) error {
	var notify func(c interface{})
	if watchConf {
		notify = func(conf interface{}) {
			s.confValue = conf
			signals <- syscall.SIGHUP
		}
	}

	var err error
	s.confValue, err = s.configure(confFile, notify)

	return err
}

func (s *Starter) run(signal <-chan os.Signal) (os.Signal, error) {
	app := s.app(s.confValue)

	startCtx, cancel1 := context.WithTimeout(context.Background(), app.StartTimeout())
	defer cancel1()

	if err := app.Start(startCtx); err != nil {
		return nil, errors.New(fmt.Sprintf("ERROR\t\tFailed to start: %v", err))
	}

	done := <-signal

	stopCtx, cancel2 := context.WithTimeout(context.Background(), app.StopTimeout())
	defer cancel2()

	if err := app.Stop(stopCtx); err != nil {
		return nil, errors.New(fmt.Sprintf("ERROR\t\tFailed to stop cleanly: %v", err))
	}

	return done, nil
}

func (s *Starter) tryRun(cmd *cobra.Command, signals chan os.Signal) (os.Signal, error) {
	if s.confValue == nil {
		flags := s.getFlags()
		if err := s.initC(flags.Config(cmd), flags.Watch(cmd), signals); err != nil {
			return nil, err
		}
	}

	return s.run(signals)
}
