package event_consumer

import (
	"log"
	"time"

	"github.com/lastZu/Esteem/internal/app/events"
)

type Consumer struct {
	fatcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fatcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fatcher:   fatcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fatcher.Fetch(c.batchSize)

		if err != nil {
			log.Printf("[ERR] consumer %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)
			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())
			continue
		}
	}

	return nil
}
