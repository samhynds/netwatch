package transporter

import "netwatch/internal/pkg/config"

type Transporter interface {
	Send() error
}

func setupTransporters(config *config.Config) {

	var transporters = make([]Transporter, 2)

	if config.Config.Transporters.Database {

	}

	if config.Config.Transporters.Queue {

	}

}
