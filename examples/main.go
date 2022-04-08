package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/cima-go/starter"
)

func main() {
	s := starter.New(
		"starter",
		func(conf interface{}) *fx.App {
			return fx.New(
				fx.Provide(
					func() *Config { return conf.(*Config) },
				),
				fx.Invoke(
					func(cfg *Config) {
						log.Printf("cfg name = %s\n", cfg.Name)
					},
				),
			)
		},
		func(vip *viper.Viper) (interface{}, error) {
			conf := &Config{}
			if err := vip.Unmarshal(conf); err != nil {
				return nil, err
			}
			return conf, nil
		},
		starter.WithConfigSearch("conf", "$HOME", "."),
	)

	root := &cobra.Command{Use: "starter"}
	root.AddCommand(starter.StdCommands(s)...)
	_ = root.Execute()
}
