package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/oci"
)

type Repository interface {
	Create(context.Context, Container) error
	Delete(context.Context, string) error
	ByID(context.Context, string) (Container, error)
}

type Service struct {
	repo   Repository
	client *containerd.Client
}

func NewService(repo Repository, client *containerd.Client) *Service {
	return &Service{
		repo:   repo,
		client: client,
	}
}

func (s *Service) createContainer(ctx context.Context, imageURL string) (string, error) {
	id := fmt.Sprintf("dev-shell-%d", time.Now().UnixNano())
	snapshotName := id + "-snapshot"

	image, err := s.client.Pull(ctx, imageURL, containerd.WithPullUnpack)
	if err != nil {
		return "", fmt.Errorf("pull image: %w", err)
	}

	ctr, err := s.client.NewContainer(ctx, id,
		containerd.WithNewSnapshot(snapshotName, image),
		containerd.WithNewSpec(
			oci.WithImageConfig(image),
			oci.WithProcessArgs("/bin/sh"),
			oci.WithTTY,
		),
	)
	if err != nil {
		return "", fmt.Errorf("create container: %w", err)
	}

	if err := s.repo.Create(ctx, Container{ID: id}); err != nil {
		if cleanupErr := ctr.Delete(ctx, containerd.WithSnapshotCleanup); cleanupErr != nil {
			return "", fmt.Errorf(
				"save container in repository: %w; cleanup also failed: %v",
				err,
				cleanupErr,
			)
		}
		return "", fmt.Errorf("save container in repository: %w", err)
	}
	return id, nil
}

func (s *Service) deleteContainer(ctx context.Context, id string) error {
	if _, err := s.repo.ByID(ctx, id); err != nil {
		return fmt.Errorf("find container in repository: %w", err)
	}
	ctr, err := s.client.LoadContainer(ctx, id)
	if err != nil {
		return fmt.Errorf("load container: %w", err)
	}
	if err := ctr.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
		return fmt.Errorf("delete container: %w", err)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("remove container from repository: %w", err)
	}

	return nil
}
