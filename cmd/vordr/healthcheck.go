package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/st0o0/vordr/internal/config"
)

func runHealthcheck() int {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		fmt.Printf("config error: %v\n", err)
		return 1
	}
	return doHealthcheck(healthAddr(cfg.ListenAddr))
}

func doHealthcheck(addr string) int {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://" + addr + "/healthz")
	if err != nil {
		fmt.Printf("healthcheck failed: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("unhealthy: status %d\n", resp.StatusCode)
		return 1
	}

	fmt.Println("ok")
	return 0
}

func healthAddr(addr string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	if host == "" || host == "0.0.0.0" || host == "::" {
		return "localhost:" + port
	}
	return addr
}
