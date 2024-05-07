package redis

type Configuration struct {
	Address       string `hcl:"address"`
	Password      string `hcl:"password"`
	DatabaseIndex int    `hcl:"database_index"`
}
