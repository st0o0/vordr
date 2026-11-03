package watcher

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/st0o0/vordr/internal/scanner"
	"github.com/st0o0/vordr/internal/store"
)

type Watcher struct {
	docker  client.APIClient
	store   *store.Store
	scanCh  chan<- scanner.ImageRef
	stagger time.Duration
	logger  *slog.Logger
}

func New(docker client.APIClient, s *store.Store, scanCh chan<- scanner.ImageRef, stagger time.Duration, logger *slog.Logger) *Watcher {
	return &Watcher{
		docker:  docker,
		store:   s,
		scanCh:  scanCh,
		stagger: stagger,
		logger:  logger,
	}
}

func (w *Watcher) Run(ctx context.Context) error {
	containers, err := w.docker.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return err
	}

	var initial []scanner.ImageRef
	for _, c := range containers {
		ref := ParseImageRef(c.Image)
		ref.ImageID = c.ImageID
		name := containerName(c.Names)
		if w.store.Add(ref, name) {
			initial = append(initial, ref)
		}
	}

	w.logger.Info("initial container listing", "containers", len(containers), "unique_images", len(initial))

	for i, ref := range initial {
		if i > 0 {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(w.stagger):
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case w.scanCh <- ref:
		}
	}

	return w.watch(ctx)
}

func (w *Watcher) watch(ctx context.Context) error {
	f := filters.NewArgs()
	f.Add("type", string(events.ContainerEventType))
	f.Add("event", "start")
	f.Add("event", "die")

	msgCh, errCh := w.docker.Events(ctx, events.ListOptions{Filters: f})

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			if ctx.Err() != nil {
				return nil
			}
			return err
		case msg := <-msgCh:
			w.handleEvent(ctx, msg)
		}
	}
}

func (w *Watcher) handleEvent(ctx context.Context, msg events.Message) {
	name := msg.Actor.Attributes["name"]
	image := msg.Actor.Attributes["image"]
	if name == "" || image == "" {
		return
	}

	ref := ParseImageRef(image)

	switch msg.Action {
	case events.ActionStart:
		if imageID := w.inspectImageID(ctx, msg.Actor.ID); imageID != "" {
			ref.ImageID = imageID
		}

		if w.store.Add(ref, name) {
			w.logger.Info("new image discovered", "image", ref.String(), "container", name)
			select {
			case <-ctx.Done():
			case w.scanCh <- ref:
			}
		} else {
			w.logger.Debug("container started with known image", "image", ref.String(), "container", name)
		}

	case events.ActionDie:
		w.store.Remove(name)
		w.logger.Debug("container stopped", "container", name)
	}
}

func (w *Watcher) inspectImageID(ctx context.Context, containerID string) string {
	info, err := w.docker.ContainerInspect(ctx, containerID)
	if err != nil {
		w.logger.Debug("failed to inspect container for ImageID", "container", containerID, "error", err)
		return ""
	}
	return info.Image
}

func ParseImageRef(image string) scanner.ImageRef {
	if idx := strings.Index(image, "@sha256:"); idx >= 0 {
		image = image[:idx]
	}

	lastSlash := strings.LastIndex(image, "/")
	tagPart := image
	if lastSlash >= 0 {
		tagPart = image[lastSlash:]
	}

	if idx := strings.LastIndex(tagPart, ":"); idx >= 0 {
		absIdx := len(image) - len(tagPart) + idx
		return scanner.ImageRef{
			Image: image[:absIdx],
			Tag:   image[absIdx+1:],
		}
	}

	return scanner.ImageRef{
		Image: image,
		Tag:   "latest",
	}
}

func containerName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return strings.TrimPrefix(names[0], "/")
}
