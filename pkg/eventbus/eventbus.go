package eventbus

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/simplcommerce-go/pkg/config"
)

type Event struct {
	Type      string    `json:"type"`
	Payload   []byte    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler func(ctx context.Context, event Event) error) error
	Close() error
}

type RabbitMQEventBus struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMQEventBus(cfg config.RabbitMQConfig) (EventBus, error) {
	dsn := "amqp://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host + ":" + cfg.Port + "/"
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"simplcommerce_events",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQEventBus{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (r *RabbitMQEventBus) Publish(ctx context.Context, event Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(ctx,
		"",
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQEventBus) Subscribe(eventType string, handler func(ctx context.Context, event Event) error) error {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var event Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("error unmarshalling event: %v", err)
				continue
			}
			if event.Type == eventType {
				if err := handler(context.Background(), event); err != nil {
					log.Printf("error handling event %s: %v", event.Type, err)
				}
			}
		}
	}()

	return nil
}

func (r *RabbitMQEventBus) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}

type InMemoryEventBus struct {
	mu       sync.RWMutex
	handlers map[string][]func(ctx context.Context, event Event) error
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]func(ctx context.Context, event Event) error),
	}
}

func (b *InMemoryEventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers, ok := b.handlers[event.Type]
	b.mu.RUnlock()

	if !ok {
		return nil
	}

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (b *InMemoryEventBus) Subscribe(eventType string, handler func(ctx context.Context, event Event) error) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

func (b *InMemoryEventBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = nil
	return nil
}
