package config

type CorsConfig struct {
	AccessControlAllowOrigin      []string `default:"*"`
	AccessControlAllowCredentials bool     `default:"true"`
	AccessControlAllowMethods     []string `default:"POST, GET, OPTIONS"`
	AccessControlAllowHeaders     []string `default:"*"`
}
