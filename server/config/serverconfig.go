package config

type ServerConfig struct {
	ServerPort        int
	PlaygroundEnabled bool
	CorsConfig        CorsConfig
}
