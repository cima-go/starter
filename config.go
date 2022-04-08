package starter

import (
	"errors"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var ErrConfigNotProvided = errors.New("config not provided")

type ConfigInit interface {
	Prepare() error
}

func (s *Starter) configure(confFile string, watchConf func(c interface{})) (interface{}, error) {
	vip := viper.New()

	if confFile != "" {
		vip.SetConfigFile(confFile)
	} else if s.confName != "" {
		if home, err := homedir.Dir(); err != nil {
			return nil, err
		} else {
			vip.AddConfigPath(home)
			vip.SetConfigName(s.confName)
		}
	} else {
		return nil, ErrConfigNotProvided
	}

	vip.AutomaticEnv()

	if err := vip.ReadInConfig(); err != nil {
		return nil, err
	}

	decoder := func() (interface{}, error) {
		cfg, err := s.conf(vip)
		if err != nil {
			return nil, err
		}

		if ci, is := cfg.(ConfigInit); is {
			if err := ci.Prepare(); err != nil {
				return nil, err
			}
		}

		return cfg, nil
	}

	if watchConf != nil {
		vip.WatchConfig()
		vip.OnConfigChange(func(_ fsnotify.Event) {
			if conf, err := decoder(); err == nil {
				watchConf(conf)
			} else {
				log.Printf("error while decoding conf %v\n", err)
			}
		})
	}

	return decoder()
}
