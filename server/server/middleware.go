package server

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/codecrafter404/bubble/config"
)

func CorsMiddleware(config *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if header := r.Header.Get("Origin"); header != "" {
			url, err := url.Parse(header)
			if err == nil {
				url.Path = ""
				url.User = nil
				url.Fragment = ""

				if slices.Contains(config.ServerConfig.CorsConfig.AccessControlAllowOrigin, url.String()) || slices.Contains(config.ServerConfig.CorsConfig.AccessControlAllowOrigin, "*") {
					w.Header().Add("Access-Control-Allow-Origin", url.String())
				}
			}
		}

		w.Header().Add("Access-Control-Allow-Credentials", fmt.Sprintf("%t", config.ServerConfig.CorsConfig.AccessControlAllowCredentials))

		w.Header().Add("Access-Control-Allow-Methods", strings.Join(config.ServerConfig.CorsConfig.AccessControlAllowMethods, ", "))

		if header := r.Header.Get("Access-Control-Request-Headers"); header != "" {
			if slices.Contains(config.ServerConfig.CorsConfig.AccessControlAllowHeaders, header) || slices.Contains(config.ServerConfig.CorsConfig.AccessControlAllowHeaders, "*") {
				w.Header().Add("Access-Control-Allow-Headers", header)
			}
		}

		next.ServeHTTP(w, r)
	})
}
