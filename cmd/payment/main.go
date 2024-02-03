package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/almirpask/payment_api/internal/entity"
	"github.com/almirpask/payment_api/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ctx := context.Background()

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgs := make(chan amqp.Delivery)

	go rabbitmq.Consume(ch, msgs, "orders")

	for msg := range msgs {
		var orderRequest entity.OrderRequest
		err := json.Unmarshal(msg.Body, &orderRequest)

		if err != nil {
			slog.Error(err.Error())
			break
		}

		response, err := orderRequest.Process()
		if err != nil {
			slog.Error(err.Error())
			break
		}

		responseJson, err := json.Marshal(response)
		if err != nil {
			slog.Error(err.Error())
			break
		}

		err = rabbitmq.Publish(ctx, ch, string(responseJson), "amq.direct")
		if err != nil {
			slog.Error(err.Error())
			break
		}

		msg.Ack(false)

		slog.Info("Order processed", "order_id", orderRequest.OrderID)
	}

}
