package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

type Prometheus struct {
	registry   *prometheus.Registry
	listenAddr string
}

func NewPrometheus(config Configuration) *Prometheus {

	registry := prometheus.NewRegistry()

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &Prometheus{
		registry:   registry,
		listenAddr: config.Address,
	}
}

func Listen(prom *Prometheus) {
	srv := http.NewServeMux()
	srv.Handle("/metrics", promhttp.HandlerFor(prom.registry, promhttp.HandlerOpts{Registry: prom.registry}))
	go http.ListenAndServe(prom.listenAddr, srv)
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
