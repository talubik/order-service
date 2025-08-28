package main

import (
	"context"
	"log"
	"time"

	"myapp/src/internal/json"
	"myapp/src/testutils"

	"github.com/segmentio/kafka-go"
)

func main() {

	orders := generator.GenerateOrders(500)

	// Создаем продюсера Kafka
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"), 
		Topic:    "orders",                    
		BatchTimeout: 10 *time.Millisecond,      
	}
	defer writer.Close()
	for i, order:= range orders {
		data , _ :=json_order.Marshall_order(&order)
		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Value: data,           // Тело сообщения (JSON)
			},
		)
		if err != nil {
			log.Printf("Ошибка отправки сообщения %d: %v", i, err)
		}else{
			log.Printf("Сообщение отправлено %d",i)
		}
	}
}