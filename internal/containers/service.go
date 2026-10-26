package containers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/oci"
)

type Repository interface {
	Create(context.Context, Container) error
	ByID(context.Context, string) (Container, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) createContainer(ctx context.Context, imageURL string) {

}

func containers() {
	ctx := context.Background()

	uniqueID := fmt.Sprintf("dev-shell-%d", time.Now().Unix())
	snapshotName := fmt.Sprintf("%s-snapshot", uniqueID)

	client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace("default"))
	if err != nil {
		log.Fatalf("connection to containerd failed: %v", err)
	}
	defer client.Close()

	image, err := client.Pull(ctx, "docker.io/library/golang:1.27.3-bookworm", containerd.WithPullUnpack)
	if err != nil {
		log.Fatalf("image pull failed: %v", err)
	}

	ctr, err := client.NewContainer(ctx, "golang",
		containerd.WithNewSnapshot(snapshotName, image),
		containerd.WithNewSpec(
			oci.WithImageConfig(image),
			oci.WithProcessArgs("/bin/sh"),
			oci.WithTTY,
		),
	)
	if err != nil {
		log.Fatalf("container creation failed: %v", err)
	}
	defer ctr.Delete(ctx, containerd.WithSnapshotCleanup)

	exitCode, err := runInteractiveShell(ctx, ctr)
	if err != nil {
		log.Fatalf("\nFehler in der Shell: %v\n", err)
	}

	fmt.Printf("container stopped with exit code: %d\n", exitCode)
}
