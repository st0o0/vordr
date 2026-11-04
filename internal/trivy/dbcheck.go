package trivy

import (
	"context"
	"encoding/json"
	"log/slog"
	"os/exec"
	"sync/atomic"
	"time"
)

type DBChecker struct {
	cacheDir string
	interval time.Duration
	enabled  bool
	dbTime   atomic.Pointer[time.Time]
	logger   *slog.Logger
}

func NewDBChecker(cacheDir string, interval time.Duration, logger *slog.Logger, enabled bool) *DBChecker {
	return &DBChecker{
		cacheDir: cacheDir,
		interval: interval,
		enabled:  enabled,
		logger:   logger,
	}
}

func (c *DBChecker) DBUpdatedAt() time.Time {
	if t := c.dbTime.Load(); t != nil {
		return *t
	}
	return time.Time{}
}

func (c *DBChecker) Run(ctx context.Context) error {
	if !c.enabled {
		c.logger.Info("db checker disabled (server mode)", "event", "db_check_disabled")
		<-ctx.Done()
		return nil
	}

	c.check(ctx)

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			c.check(ctx)
		}
	}
}

func (c *DBChecker) check(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "trivy", "version", "--format", "json", "--cache-dir", c.cacheDir)
	out, err := cmd.Output()
	if err != nil {
		c.logger.Error("db check failed", "event", "db_check_failed", "error", err)
		return
	}

	var ver versionOutput
	if err := json.Unmarshal(out, &ver); err != nil {
		c.logger.Error("parse trivy version failed", "event", "db_check_failed", "error", err)
		return
	}

	t := ver.VulnerabilityDB.UpdatedAt
	c.dbTime.Store(&t)
	c.logger.Info("db check completed", "event", "db_check", "updated_at", t)
}

type versionOutput struct {
	VulnerabilityDB struct {
		UpdatedAt time.Time `json:"UpdatedAt"`
	} `json:"VulnerabilityDB"`
}
