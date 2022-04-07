package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/cima-go/starter"
)

func main() {
	s := starter.NewStarter(func(conf interface{}) *fx.App {
		return fx.New(
			fx.Provide(
				func() *Config { return conf.(*Config) },
			),
			fx.Invoke(
				echo,
			),
		)
	}, func(vip *viper.Viper) (interface{}, error) {
		cfg := &Config{}
		if err := vip.Unmarshal(cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	})

	root := &cobra.Command{Use: "test"}
	root.AddCommand(s.StdCmds()...)
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}

func echo(cfg *Config) {
	log.Printf("cfg name = %s\n", cfg.Name)
}
