package transporter

import "netwatch/internal/pkg/config"

type Manager struct {
	Queue *TransportQueue
}

func NewManager(config *config.Config) *Manager {
	queue := NewTransportQueue(config.Config.Queue.Capacity)

	return &Manager{
		Queue: queue,
	}
}
