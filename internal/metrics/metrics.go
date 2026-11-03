package metrics

import (
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/st0o0/vordr/internal/scanner"
)

type Snapshotter interface {
	Snapshot() []scanner.ImageState
}

type DBTimeFunc func() time.Time

type Collector struct {
	store         Snapshotter
	activeWorkers *atomic.Int64
	queueLen      func() int
	dbTime        DBTimeFunc
	dbUpdated     *prometheus.Desc

	scanStatus    *prometheus.Desc
	vulnCount     *prometheus.Desc
	fixableCount  *prometheus.Desc
	scanTimestamp *prometheus.Desc
	scanDuration  *prometheus.Desc
	containerInfo *prometheus.Desc
	discovered    *prometheus.Desc
	pending       *prometheus.Desc
	scanned       *prometheus.Desc
	errors        *prometheus.Desc
	workersActive *prometheus.Desc
	queueLength   *prometheus.Desc
}

func NewCollector(store Snapshotter, activeWorkers *atomic.Int64, queueLen func() int, dbTime DBTimeFunc) *Collector {
	return &Collector{
		store:         store,
		activeWorkers: activeWorkers,
		queueLen:      queueLen,
		dbTime:        dbTime,
		dbUpdated:     prometheus.NewDesc("vordr_trivy_db_updated_timestamp_seconds", "Unix timestamp of last Trivy vulnerability DB update", nil, nil),
		scanStatus:    prometheus.NewDesc("vordr_image_scan_status", "Image scan state (0=pending, 1=scanning, 2=scanned, 3=error)", []string{"image", "tag"}, nil),
		vulnCount:     prometheus.NewDesc("vordr_vulnerability_count", "Vulnerability count by severity", []string{"image", "tag", "severity"}, nil),
		fixableCount:  prometheus.NewDesc("vordr_fixable_count", "Fixable vulnerability count by severity", []string{"image", "tag", "severity"}, nil),
		scanTimestamp: prometheus.NewDesc("vordr_scan_timestamp_seconds", "Unix timestamp of last scan", []string{"image", "tag"}, nil),
		scanDuration:  prometheus.NewDesc("vordr_scan_duration_seconds", "Duration of last scan in seconds", []string{"image", "tag"}, nil),
		containerInfo: prometheus.NewDesc("vordr_container_image_info", "Container to image mapping", []string{"container", "image", "tag"}, nil),
		discovered:    prometheus.NewDesc("vordr_images_discovered", "Total images discovered", nil, nil),
		pending:       prometheus.NewDesc("vordr_images_pending", "Images pending scan", nil, nil),
		scanned:       prometheus.NewDesc("vordr_images_scanned", "Images successfully scanned", nil, nil),
		errors:        prometheus.NewDesc("vordr_images_errors", "Images in error state", nil, nil),
		workersActive: prometheus.NewDesc("vordr_scan_workers_active", "Currently active scan workers", nil, nil),
		queueLength:   prometheus.NewDesc("vordr_scan_queue_length", "Current scan queue length", nil, nil),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.scanStatus
	ch <- c.vulnCount
	ch <- c.fixableCount
	ch <- c.scanTimestamp
	ch <- c.scanDuration
	ch <- c.containerInfo
	ch <- c.discovered
	ch <- c.pending
	ch <- c.scanned
	ch <- c.errors
	ch <- c.workersActive
	ch <- c.queueLength
	ch <- c.dbUpdated
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	states := c.store.Snapshot()

	var nPending, nScanned, nErrors int

	for _, s := range states {
		img := s.Ref.Image
		tag := s.Ref.Tag

		ch <- prometheus.MustNewConstMetric(c.scanStatus, prometheus.GaugeValue, float64(s.Status), img, tag)

		switch s.Status {
		case scanner.Pending:
			nPending++
		case scanner.Scanned:
			nScanned++
		case scanner.Error:
			nErrors++
		}

		if s.Status == scanner.Scanned {
			for sev, count := range s.Counts {
				ch <- prometheus.MustNewConstMetric(c.vulnCount, prometheus.GaugeValue, float64(count), img, tag, sev)
			}
			for sev, count := range s.Fixable {
				ch <- prometheus.MustNewConstMetric(c.fixableCount, prometheus.GaugeValue, float64(count), img, tag, sev)
			}
			ch <- prometheus.MustNewConstMetric(c.scanTimestamp, prometheus.GaugeValue, float64(s.ScannedAt.Unix()), img, tag)
			ch <- prometheus.MustNewConstMetric(c.scanDuration, prometheus.GaugeValue, s.Duration.Seconds(), img, tag)
		}

		for _, ctr := range s.Containers {
			ch <- prometheus.MustNewConstMetric(c.containerInfo, prometheus.GaugeValue, 1, ctr, img, tag)
		}
	}

	ch <- prometheus.MustNewConstMetric(c.discovered, prometheus.GaugeValue, float64(len(states)))
	ch <- prometheus.MustNewConstMetric(c.pending, prometheus.GaugeValue, float64(nPending))
	ch <- prometheus.MustNewConstMetric(c.scanned, prometheus.GaugeValue, float64(nScanned))
	ch <- prometheus.MustNewConstMetric(c.errors, prometheus.GaugeValue, float64(nErrors))
	ch <- prometheus.MustNewConstMetric(c.workersActive, prometheus.GaugeValue, float64(c.activeWorkers.Load()))
	ch <- prometheus.MustNewConstMetric(c.queueLength, prometheus.GaugeValue, float64(c.queueLen()))

	if t := c.dbTime(); !t.IsZero() {
		ch <- prometheus.MustNewConstMetric(c.dbUpdated, prometheus.GaugeValue, float64(t.Unix()))
	}
}
