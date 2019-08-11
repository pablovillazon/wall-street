package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// Quote data type
type Quote struct {
	QuoteType string
	Symbol    string
	DateTime  string
	Open      float32
	High      float32
	Low       float32
	Close     float32
	Volume    int
}

func main() {
	conn, err := amqp.Dial("amqp://admin:Password123@159.65.220.217:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	qInput, err := ch.QueueDeclare(
		"leonardo-dispatcher-qeue", // name
		false,                      // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		qInput.Name, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	quote := Quote{}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			err := json.Unmarshal(d.Body, &quote)
			failOnError(err, "Failed unmarshaling json data")
			data := getQuoteAsString(quote)

			log.Printf("Received: %s", data)

			log.Printf("Sending to A: %s", data)
			dispatch(conn, "leonardo-client-a-queue", d.Body)

			log.Printf("Sending to B: %s", data)
			dispatch(conn, "leonardo-client-b-queue", d.Body)

			log.Printf("Sending to C: %s", data)
			dispatch(conn, "leonardo-client-c-queue", d.Body)

			log.Printf("Sending to D: %s", data)
			dispatch(conn, "leonardo-client-d-queue", d.Body)
		}
	}()

	log.Printf(" > Waiting for messages. To exit press CTRL+C")
	<-forever

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getQuoteAsString(quote Quote) string {
	return "[" + quote.QuoteType + " " + quote.Symbol + " " + quote.DateTime + " " + fmt.Sprintf("%.2f", quote.Open) + " " + fmt.Sprintf("%.2f", quote.High) + " " + fmt.Sprintf("%.2f", quote.Low) + " " + fmt.Sprintf("%.2f", quote.Close) + " " + fmt.Sprintf("%d", quote.Volume) + "]"
}

func dispatch(conn *amqp.Connection, queueName string, body []byte) {
	// conn, err := amqp.Dial("amqp://admin:Password123@159.65.220.217:5672")
	// failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
}
