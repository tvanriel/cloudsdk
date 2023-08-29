# CloudSDK

Provides a small standard set of tools for long running services. Take a look at the examples directory how to kickstart an application.

```go
import (
        "go.uber.org/fx"
        "github.com/tvanriel/cloudsdk/http"
        "github.com/tvanriel/cloudsdk/logging"
        "github.com/tvanriel/cloudsdk/mysql"
        "github.com/tvanriel/cloudsdk/kubernetes"
        "github.com/tvanriel/cloudsdk/s3"
        "github.com/tvanriel/cloudsdk/amqp"
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
                        ViperConfiguration,
                        KubernetesConfiguration,
                        amqpConfiguration,
                        HttpConfiguration,
                        S3Configuration,
                        MySQLConfiguration,
                        LoggingConfiguration,

                        http.AsRouteGroup(NewHttpController),
                ),
        ).Run()
}

```
