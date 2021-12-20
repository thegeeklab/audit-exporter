package collector

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	dtypes "github.com/docker/docker/api/types"
	dclient "github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/thegeeklab/audit-exporter/pkg/client"
	"github.com/thegeeklab/audit-exporter/pkg/utils"
	"golang.org/x/xerrors"
)

const (
	namespace = "trivy"
	name      = "trivy"
)

// TrivyCollector defines the trivy collector instance
type TrivyCollector struct {
	trivyClient        client.TrivyClient
	settings           Settings
	Vulnerabilities    *prometheus.GaugeVec
	VulnerabilitiesSum *prometheus.GaugeVec
	logger             *logrus.Logger
}

// NewTrivyCollector creates a new collector instance
func NewTrivyCollector(
	trivyClient client.TrivyClient,
	settings Settings,
	logger *logrus.Logger,
) *TrivyCollector {
	return &TrivyCollector{
		logger:   logger,
		settings: settings,
		Vulnerabilities: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "vulnerabilities",
			Help:      "Vulnerabilities detected by trivy",
		}, []string{"artifactName", "artifactType", "vulnerabilityId", "pkgName", "installedVersion", "severity", "fixedVersion"}),
		VulnerabilitiesSum: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "vulnerabilities_sum",
			Help:      "Vulnerabilities detected by trivy",
		}, []string{"artifactName", "artifactType", "severity"}),
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

// Scan checks the discovered docker images in parallel and maps the results to prometheus metrics
func (c *TrivyCollector) Scan(ctx context.Context) error {
	if _, err := c.trivyClient.UpdateDatabase(ctx); err != nil {
		return xerrors.Errorf("failed to update database: %w", err)
	}

	cli, err := dclient.NewClientWithOpts(dclient.FromEnv, dclient.WithAPIVersionNegotiation())
	if err != nil {
		return xerrors.Errorf("failed to connect to docker daemon: %w", err)
	}
	containers, err := cli.ContainerList(ctx, dtypes.ContainerListOptions{})
	if err != nil {
		return xerrors.Errorf("failed to get containers: %w", err)
	}

	semaphore := make(chan struct{}, c.settings.ClientConcurrency)
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
				c.logger.Errorf("Failed to detect vulnerability at %s: %s", image, err.Error())
				return
			}

			var response client.TrivyResponse
			if err := json.Unmarshal([]byte(out), &response); err != nil {
				c.logger.Errorf("Failed to parse trivy response at %s: %s", image, err.Error())
				return
			}
			func() {
				mutex.Lock()
				defer mutex.Unlock()
				trivyResponses = append(trivyResponses, response)
			}()
		}(image)
	}
	wg.Wait()

	c.Vulnerabilities.Reset()
	c.VulnerabilitiesSum.Reset()
	for _, trivyResponse := range trivyResponses {
		for _, results := range trivyResponse.Results {
			sevList := []string{}
			for _, vulnerability := range results.Vulnerabilities {
				if vulnerability.Severity != "" {
					sevList = append(sevList, vulnerability.Severity)
				}
				labels := []string{
					trivyResponse.ArtifactName,
					trivyResponse.ArtifactType,
					vulnerability.VulnerabilityID,
					vulnerability.PkgName,
					vulnerability.InstalledVersion,
					vulnerability.Severity,
					vulnerability.FixedVersion,
				}
				c.Vulnerabilities.WithLabelValues(labels...).Set(1)
			}
			sevMap := utils.DupCount(sevList)
			for sev, sevSum := range sevMap {
				labels := []string{
					trivyResponse.ArtifactName,
					trivyResponse.ArtifactType,
					sev,
				}
				c.VulnerabilitiesSum.WithLabelValues(labels...).Set(float64(sevSum))
			}
		}
	}

	if _, err := c.trivyClient.ClearCache(ctx); err != nil {
		return xerrors.Errorf("failed to clear cache: %w", err)
	}

	return nil
}

// StartLoop re-schedules the next collectors run depending on the configured interval
func (c *TrivyCollector) StartLoop(ctx context.Context, interval time.Duration) {
	go func(ctx context.Context) {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				if err := c.Scan(ctx); err != nil {
					c.logger.Errorf("Failed to scan: %s", err.Error())
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}

func (c *TrivyCollector) collectors() []prometheus.Collector {
	return []prometheus.Collector{
		c.Vulnerabilities,
		c.VulnerabilitiesSum,
	}
}

// Describe implements the prometheus colletor interface method
func (c *TrivyCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range c.collectors() {
		collector.Describe(ch)
	}
}

// Collect implements the prometheus colletor interface method
func (c *TrivyCollector) Collect(ch chan<- prometheus.Metric) {
	for _, collector := range c.collectors() {
		collector.Collect(ch)
	}
}

// Name returns the collectors friendly name
func (c *TrivyCollector) Name() string {
	return name
}

// Settings returns the collectors settings
func (c *TrivyCollector) Settings() Settings {
	return c.settings
}
