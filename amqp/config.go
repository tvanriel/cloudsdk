package amqp

import "strings"

type Configuration struct {
	Address      string `hcl:"address"`
	ConsumerName string `hcl:"consumer_name"`
	TLS          bool   `hcl:"tls"`
	Username     string `hcl:"username"`
	Password     string `hcl:"password"`
}

func (config Configuration) Dsn() string {
	scheme := "amqp"
	if config.TLS {
		scheme = "amqps"
	}
	var sb strings.Builder
	sb.WriteString(scheme)
	sb.WriteString("://")
	if config.Username != "" {
		sb.WriteString(config.Username)
		if config.Password != "" {
			sb.WriteString(":")
			sb.WriteString(config.Password)
		}
		sb.WriteString("@")
	}
	sb.WriteString(config.Address)
	sb.WriteString("/")
	return sb.String()
}
