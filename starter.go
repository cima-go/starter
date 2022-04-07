package starter

import (
	"log"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type App func(conf interface{}) *fx.App

type Conf func(vip *viper.Viper) (interface{}, error)

type Starter struct {
	app    App
	conf   Conf
	config interface{}
}

func NewStarter(app App, conf Conf) *Starter {
	s := &Starter{
		app:  app,
		conf: conf,
	}
	return s
}

func (s *Starter) initC(cmd *cobra.Command, signals chan<- os.Signal) {
	var watcher func()
	if cmd.Flag("watch").Value.String() == "true" {
		watcher = func() {
			signals <- syscall.SIGHUP
		}
	}

	if conf2, err := s.configure(cmd.Flag("config").Value.String(), watcher); err != nil {
		log.Fatal(err)
	} else {
		s.config = conf2
	}
}

func (s *Starter) StdCmds() []*cobra.Command {
	serve := &cobra.Command{
		Use:   "serve",
		Short: "run server in frontend",
		RunE:  s.cmdServe,
	}
	start := &cobra.Command{
		Use:   "start",
		Short: "run server as daemon",
		RunE:  s.cmdStart,
	}
	stop := &cobra.Command{
		Use:   "stop",
		Short: "stop daemon server",
		RunE:  s.cmdStop,
	}
	reload := &cobra.Command{
		Use:   "reload",
		Short: "reload daemon server",
		RunE:  s.cmdReload,
	}

	server := []*cobra.Command{serve, start}
	daemon := []*cobra.Command{start, stop, reload}

	for _, cmd := range server {
		cmd.PersistentFlags().StringP("config", "c", "", "config path")
		cmd.PersistentFlags().BoolP("watch", "w", false, "live reload")
	}

	for _, cmd := range daemon {
		cmd.PersistentFlags().String("pid", "", "pid file")
		cmd.PersistentFlags().String("log", "", "log file")
	}

	return []*cobra.Command{
		serve,
		start,
		stop,
		reload,
	}
}
