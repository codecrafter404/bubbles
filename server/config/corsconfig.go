package config

type CorsConfig struct {
	AccessControlAllowOrigin      []string
	AccessControlAllowCredentials bool
	AccessControlAllowMethods     []string
	AccessControlAllowHeaders     []string
}
