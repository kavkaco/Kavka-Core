package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NATSQueue struct {
	js nats.JetStreamContext
}

func NewNATSQueue(js nats.JetStreamContext) (*NATSQueue, error) {
	if js == nil {
		return nil, fmt.Errorf("jetstream context is nil")
	}

	_, err := js.AddStream(&nats.StreamConfig{
		Name:      "kavka-tasks",
		Subjects:  []string{"tasks.>"},
		Storage:   nats.FileStorage,
		MaxAge:    24 * 7 * 3, // 3 weeks
		Replicas:  1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tasks stream: %w", err)
	}

	return &NATSQueue{js: js}, nil
}

func (q *NATSQueue) Enqueue(ctx context.Context, task Task) error {
	subject := fmt.Sprintf("tasks.%s", task.Type)

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	_, err = q.js.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish task: %w", err)
	}

	return nil
}

func (q *NATSQueue) EnqueueBatch(ctx context.Context, tasks []Task) error {
	for _, task := range tasks {
		err := q.Enqueue(ctx, task)
		if err != nil {
			return err
		}
	}

	return nil
}

func (q *NATSQueue) Close() error {
	return nil
}
