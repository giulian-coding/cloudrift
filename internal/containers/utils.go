package containers

import (
	"context"
	"fmt"
	"os"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"golang.org/x/term"
)

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
