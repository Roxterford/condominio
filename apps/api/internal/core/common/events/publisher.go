package events

import "context"

type EventPublisher interface {
	// Publish sends an event to the transport (Redis, RabbitMQ, etc.)
	Publish(ctx context.Context, event Event) error
}
