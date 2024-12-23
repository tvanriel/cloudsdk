package mysql

import (
	"strings"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	drivermysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type GormParams struct {
	fx.In

	Config Configuration
	Log    *zap.Logger
}

func NewGorm(in GormParams, lc fx.Lifecycle) (*gorm.DB, error) {
	db, err := gorm.Open(
		drivermysql.Open(
			makeDSN(in.Config.Host, in.Config.User, in.Config.Password, in.Config.DBName),
		),
		&gorm.Config{Logger: zapgorm2.New(in.Log.Named("mysql"))},
	)
	if err != nil {
		return db, err
	}

	lc.Append(fx.StartHook(func() {
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			for range ticker.C {
				if db == nil {
					return
				}
				sqlDB, err := db.DB()
				if err != nil {
					return
				}
				err = sqlDB.Ping()
				if err != nil {
					in.Log.Fatal("DB Connection is closed", zap.Error(err))
					return
				}
			}
		}()
	}))
	return db, err
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
