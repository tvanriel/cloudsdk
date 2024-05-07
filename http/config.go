package http

type TLSOptions struct {
	CertFile string `hcl:"cert_file"`
	KeyFile  string `hcl:"key_file"`
}

type Configuration struct {
	Address   string      `hcl:"address"`
	Ratelimit int         `hcl:"rate_limit"`
	Debug     bool        `hcl:"debug"`
	TLS       *TLSOptions `hcl:"tls,block"`
}
