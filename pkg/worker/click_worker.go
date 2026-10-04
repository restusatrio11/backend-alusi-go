package worker

import (
	"context"
	"time"

	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"

	"github.com/rs/zerolog/log"
)

type ClickWorker struct {
	clickLogRepo *postgres.ClickLogRepo
	queue        chan *domain.ClickLog
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewClickWorker(clickLogRepo *postgres.ClickLogRepo, bufferSize int, workerCount int) *ClickWorker {
	ctx, cancel := context.WithCancel(context.Background())
	cw := &ClickWorker{
		clickLogRepo: clickLogRepo,
		queue:        make(chan *domain.ClickLog, bufferSize),
		ctx:          ctx,
		cancel:       cancel,
	}

	for i := 1; i <= workerCount; i++ {
		go cw.startWorker(i)
	}

	return cw
}

func (w *ClickWorker) startWorker(workerID int) {
	for {
		select {
		case <-w.ctx.Done():
			return
		case item, ok := <-w.queue:
			if !ok {
				return
			}
			if w.clickLogRepo != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				if err := w.clickLogRepo.Record(ctx, item); err != nil {
					log.Warn().Err(err).Int("app_id", item.AppID).Msg("Failed to persist click log")
				}
				cancel()
			}
		}
	}
}

// Enqueue queues a click log item asynchronously without blocking
func (w *ClickWorker) Enqueue(logItem *domain.ClickLog) {
	select {
	case w.queue <- logItem:
	default:
		log.Warn().Int("app_id", logItem.AppID).Msg("Click log queue is full, dropping item to preserve API latency")
	}
}

// Stop gracefully shuts down the worker
func (w *ClickWorker) Stop() {
	w.cancel()
}
