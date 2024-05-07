package s3

type Configuration struct {
	Endpoint  string `hcl:"endpoint"`
	AccessKey string `hcl:"access_key"`
	SecretKey string `hcl:"secret_key"`
	SSL       bool   `hcl:"ssl"`
}
