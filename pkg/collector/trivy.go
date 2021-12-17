package collector

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	dtypes "github.com/docker/docker/api/types"
	dclient "github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/thegeeklab/audit-exporter/pkg/client"
	"golang.org/x/xerrors"
)

const (
	namespace = "trivy"
)

type TrivyCollector struct {
	trivyClient     client.TrivyClient
	concurrency     int64
	vulnerabilities *prometheus.GaugeVec
}

func NewTrivyCollector(
	trivyClient client.TrivyClient,
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

func uniqueContainerImages(containers []dtypes.Container) []string {
	keys := make(map[string]bool)
	var images []string
	for _, container := range containers {
		image := container.Image
		if _, value := keys[image]; !value {
			keys[image] = true
			images = append(images, image)
		}
	}
	return images
}

func (c *TrivyCollector) Scan(ctx context.Context) error {
	if _, err := c.trivyClient.UpdateDatabase(ctx); err != nil {
		return xerrors.Errorf("failed to update database: %w", err)
	}

	cli, err := dclient.NewClientWithOpts(dclient.FromEnv, dclient.WithAPIVersionNegotiation())
	containers, err := cli.ContainerList(ctx, dtypes.ContainerListOptions{})
	if err != nil {
		return xerrors.Errorf("failed to get containers: %w", err)
	}

	semaphore := make(chan struct{}, c.concurrency)
	defer close(semaphore)

	wg := sync.WaitGroup{}
	mutex := &sync.Mutex{}

	var trivyResponses []client.TrivyResponse
	for _, image := range uniqueContainerImages(containers) {
		wg.Add(1)
		go func(image string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() {
				<-semaphore
			}()
			out, err := c.trivyClient.Do(ctx, image)
			if err != nil {
				// c.logger.Errorf("Failed to detect vulnerability at %s: %s\n", image, err.Error())
				return
			}

			var responses []client.TrivyResponse
			if err := json.Unmarshal(out, &responses); err != nil {
				// c.logger.Errorf("Failed to parse trivy response at %s: %s\n", image, err.Error())
				return
			}
			func() {
				mutex.Lock()
				defer mutex.Unlock()
				trivyResponses = append(trivyResponses, responses...)
			}()
		}(image)
	}
	wg.Wait()

	c.vulnerabilities.Reset()
	for _, trivyResponse := range trivyResponses {
		for _, vulnerability := range trivyResponse.Vulnerabilities {
			labels := []string{
				trivyResponse.ExtractImage(),
				vulnerability.VulnerabilityID,
				vulnerability.PkgName,
				vulnerability.InstalledVersion,
				vulnerability.Severity,
				vulnerability.FixedVersion,
			}
			c.vulnerabilities.WithLabelValues(labels...).Set(1)
		}
	}

	if _, err := c.trivyClient.ClearCache(ctx); err != nil {
		return xerrors.Errorf("failed to clear cache: %w", err)
	}

	return nil
}

func (c *TrivyCollector) StartLoop(ctx context.Context, interval time.Duration) {
	go func(ctx context.Context) {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				if err := c.Scan(ctx); err != nil {
					// c.logger.Errorf("Failed to scan: %s\n", err.Error())
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}

func (c *TrivyCollector) collectors() []prometheus.Collector {
	return []prometheus.Collector{
		c.vulnerabilities,
	}
}

func (c *TrivyCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range c.collectors() {
		collector.Describe(ch)
	}
}

func (c *TrivyCollector) Collect(ch chan<- prometheus.Metric) {
	for _, collector := range c.collectors() {
		collector.Collect(ch)
	}
}
