package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigtoml"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/rs/zerolog"
)

type Config struct {
	OrderConfig  OrderConfig
	ServerConfig ServerConfig
	DbPath       string        `default:"file:bubbles.db?_foreign_keys=on"`
	LogLevel     zerolog.Level `default:"info"`
}

func LoadConfig() (Config, error) {
	var conf Config
	loader := aconfig.LoaderFor(&conf, aconfig.Config{
		EnvPrefix: "BUBBLE",
		SkipFlags: true,
		Files:     []string{"config.yaml", "config.toml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".toml": aconfigtoml.New(),
		},
	})
	err := loader.Load()
	return conf, err
}
