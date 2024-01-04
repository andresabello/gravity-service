package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type QueueItem struct {
	ID   int
	Data interface{}
}

type ProcessFunc func(QueueItem) error

type QueueService struct {
	queue          chan QueueItem
	failedQueue    chan QueueItem
	maxRetries     int
	retryThreshold int
	WG             sync.WaitGroup
	processFunc    ProcessFunc
}

func NewQueueService(
	maxRetries,
	retryThreshold int,
	processFunc ProcessFunc,
) *QueueService {
	return &QueueService{
		queue:          make(chan QueueItem),
		failedQueue:    make(chan QueueItem),
		maxRetries:     maxRetries,
		retryThreshold: retryThreshold,
		processFunc:    processFunc,
	}
}

func (qs *QueueService) StartProcessor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case item := <-qs.queue:
			qs.WG.Add(1)
			go func(item QueueItem) {
				qs.processQueueItem(ctx, item)
			}(item)
		case failedItem := <-qs.failedQueue:
			qs.WG.Add(1)
			go func(failedItem QueueItem) {
				qs.retryFailedItem(ctx, failedItem)
			}(failedItem)
		}
	}
}

func (qs *QueueService) Enqueue(item QueueItem) {
	qs.queue <- item
}

func (qs *QueueService) processQueueItem(ctx context.Context, item QueueItem) {
	defer qs.WG.Done()

	// Simulate processing
	fmt.Printf("Processing item %d: %s\n", item.ID, item.Data)
	time.Sleep(time.Second)

	if err := qs.processFunc(item); err != nil {
		qs.failedQueue <- item
	} else {
		fmt.Printf("Processed item %d: %s\n", item.ID, item.Data)
	}
}

func (qs *QueueService) processFailedQueueItem(ctx context.Context, item QueueItem, retries int) {
	defer qs.WG.Done()

	// Simulate processing of failed item
	fmt.Printf("Retrying failed item %d (Retry %d): %s\n", item.ID, retries, item.Data)
	time.Sleep(time.Second)

	// Simulate success after retrying
	if err := qs.processFunc(item); err != nil {
		fmt.Printf("permanently failed item %d: %s\n", item.ID, item.Data)
		fmt.Printf("found error %s \n", err.Error())
	}
}

func (qs *QueueService) retryFailedItem(ctx context.Context, item QueueItem) {
	retries := 0
	for retries < qs.maxRetries {
		select {
		case <-ctx.Done():
			return
		default:
			qs.processFailedQueueItem(ctx, item, retries)
			retries++

			// Check if the item succeeded after retry
			if retries < qs.retryThreshold {
				// Simulate success after retrying
				fmt.Printf("successfully retried item %d (Retry %d): %s\n", item.ID, retries, item.Data)
				return
			}

			// If the item still fails after the maximum retries, move it to the permanently failed queue
			fmt.Printf("permanently failed item %d: %s\n", item.ID, item.Data)
			return
		}
	}
}
