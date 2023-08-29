package mysql

import (
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"
	drivermysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)


type GormParams struct {
        fx.In
        
        Config Configuration
        Log *zap.Logger


}

func NewGorm(in GormParams) (*gorm.DB, error) {
        return gorm.Open(
                drivermysql.Open(
                        makeDSN(in.Config.Host, in.Config.User, in.Config.Password, in.Config.DBName),
                ),
                &gorm.Config{Logger: zapgorm2.New(in.Log)},
        )
}

func makeDSN(host string, username string, password string, dbName string) string {
        return strings.Join(
                []string{
                        username,
                        ":",
                        password,
                        "@tcp(",
                        host,
                        ")/",
                        dbName,
                        "?charset=utf8&parseTime=true&loc=Local",
                },
                "",
        )
}
