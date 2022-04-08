package starter

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Option func(s *Starter)

func WithConfigSearch(name string, dirs ...string) Option {
	return func(s *Starter) {
		s.confInits = append(s.confInits, func(vip *viper.Viper) error {
			vip.SetConfigName(name)
			for _, dir := range dirs {
				if dir == "$HOME" { // special dir
					home, err := homedir.Dir()
					if err != nil {
						return err
					}
					vip.AddConfigPath(home)
				} else {
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
