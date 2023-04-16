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

type Config struct {
	ProjectId string
	TopicId   string
	SubId     string
}

var cfg *Config

func init() {
	cfg = &Config{
		ProjectId: os.Getenv("PROJECT_ID"),
		TopicId:   os.Getenv("TOPIC_ID"),
		SubId:     os.Getenv("SUB_ID"),
	}
}

func pubsubHandlerNack(ctx context.Context, m *pubsub.Message) {
	fmt.Println("[SUB] received message: ", m.ID, " now: ", time.Now().Unix())
	fmt.Println(string(m.Data), "\n", "\n", m.DeliveryAttempt)
	m.Nack()

	// mark item process success
	//m.Ack()

	// mark item process fail, resend
	//m.Nack()
}

// 1640423777_7
func pubsubHandler(ctx context.Context, m *pubsub.Message) {
	fmt.Println("[SUB] received message: ", m.ID, " now: ", time.Now().Unix())
	fmt.Println(string(m.Data), "\n", "\n")

	// mark item process success
	//m.Ack()

	// mark item process fail, resend
	//m.Nack()
}

func blockAndWait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	for {
		select {
		case sig := <-c:
			if sig == syscall.SIGKILL || sig == syscall.SIGINT {
				log.Println("Received signal kill process")
				return
			}
		}
	}
}

func demoDeadLetter(cli *pubsub.Client) {
	defer cli.Close()
	topic := cli.Topic(cfg.TopicId)
	pubClient := pub.NewPublisher(topic)
	err := pubClient.PublishMessage(`{"user_name":"minh.nguyen3"}`)
	if err != nil {
		panic(err)
	}
	subscriber := sub.NewSubscriber(cli, pubsubHandlerNack, cfg.SubId)
	go func() {
		log.Println("subcription start")
		err = subscriber.Start()
		if err != nil {
			log.Println("start consumer error", err)
			return
		}

		log.Println("start consumer success")
	}()

	blockAndWait()
}

func main() {
	fmt.Println(cfg.ProjectId, cfg.TopicId, cfg.SubId)
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, cfg.ProjectId)
	if err != nil {
		log.Fatal(err)
	}

	args := os.Args
	if len(args) == 2 {

		switch args[1] {
		case "dead_letter":
			demoDeadLetter(client)
			return
		}
	}

	log.Println("create client success")
	topic := client.Topic(cfg.TopicId)
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

	subscriber := sub.NewSubscriber(client, pubsubHandler, cfg.SubId)
	go func() {
		log.Println("subcription start")
		err = subscriber.Start()
		if err != nil {
			log.Println("start consumer error", err)
			return
		}

		log.Println("start consumer success")
	}()
	blockAndWait()
	client.Close()
}
