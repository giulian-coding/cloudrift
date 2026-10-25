package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/giulian-coding/cloudrift/internal/order"
	"github.com/giulian-coding/cloudrift/internal/user"
	"golang.org/x/term"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/oci"
)

func main_2() {
	userRepo := user.NewMemoryRepository()
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	orderRepo := order.NewMemoryRepository()
	orderSvc := order.NewService(orderRepo)
	ordersHandler := order.NewHandler(orderSvc)

	mux := http.NewServeMux()
	userHandler.Register(mux)
	ordersHandler.Register(mux)
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server is running...")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("server failed to start", "error", err)
	}
}

func runInteractiveShell(ctx context.Context, container containerd.Container) (uint32, error) {
	// Task mit Standard-IO und Terminal-Emulation erstellen
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio, cio.WithTerminal))
	if err != nil {
		return 0, fmt.Errorf("fehler beim Erstellen des Tasks: %w", err)
	}
	defer task.Delete(ctx)

	// Host-Terminal in den Raw-Modus versetzen (für interaktive Eingaben)
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return 0, fmt.Errorf("fehler beim Setzen des Raw-Modus: %w", err)
	}
	// Sicherstellen, dass das Terminal beim Verlassen der Funktion wieder normal wird
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Auf Task-Ende warten (Muss vor Start() aufgerufen werden)
	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		return 0, fmt.Errorf("fehler beim Warten auf den Task: %w", err)
	}

	// Task (die Shell) starten
	if err := task.Start(ctx); err != nil {
		return 0, fmt.Errorf("fehler beim Starten des Tasks: %w", err)
	}

	// Blockieren, bis der Task beendet ist
	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		return 0, fmt.Errorf("task beendet mit fehlerhaftem Status: %w", err)
	}
	return code, nil
}

func main() {
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

/*

Connect to Container
ctr -n default tasks exec --exec-id shell -t golang sh

Kill container
ctr -n default tasks kill --signal SIGKILL golang

Delete Task
ctr -n default tasks delete golang

Delete Container
ctr -n default containers delete golang

*/
