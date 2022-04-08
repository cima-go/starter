package starter

import (
	"github.com/spf13/cobra"
)

type CmdFlag struct {
	Name  string      // flag full name
	Short string      // shorthand letter
	Usage string      // usage description
	Value interface{} // string | bool | int
}

type Flags interface {
	Flags() []CmdFlag
	Config(cmd *cobra.Command) string
	Watch(cmd *cobra.Command) bool
	Pid(cmd *cobra.Command) string
	Log(cmd *cobra.Command) string
}

func AddFlags(cmd *cobra.Command, flags []CmdFlag) {
	for _, flag := range flags {
		switch val := flag.Value.(type) {
		case string:
			if flag.Short != "" {
				cmd.PersistentFlags().StringP(flag.Name, flag.Short, val, flag.Usage)
			} else {
				cmd.PersistentFlags().String(flag.Name, val, flag.Usage)
			}
			break
		case bool:
			if flag.Short != "" {
				cmd.PersistentFlags().BoolP(flag.Name, flag.Short, val, flag.Usage)
			} else {
				cmd.PersistentFlags().Bool(flag.Name, val, flag.Usage)
			}
			break
		}
	}
}

type DefaultFlags struct {
	config CmdFlag
	watch  CmdFlag
	pid    CmdFlag
	log    CmdFlag
}

func NewDefaultFlags() DefaultFlags {
	return DefaultFlags{
		config: CmdFlag{
			Name:  "config",
			Short: "c",
			Usage: "config path",
			Value: "",
		},
		watch: CmdFlag{
			Name:  "watch",
			Short: "w",
			Usage: "live reload",
			Value: false,
		},
		pid: CmdFlag{
			Name:  "pid",
			Usage: "pid file",
			Value: "",
		},
		log: CmdFlag{
			Name:  "log",
			Usage: "log file",
			Value: "",
		},
	}
}

func (f DefaultFlags) Flags() []CmdFlag {
	return []CmdFlag{f.config, f.watch, f.pid, f.log}
}

func (f DefaultFlags) Config(cmd *cobra.Command) string {
	return cmd.Flag(f.config.Name).Value.String()
}

func (f DefaultFlags) Watch(cmd *cobra.Command) bool {
	return cmd.Flag(f.watch.Name).Value.String() == "true"
}

func (f DefaultFlags) Pid(cmd *cobra.Command) string {
	return cmd.Flag(f.pid.Name).Value.String()
}

func (f DefaultFlags) Log(cmd *cobra.Command) string {
	return cmd.Flag(f.log.Name).Value.String()
}
