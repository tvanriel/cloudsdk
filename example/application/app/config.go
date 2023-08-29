package app

import (
	"github.com/spf13/viper"
	"github.com/tvanriel/cloudsdk/amqp"
	"github.com/tvanriel/cloudsdk/http"
	"github.com/tvanriel/cloudsdk/kubernetes"
	"github.com/tvanriel/cloudsdk/logging"
	"github.com/tvanriel/cloudsdk/mysql"
	"github.com/tvanriel/cloudsdk/s3"
)

type Configuration struct {
	Http       http.Configuration       `mapstructure:"http"`
	MySQL      mysql.Configuration      `mapstructure:"mysql"`
	S3         s3.Configuration         `mapstructure:"s3"`
	AMQP       amqp.Configuration       `mapstructure:"amqp"`
	Kubernetes kubernetes.Configuration `mapstructure:"kubernetes"`
	Logging    logging.Configuration    `mapstructure:"log"`
}

func ViperConfiguration() (Configuration, error) {
	var config Configuration
	viper.AddConfigPath(".")
	viper.SetConfigName("application")
	viper.SetConfigType(".yaml")
	err := viper.Unmarshal(&config)
	if err != nil {
		print(err)
		panic(err)
	}

	viper.AutomaticEnv()
	return config, err
}

func KubernetesConfiguration(config Configuration) kubernetes.Configuration {
	return config.Kubernetes
}

func MySQLConfiguration(config Configuration) mysql.Configuration {
	return config.MySQL
}

func HttpConfiguration(config Configuration) http.Configuration {
	return config.Http
}

func S3Configuration(config Configuration) s3.Configuration {
	return config.S3
}

func AmqpConfiguration(config Configuration) amqp.Configuration {
	return config.AMQP
}

func LoggingConfiguration(config Configuration) logging.Configuration {
	return config.Logging
}
