package worker

import (
	"context" //сигналы отмены
	"log"     //вывод сообщений в консоль
	"time"
)

type Worker struct {
	interval time.Duration
}

type Request struct {
	direction string
	date      string
	from      string
	numberOfPeople int
}

requests := []Request{}

func (r *Request) NewRequest(direction, date, from string, numberOfPeople int) {
	requests = append(requests, Request{
		direction:      direction,
		date:           date,
		from:           from,
		numberOfPeople: numberOfPeople,
	})
}


// NewWorker — конструктор
func NewWorker(interval time.Duration) *Worker {
	return &Worker{
		interval: interval,
	}
}

// Run запускает бесконечный цикл проверки по таймеру
func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop() // срабатывает при выходе из функции

	log.Println("Воркер запущен. Интервал:", w.interval)

	for {
		select {
		case <-ctx.Done():
			// Сработает, когда ты нажмёшь Ctrl+C
			log.Println("Получен сигнал остановки, выключаю воркер...")
			return
		case <-ticker.C:
			// Срабатывает каждый раз, когда проходит interval
			log.Println("Тик таймера: проверяю билеты...")
		}
	}
}
