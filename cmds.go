package starter

import (
	"fmt"

	"github.com/spf13/cobra"
)

func StdCommands(s *Starter) []*cobra.Command {
	serve := &cobra.Command{
		Use:   "serve",
		Short: fmt.Sprintf("run %s in frontend", s.AppName()),
		RunE:  s.cmdServe,
	}
	start := &cobra.Command{
		Use:   "start",
		Short: fmt.Sprintf("run %s as daemon", s.AppName()),
		RunE:  s.cmdStart,
	}
	stop := &cobra.Command{
		Use:   "stop",
		Short: fmt.Sprintf("stop daemon of %s", s.AppName()),
		RunE:  s.cmdStop,
	}
	reload := &cobra.Command{
		Use:   "reload",
		Short: fmt.Sprintf("reload daemon of %s", s.AppName()),
		RunE:  s.cmdReload,
	}

	mixed := []*cobra.Command{serve, start, stop, reload}

	flags := s.getFlags().Flags()
	for _, cmd := range mixed {
		AddFlags(cmd, flags)
	}

	return mixed
}
