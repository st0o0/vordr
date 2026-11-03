package metrics

import (
	"encoding/json"
	"net/http"

	"github.com/st0o0/vordr/internal/scanner"
)

type detailsResponse struct {
	Image           string            `json:"image"`
	Tag             string            `json:"tag"`
	Status          string            `json:"status"`
	ScannedAt       string            `json:"scanned_at,omitempty"`
	Vulnerabilities []detailsVuln     `json:"vulnerabilities"`
}

type detailsVuln struct {
	ID        string `json:"id"`
	Package   string `json:"package"`
	Installed string `json:"installed"`
	Fixed     string `json:"fixed,omitempty"`
	Severity  string `json:"severity"`
}

func detailsHandler(store Snapshotter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		states := store.Snapshot()
		filter := r.URL.Query().Get("image")

		var resp []detailsResponse
		for _, s := range states {
			key := s.Ref.String()
			if filter != "" && key != filter {
				continue
			}

			entry := detailsResponse{
				Image:  s.Ref.Image,
				Tag:    s.Ref.Tag,
				Status: s.Status.String(),
			}

			if !s.ScannedAt.IsZero() {
				entry.ScannedAt = s.ScannedAt.UTC().Format("2006-01-02T15:04:05Z")
			}

			entry.Vulnerabilities = make([]detailsVuln, 0, len(s.Vulnerabilities))
			for _, v := range s.Vulnerabilities {
				entry.Vulnerabilities = append(entry.Vulnerabilities, detailsVuln{
					ID:        v.ID,
					Package:   v.PkgName,
					Installed: v.InstalledVersion,
					Fixed:     v.FixedVersion,
					Severity:  v.Severity,
				})
			}

			resp = append(resp, entry)
		}

		if filter != "" && len(resp) == 0 {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "image not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func snapshotToDetails(states []scanner.ImageState) []detailsResponse {
	resp := make([]detailsResponse, 0, len(states))
	for _, s := range states {
		entry := detailsResponse{
			Image:  s.Ref.Image,
			Tag:    s.Ref.Tag,
			Status: s.Status.String(),
		}
		if !s.ScannedAt.IsZero() {
			entry.ScannedAt = s.ScannedAt.UTC().Format("2006-01-02T15:04:05Z")
		}
		entry.Vulnerabilities = make([]detailsVuln, 0, len(s.Vulnerabilities))
		for _, v := range s.Vulnerabilities {
			entry.Vulnerabilities = append(entry.Vulnerabilities, detailsVuln{
				ID:        v.ID,
				Package:   v.PkgName,
				Installed: v.InstalledVersion,
				Fixed:     v.FixedVersion,
				Severity:  v.Severity,
			})
		}
		resp = append(resp, entry)
	}
	return resp
}
