package config

type ServerConfig struct {
	ServerPort        int  `default:"8080"`
	PlaygroundEnabled bool `default:"true"`
	CorsConfig        CorsConfig
}
