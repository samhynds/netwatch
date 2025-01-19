package transporter

import (
	"netwatch/internal/pkg/config"
	transporterqueue "netwatch/internal/pkg/transporter/queue"
)

type Manager struct {
	Queue *transporterqueue.Queue
}

func NewManager(config *config.Config) *Manager {
	queue := transporterqueue.NewQueue(config.Config.Queue.Capacity)

	return &Manager{
		Queue: queue,
	}
}
