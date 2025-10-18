package bus

import (
	"context"
)

type Kind string

const (
	KindDuty Kind = "duty"
)

type Event struct {
	Kind   Kind
	Height uint64
	Round  uint64
	Body   any
}

type Subscriber chan Event

type Bus struct {
	pub chan Event
}

func New(size int) *Bus {
	if size <= 0 { size = 128 }
	return &Bus{pub: make(chan Event, size)}
}

func (b *Bus) Publish(_ context.Context, ev Event) {
	select { case b.pub <- ev: default: /* drop on backpressure */ }
}

func (b *Bus) Subscribe() Subscriber { return b.pub }
