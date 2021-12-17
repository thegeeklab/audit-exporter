package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "trivy"
)

type TrivyCollector struct {
	concurrency     int64
	vulnerabilities *prometheus.GaugeVec
}

func NewTrivyCollector(
	concurrency int64,
) *TrivyCollector {
	return &TrivyCollector{
		concurrency: concurrency,
		vulnerabilities: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "vulnerabilities",
			Help:      "Vulnerabilities detected by trivy",
		}, []string{"image", "vulnerabilityId", "pkgName", "installedVersion", "severity", "fixedVersion"}),
	}
}
