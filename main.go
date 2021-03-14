package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
)

func main() {
	log.Println("producer is starting...")
	defer log.Println("producer is going down...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handleSignals(cancel)

	gcpProjectName := os.Getenv("GCP_PROJECT")
	pubsubTopicName := os.Getenv("PAYMENT_TOPIC")

	client, err := pubsub.NewClient(ctx, gcpProjectName)
	if err != nil {
		log.Fatalf("failed creating pubsub client: %v", err)
	}
	defer client.Close()

	paymentTopic := client.Topic(pubsubTopicName)
	for {
		log.Printf("publishing event...")

		_, err := paymentTopic.Publish(ctx, &pubsub.Message{
			Data: []byte("payment event!"),
		}).Get(ctx)
		if err != nil {
			log.Fatalf("failed to publish payment event: %v", err)
		}

		select {
		case <-time.After(10 * time.Second):

		case <-ctx.Done():
			return
		}
	}
}

func handleSignals(doneFunc func()) {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalC
		log.Printf("got signal: %v", sig)
		doneFunc()
	}()
}
