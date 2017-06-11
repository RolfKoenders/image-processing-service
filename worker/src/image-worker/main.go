package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nfnt/resize"
	"github.com/streadway/amqp"
)

type ImageMessage struct {
	Path  string
	Image string
}

func connectToRabbit(uri string) *amqp.Connection {
	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672")
		if err == nil {
			log.Println("Connected to RabbitMQ")
			// Create channel with rabbit
			// Stop trying to connect
			return conn
		}

		fmt.Println("Trying to reconnect to RabbitMQ")
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	var sizesToScale = []int{50, 125, 250, 500}

	IMAGEQUEUE := os.Getenv("IMAGE_SERVICE_QUEUE")

	conn := connectToRabbit("amqp://guest:guest@rabbitmq:5672")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		IMAGEQUEUE, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive,
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		os.Getenv("WORKER_NAME"), // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer/worker")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			body := d.Body

			var m ImageMessage
			err := json.Unmarshal(body, &m)
			if err != nil {
				log.Fatal("Error parsing json in message")
			}

			log.Printf("Received a message: %s", m.Image)

			var wg sync.WaitGroup
			for _, size := range sizesToScale {
				log.Println("Create routine ")
				wg.Add(1)
				go resizeImage(&m, size)
			}

			d.Ack(false)
		}
	}()

	log.Printf("[*] %v Waiting for messages. To exit press CTRL+C", os.Getenv("WORKER_NAME"))
	<-forever
}

func resizeImage(im *ImageMessage, width int) {
	imagePath := im.Path + "/" + im.Image
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resize.Resize(uint(width), 0, img, resize.Lanczos3)

	newImagePath := fmt.Sprintf("%s/processed/%s_%v.jpg", im.Path, im.Image, width)
	out, err := os.Create(newImagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	jpeg.Encode(out, m, nil)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
