package queue

import (
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/memqueue"
)

type Worker struct {
	Name  string
	Queue taskq.Queue
}

func (w *Worker) Add(m *taskq.Message) error {
	return w.Queue.Add(m)
}

func NewWorker(name string) Worker {
	factory := memqueue.NewFactory()
	mainqueue := factory.RegisterQueue(&taskq.QueueOptions{
		Name: name,
	})

	return Worker{
		Name:  name,
		Queue: mainqueue,
	}
}
