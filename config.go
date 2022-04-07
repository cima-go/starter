package starter

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type ConfigInit interface {
	Prepare() error
}

func (s *Starter) configure(confFile string, watchConf func()) (interface{}, error) {
	vip := viper.New()

	if confFile != "" {
		vip.SetConfigFile(confFile)
		goto INIT
	}

	if home, err := homedir.Dir(); err != nil {
		return nil, err
	} else {
		vip.AddConfigPath(home)
		vip.SetConfigName("config.yaml")
	}

INIT:

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
		vip.OnConfigChange(func(e fsnotify.Event) {
			if cfg, err := decoder(); err == nil {
				s.config = cfg
				watchConf()
			} else {
				log.Printf("error while decoding conf %v\n", err)
			}
		})
	}

	return decoder()
}
