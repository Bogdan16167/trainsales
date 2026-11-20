package worker

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Worker struct {
	interval time.Duration
	requests []Request
	addCh    chan Request
}

type Request struct {
	Direction      string
	Date           string
	From           string
	NumberOfPeople int
	ID             int
}

func NewWorker(interval time.Duration) *Worker {
	return &Worker{
		interval: interval,
		requests: make([]Request, 0),
		addCh:    make(chan Request),
	}
}

func (w *Worker) AddRequest(req Request) {
	w.addCh <- req
}

func (w *Worker) RemoveRequest(id int) {
	for i, req := range w.requests {
		if req.ID == id {
			w.requests = append(w.requests[:i], w.requests[i+1:]...)
			return
		}
	}
}

func (w *Worker) Run(ctx context.Context) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	ticker := time.NewTicker(w.interval)
	removeList := []Request{}
	defer ticker.Stop()

	log.Println("Воркер запущен. Интервал:", w.interval)

	for {
		select {
		case <-ctx.Done():
			//Ctrl+C
			log.Println("Получен сигнал остановки, выключаю воркер...")
			return
		case <-ticker.C:
			// Срабатывает каждый раз, когда проходит interval
			if len(w.requests) == 0 {
				log.Println("Нет запросов для проверки.")
				continue
			}
			for _, req := range w.requests {
				wg.Add(1)
				go func(req Request) {
					defer wg.Done()
					log.Printf("Проверяю: %s на %s (%d чел.)", req.Direction, req.Date, req.NumberOfPeople)
					if checkTicket(req) {
						log.Printf("Билеты найдены! %s на %s (%d чел.)", req.Direction, req.Date, req.NumberOfPeople)
						mu.Lock()
						removeList = append(removeList, req)
						mu.Unlock()
					}
				}(req)
			}
			wg.Wait()
			for _, req := range removeList {
				log.Printf("Удаляю запрос: %s на %s (%d чел.)", req.Direction, req.Date, req.NumberOfPeople)
				w.RemoveRequest(req.ID)
			}
			removeList = []Request{}
		case newreq := <-w.addCh:
			log.Printf("Добавлен новый запрос: %s на %s (%d чел.)", newreq.Direction, newreq.Date, newreq.NumberOfPeople)
			w.requests = append(w.requests, newreq)
		}
	}
}

func checkTicket(req Request) bool {
	time.Sleep(1 * time.Second)
	return rand.Intn(10) < 3
}
