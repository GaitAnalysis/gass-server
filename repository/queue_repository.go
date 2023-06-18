package repository

import (
	"log"

	"github.com/streadway/amqp"
)

const uploadQueueIdentifier = "UploadQueue"

type QueueRepository struct {
	queueClient  *amqp.Connection
	queueChannel *amqp.Channel
}

func (r *QueueRepository) PublishUpload(message []byte) error {
	return r.queueChannel.Publish(
		"",
		uploadQueueIdentifier,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
}

func NewQueueRepository(queueClient *amqp.Connection) (*QueueRepository, error) {
	queueChannel, err := queueClient.Channel()
	if err != nil {
		return nil, err
	}

	uploadQueue, err := queueChannel.QueueDeclare(
		uploadQueueIdentifier,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	log.Println(uploadQueue)

	return &QueueRepository{
		queueClient:  queueClient,
		queueChannel: queueChannel,
	}, nil
}
