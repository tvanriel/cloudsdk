package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
        registry *prometheus.Registry
        listenAddr string
}

func NewPrometheus(config Configuration) *Prometheus {

        registry := prometheus.NewRegistry()

        registry.MustRegister(
                collectors.NewGoCollector(),
                collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
        )

        return &Prometheus{
                registry: registry,
                listenAddr: config.Address,
        }
}


func Listen(prom *Prometheus) {
        srv := http.NewServeMux()
        srv.Handle("/metrics", promhttp.HandlerFor(prom.registry, promhttp.HandlerOpts{Registry: prom.registry}))
        go http.ListenAndServe(prom.listenAddr, srv)
}
