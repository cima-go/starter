package starter

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ServiceCommands(s *Starter) []*cobra.Command {
	install := &cobra.Command{
		Use:   "install",
		Short: fmt.Sprintf("install service of %s", s.AppName()),
		RunE:  s.cmdInstall,
	}
	uninstall := &cobra.Command{
		Use:   "uninstall",
		Short: fmt.Sprintf("uninstall service of %s", s.AppName()),
		RunE:  s.cmdUnInstall,
	}
	return []*cobra.Command{install, uninstall}
}

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

	return append(mixed, ServiceCommands(s)...)
}
