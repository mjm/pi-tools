package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	backupKind    = flag.String("kind", "", "Which kind of backup to perform (borg or tarsnap)")
	metricJobName = flag.String("metric-job-name", "backup", "Job name to use when reporting metrics")
)

var (
	lastStartTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backup_last_start_time",
		Help: "The timestamp of the last time a backup started.",
	})
	lastCompletionTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backup_last_completion_time",
		Help: "The timestamp of the last time a backup completed, regardless of result.",
	})
	lastSuccessTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backup_last_success_time",
		Help: "The timestamp of the last time a backup completed successfully.",
	})
)

func main() {
	flag.Parse()

	registry := prometheus.NewRegistry()
	registry.MustRegister(lastStartTime, lastCompletionTime, lastSuccessTime)

	pushGatewayURL := os.Getenv("PUSHGATEWAY_URL")
	if pushGatewayURL == "" {
		log.Panicf("no PUSHGATEWAY_URL provided")
	}

	pusher := push.New(pushGatewayURL, *metricJobName).
		Gatherer(registry).
		Grouping("kind", *backupKind)
	defer func() {
		if err := pusher.Add(); err != nil {
			log.Panicf("Could not push metrics: %v", err)
		}
	}()

	backupTimestamp := time.Now().UTC().Format("2006-01-02_15-04-05")

	switch *backupKind {
	case "borg":
		log.Printf("Performing a borg backup")
		lastStartTime.SetToCurrentTime()

		backupName := fmt.Sprintf("backup-%s", backupTimestamp)
		backupDest := fmt.Sprintf("/dest/backup::%s", backupName)
		cmd := exec.Command(
			"borg",
			"create",
			"--stats",
			backupDest,
			"data")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = os.Getenv("NOMAD_ALLOC_DIR")
		cmd.Env = append([]string{}, os.Environ()...)
		cmd.Env = append(cmd.Env, "BORG_UNKNOWN_UNENCRYPTED_REPO_ACCESS_IS_OK=yes")

		if err := cmd.Run(); err != nil {
			lastCompletionTime.SetToCurrentTime()
			log.Panicf("Running tarsnap failed: %v", err)
		}

		lastCompletionTime.SetToCurrentTime()
		lastSuccessTime.SetToCurrentTime()

		log.Printf("Borg backup completed")

	case "tarsnap":
		log.Printf("Performing a tarsnap backup")

		backupName := fmt.Sprintf("daily-backup-%s", backupTimestamp)
		keyfile := filepath.Join(os.Getenv("NOMAD_SECRETS_DIR"), "tarsnap.key")
		cmd := exec.Command(
			"tarsnap",
			"-c",
			"--keyfile", keyfile,
			"--cachedir", "/var/lib/tarsnap/cache",
			"-f", backupName,
			"--no-default-config",
			"--checkpoint-bytes", "1G",
			"--print-stats",
			"-v",
			"data")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = os.Getenv("NOMAD_ALLOC_DIR")

		if err := cmd.Run(); err != nil {
			lastCompletionTime.SetToCurrentTime()
			log.Panicf("Running tarsnap failed: %v", err)
		}

		lastCompletionTime.SetToCurrentTime()
		lastSuccessTime.SetToCurrentTime()

		log.Printf("Tarsnap backup completed")

	default:
		log.Panicf("Unrecognized backup kind %q", *backupKind)
	}
}
