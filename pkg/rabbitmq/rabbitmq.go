package rabbitmq

import (
	"context"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

func OpenChannel() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://admin:admin@rabbitmq:5672/")

	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch, nil
}

func Consume(ch *amqp.Channel, out chan amqp.Delivery, queue string, exchange string, routingKey string) error {

	err := ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		slog.Error("Failed to declare exchange: %s", err)
		slog.Error(err.Error())
		return err
	}

	err = ch.ExchangeDeclare(queue, "direct", true, false, false, false, nil)

	if err != nil {
		slog.Error("Failed to declare a queue: %s", err)
		slog.Error(err.Error())
		return err
	}

	_, err = ch.QueueDeclare(queue, true, false, false, false, nil)

	if err != nil {
		slog.Error("Failed to declare a queue: %s", err)
		slog.Error(err.Error())
		return err
	}

	err = ch.QueueBind(queue, routingKey, exchange, false, nil)

	if err != nil {
		slog.Error("Failed to bind: %s", err)
		slog.Error(err.Error())
		return err
	}

	msgs, err := ch.Consume(
		queue,
		"go-payment",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	for msg := range msgs {
		out <- msg
	}

	return nil
}

func Publish(ctx context.Context, ch *amqp.Channel, body, exchange string) error {
	err := ch.PublishWithContext(ctx, exchange, "PaymentDone", false, false, amqp.Publishing{
		ContentType: "text/json",
		Body:        []byte(body),
	})

	if err != nil {
		return err
	}

	return nil
}
