package stats

import (
	"context"
)

type Collector struct {
	EventCh chan Event
	Store   *Store
}

func NewCollector(eventCh chan Event, store *Store) *Collector {
	return &Collector{
		EventCh: eventCh,
		Store:   store,
	}
}

func (c *Collector) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-c.EventCh:
			c.Store.Apply(event)
		}
	}
}
