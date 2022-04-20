package starter

import (
	"github.com/kardianos/osext"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Option func(s *Starter)

func WithConfigSearch(name string, dirs ...string) Option {
	return func(s *Starter) {
		s.confInits = append(s.confInits, func(vip *viper.Viper) error {
			vip.SetConfigName(name)
			for _, dir := range dirs {
				switch dir {
				case "$EXE":
					exe, err := osext.ExecutableFolder()
					if err != nil {
						return err
					}
					vip.AddConfigPath(exe)
				case "$HOME":
					home, err := homedir.Dir()
					if err != nil {
						return err
					}
					vip.AddConfigPath(home)
				default:
					vip.AddConfigPath(dir)
				}
			}
			return nil
		})
	}
}

func WithCustomConfigInit(init ConfInit) Option {
	return func(s *Starter) {
		s.confInits = []ConfInit{init}
	}
}

func WithCustomFlags(flags Flags) Option {
	return func(s *Starter) {
		s.flags = flags
	}
}
