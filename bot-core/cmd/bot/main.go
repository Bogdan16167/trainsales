package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Trainsales/internal/worker"
)

func main() {
	// Контекст автоматически отменится, если нажать Ctrl+C в терминале
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Для теста ставим 5 секунд, чтобы долго не ждать
	w := worker.NewWorker(5 * time.Second)
	w.Run(ctx)
}
