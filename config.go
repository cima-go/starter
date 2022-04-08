package starter

import (
	"errors"
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var ErrConfigNotProvided = errors.New("config not provided")

type ConfigInit interface {
	Prepare() error
}

func (s *Starter) Configure(confFile string, watchConf func(c interface{})) (interface{}, error) {
	vip := viper.New()

	if confFile != "" {
		vip.SetConfigFile(confFile)
	} else if len(s.confInits) > 0 {
		for _, init := range s.confInits {
			if err := init(vip); err != nil {
				return nil, fmt.Errorf("config init: %w", err)
			}
		}
	} else {
		return nil, ErrConfigNotProvided
	}

	vip.AutomaticEnv()

	if err := vip.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	decoder := func() (interface{}, error) {
		cfg, err := s.conf(vip)
		if err != nil {
			return nil, err
		}

		if ci, is := cfg.(ConfigInit); is {
			if err := ci.Prepare(); err != nil {
				return nil, fmt.Errorf("config prepare: %w", err)
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
				log.Printf("[starter] error while decoding conf %v\n", err)
			}
		})
	}

	return decoder()
}
