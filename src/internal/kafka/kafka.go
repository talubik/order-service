package kafka_consumer

import (
	"context"
	"log"
	"myapp/src/internal/memory"
	"myapp/src/internal/json"

	"github.com/segmentio/kafka-go"
)

func RunConsumer(mem *memory.Memory) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
		GroupID: "my-group",
	})
	defer reader.Close()
	log.Println("Consumer started")
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}
		order , err := json_order.Unmarshal_order(msg.Value)
		if  err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}
		if err := mem.Add(order); err != nil {
			log.Println("Error adding order to memory:", err)
			continue
		}
	}
}
