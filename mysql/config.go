package mysql

type Configuration struct {
	Host     string `hcl:"host"`
	User     string `hcl:"user"`
	Password string `hcl:"password"`
	DBName   string `hcl:"db_name"`
}
