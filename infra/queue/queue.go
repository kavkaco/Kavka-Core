package queue

import (
	"context"
	"encoding/json"
	"time"
)

type TaskType string

const (
	TaskSendEmail        TaskType = "email.send"
	TaskSendVerification TaskType = "email.verify"
	TaskResetPassword    TaskType = "email.reset_password"
	TaskProcessImage     TaskType = "image.process"
	TaskNotifyUser       TaskType = "notify.user"
)

type Task struct {
	ID        string            `json:"id"`
	Type      TaskType          `json:"type"`
	Payload   json.RawMessage   `json:"payload"`
	Metadata  map[string]string `json:"metadata"`
	MaxRetries int              `json:"max_retries"`
	RetryCount int              `json:"retry_count"`
	CreatedAt time.Time         `json:"created_at"`
}

type Queue interface {
	Enqueue(ctx context.Context, task Task) error
	EnqueueBatch(ctx context.Context, tasks []Task) error
	Close() error
}
