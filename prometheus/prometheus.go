package prometheus

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Prometheus struct {
	registry   *prometheus.Registry
	listenAddr string
	Log        *zap.Logger
}

func NewPrometheus(config Configuration, log *zap.Logger, lc fx.Lifecycle) *Prometheus {

	registry := prometheus.NewRegistry()
	l := log.Named("prometheus")

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	prom := &Prometheus{
		registry:   registry,
		listenAddr: config.Address,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l.Info("listen address", zap.String("addr", prom.listenAddr))

			go func() {

				srv := http.NewServeMux()
				srv.Handle("/metrics", promhttp.HandlerFor(prom.registry, promhttp.HandlerOpts{Registry: prom.registry}))

				err := http.ListenAndServe(prom.listenAddr, srv)
				if err != nil {
					l.Error("Failed to open metricslistener", zap.Error(err))
				}
			}()
			return nil
		},
	})
	return prom

}

func (p *Prometheus) Register(c prometheus.Collector) error {
	return p.registry.Register(c)
}
func (p *Prometheus) Unregister(c prometheus.Collector) bool {
	return p.registry.Unregister(c)
}
func (p *Prometheus) Gather() ([]*io_prometheus_client.MetricFamily, error) {
	return p.registry.Gather()
}
func (p *Prometheus) Collect(ch chan<- prometheus.Metric) {
	p.registry.Collect(ch)
}
func (p *Prometheus) Describe(ch chan<- *prometheus.Desc) {
	p.registry.Describe(ch)
}
