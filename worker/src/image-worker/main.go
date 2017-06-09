package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"image/jpeg"

	"github.com/nfnt/resize"
	"github.com/streadway/amqp"
)

type ImageMessage struct {
	Path  string
	Image string
}

func main() {
	imageQueue := os.Getenv("IMAGE_SERVICE_QUEUE")
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		imageQueue, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive,
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,     // queue
		"Worker 1", // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
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

			resizeImage(&m)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func resizeImage(im *ImageMessage) {
	uploadPath := os.Getenv("IMAGE_UPLOAD_PATH")
	imagePath := uploadPath + "/" + im.Image
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

	m := resize.Resize(100, 0, img, resize.Lanczos3)
	newImagePath := uploadPath + "/processed/" + im.Image + "_100.jpg"
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
