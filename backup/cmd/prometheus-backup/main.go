package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	prometheusURL      = flag.String("prometheus-url", "", "Base URL for accessing prometheus")
	prometheusDataPath = flag.String("prometheus-data-path", "", "Path to directory where Prometheus is storing its data")
)

const destPath = "/backup/prometheus"

func main() {
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	backupName, err := takeSnapshot(ctx)
	if err != nil {
		log.Panicf("Failed to take snapshot: %v", err)
	}

	if err := os.MkdirAll(destPath, 0755); err != nil {
		log.Panicf("Failed to create directory at %s: %v", destPath, err)
	}

	snapshotPath := filepath.Join(*prometheusDataPath, "snapshots", backupName)
	if err := exec.CommandContext(ctx, "rsync", "-av", snapshotPath+"/", destPath).Run(); err != nil {
		log.Panicf("Failed to copy snapshot to backup destination: %v", err)
	}

	log.Printf("Successfully copied snapshot from %s to %s", snapshotPath, destPath)

	if err := os.RemoveAll(snapshotPath); err != nil {
		log.Printf("Warning: Failed to remove snapshot at %s: %v", snapshotPath, err)
		return
	}

	log.Printf("Successfully removed snapshot from Prometheus storage")
}

func takeSnapshot(ctx context.Context) (string, error) {
	u := fmt.Sprintf("%s/api/v1/admin/tsdb/snapshot", *prometheusURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status code: %d", res.StatusCode)
	}

	var resp struct {
		Data struct {
			Name string
		}
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", err
	}

	if resp.Data.Name == "" {
		return "", fmt.Errorf("no snapshot name in response from prometheus")
	}

	return resp.Data.Name, nil
}
