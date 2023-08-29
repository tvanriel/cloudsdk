package app

import (
	"github.com/tvanriel/cloudsdk/amqp"
	"github.com/tvanriel/cloudsdk/http"
	"github.com/tvanriel/cloudsdk/kubernetes"
	"github.com/tvanriel/cloudsdk/logging"
	"github.com/tvanriel/cloudsdk/mysql"
	"github.com/tvanriel/cloudsdk/s3"
	"go.uber.org/fx"
)

func Run() {
	fx.New(
		http.Module,
		mysql.Module,
		logging.Module,
		kubernetes.Module,
		s3.Module,
		amqp.Module,
		fx.Provide(
			http.AsRouteGroup(NewHttpController),
			ViperConfiguration,
			KubernetesConfiguration,
			AmqpConfiguration,
			HttpConfiguration,
			S3Configuration,
			MySQLConfiguration,
			LoggingConfiguration,
		),
	).Run()
}
