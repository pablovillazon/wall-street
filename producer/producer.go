package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

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

	q, err := ch.QueueDeclare(
		"leonardo-dispatcher-qeue", // name
		false,                      // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	dt := time.Now()
	myRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		dateTime := dt.Format("15:04:05.000000")
		open := (myRand.Float32() * 100) + 100
		high := (myRand.Float32() * 100) + 100
		low := (myRand.Float32() * 100) + 100
		close := (myRand.Float32() * 100) + 100
		volume := myRand.Intn(1000)

		quote := Quote{
			QuoteType: "bid",
			Symbol:    "GOOG",
			DateTime:  dateTime,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}

		jsonData, err := json.Marshal(quote)
		failOnError(err, "Failed converting data to json")

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(jsonData),
			})
		data := getQuoteAsString(quote)
		log.Printf("Sent: %s", data)
		failOnError(err, "Failed to publish a message")

		time.Sleep(time.Second)
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getQuoteAsString(quote Quote) string {
	return "[" + quote.QuoteType + " " + quote.Symbol + " " + quote.DateTime + " " + fmt.Sprintf("%.2f", quote.Open) + " " + fmt.Sprintf("%.2f", quote.High) + " " + fmt.Sprintf("%.2f", quote.Low) + " " + fmt.Sprintf("%.2f", quote.Close) + " " + fmt.Sprintf("%d", quote.Volume) + "]"
}
