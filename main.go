package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"pubsub/pub"
	"pubsub/sub"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
)

// 1640423777_7
func pubsubHandler(ctx context.Context, m *pubsub.Message) {
	fmt.Println("[SUB] received message: ", m.ID, " now: ", time.Now().Unix())
	fmt.Println(string(m.Data), "\n", "\n")

	// mark item process success
	//m.Ack()

	// mark item process fail, resend
	//m.Nack()
}

func main() {
	projectId := os.Getenv("PROJECT_ID")
	topicId := os.Getenv("TOPIC_ID")
	subId := os.Getenv("SUB_ID")
	fmt.Println(projectId, topicId, subId)

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("create client success")

	topic := client.Topic(topicId)
	pubClient := pub.NewPublisher(topic)
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			messageId := fmt.Sprintf("%d_%d", time.Now().Unix(), i)
			if err := pubClient.PublishMessage(map[string]string{
				"key":        "val",
				"message_id": messageId,
			}); err != nil {
				fmt.Println("[PUB] Publish message error ", messageId)
				continue
			}

			fmt.Println("[PUB] Publish message success ", messageId)
		}
	}()

	subscriber := sub.NewSubscriber(client, pubsubHandler, subId)
	go func() {
		log.Println("subcription start")
		err = subscriber.Start()
		if err != nil {
			log.Println("start consumer error", err)
			return
		}

		log.Println("start consumer success")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c)
	for {
		select {
		case sig := <-c:
			if sig == syscall.SIGKILL || sig == syscall.SIGINT {
				log.Println("Received signal kill process")
				client.Close()
				return
			}
		}
	}
}
