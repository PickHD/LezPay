package infrastructure

import (
	"context"
	"fmt"

	"github.com/PickHD/LezPay/wallet/internal/v1/application"
	"github.com/segmentio/kafka-go"
)

// ConsumeMessages generic function to consume message from defined param topics
func ConsumeMessages(app *application.App, topicName string) error {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{app.Config.Kafka.FirstBrokerHost, app.Config.Kafka.SecondBrokerHost},
		GroupID:  app.Config.Kafka.GroupID, // Use the same consumer group ID for all readers
		Topic:    topicName,
		MaxBytes: 10e6, // Max message size
	})

	for {
		m, err := consumer.ReadMessage(context.Background())
		if err != nil {
			app.Logger.Error(" consumer.ReadMessage ERROR :", err)
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	if err := consumer.Close(); err != nil {
		app.Logger.Error("failed to close consumer:", err)
		return err
	}

	return nil
}
