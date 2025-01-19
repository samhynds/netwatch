package kafka

import "log"

func Setup() {
	// connect to db
	// setup db & tables if they don't exist
	// return reference for further use
	log.Println("Setting up Kafka")
}

func Send() {
	// saves a TransportQueueItem to the Kafka
	log.Println("Sending to Kafka")
}
