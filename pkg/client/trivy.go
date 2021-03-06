package client

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/xerrors"
)

// TrivyClient struct
type TrivyClient struct{}

// TrivyResponse top level struct to unmarshal the required fields from the json output file
type TrivyResponse struct {
	ArtifactName string `json:"ArtifactName"`
	ArtifactType string `json:"ArtifactType"`
	Results      []struct {
		Target          string               `json:"Target"`
		Vulnerabilities []TrivyVulnerability `json:"Vulnerabilities"`
	} `json:"Results"`
}

// TrivyVulnerability sub struct to unmarshal the required fields from the json output file
type TrivyVulnerability struct {
	VulnerabilityID  string   `json:"VulnerabilityID"`
	PkgName          string   `json:"PkgName"`
	InstalledVersion string   `json:"InstalledVersion"`
	FixedVersion     string   `json:"FixedVersion"`
	Title            string   `json:"Title"`
	Description      string   `json:"Description"`
	Severity         string   `json:"Severity"`
	References       []string `json:"References"`
}

// Do implements the core functionality of the client
func (c *TrivyClient) Do(ctx context.Context, image string) ([]byte, error) {
	tmpfile, err := ioutil.TempFile("", "*.json")
	if err != nil {
		return nil, xerrors.Errorf("failed to create tmpfile: %w", err)
	}
	filename := tmpfile.Name()

	defer tmpfile.Close()
	defer os.Remove(filename)

	result, err := exec.CommandContext(ctx, "trivy", "image", "--skip-update", "--no-progress", "-o", filename, "-f", "json", image).CombinedOutput()
	if err != nil {
		i := strings.Index(string(result), "error in image scan")
		if i == -1 {
			return nil, xerrors.Errorf("failed to execute trivy: %w", err)
		}
		return nil, xerrors.Errorf("failed to execute trivy: %s", result[i:len(result)-1])
	}
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, xerrors.Errorf("failed to read tmpfile: %w", err)
	}
	return body, nil
}

// UpdateDatabase fetches the trivy vulnerability database
func (c *TrivyClient) UpdateDatabase(ctx context.Context) ([]byte, error) {
	return exec.CommandContext(ctx, "trivy", "image", "--download-db-only").CombinedOutput()
}

// ClearCache clears local trivy image cache
func (c *TrivyClient) ClearCache(ctx context.Context) ([]byte, error) {
	return exec.CommandContext(ctx, "trivy", "image", "--clear-cache").CombinedOutput()
}
