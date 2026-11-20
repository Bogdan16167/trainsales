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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	w := worker.NewWorker(5 * time.Second)
	r := []worker.Request{
		{
			Direction:      "Moscow - Saint Petersburg",
			Date:           "2024-07-01",
			From:           "Moscow",
			NumberOfPeople: 2,
			ID:             1,
		},
		{
			Direction:      "Saint Petersburg - Moscow",
			Date:           "2024-07-02",
			From:           "Saint Petersburg",
			NumberOfPeople: 1,
			ID:             2,
		},
		{
			Direction:      "Moscow - Kazan",
			Date:           "2024-07-03",
			From:           "Moscow",
			NumberOfPeople: 3,
			ID:             3,
		},
		{
			Direction:      "Kazan - Moscow",
			Date:           "2024-07-04",
			From:           "Kazan",
			NumberOfPeople: 2,
			ID:             4,
		},
	}

	go w.Run(ctx)

	for _, req := range r {
		w.AddRequest(req)
	}
	<-ctx.Done()
}
