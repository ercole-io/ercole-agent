package queue

import (
	"context"
	"time"

	"github.com/ercole-io/ercole-agent/v2/client"
	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/vmihailenco/taskq/v3"
)

type Consumer struct {
	Name string
	Task *taskq.Task
}

func (c *Consumer) NewMessage(ctx context.Context, args ...interface{}) *taskq.Message {
	return c.Task.WithArgs(ctx, args...)
}

func NewConsumer(name string, waitingTime time.Duration, retry int,
	f func(logger.Logger, *client.Client, config.Configuration, interface{}, string) error) Consumer {
	opts := &taskq.TaskOptions{
		Name:       name,
		MinBackoff: waitingTime,
		MaxBackoff: waitingTime,
		RetryLimit: retry,
		Handler:    f,
	}

	return Consumer{
		Name: name,
		Task: taskq.RegisterTask(opts),
	}
}
