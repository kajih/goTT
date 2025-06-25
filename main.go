package main

import (
	"context"
	"fmt"
	"goTT/mqtt"
	"goTT/web"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	broker := os.Getenv("MQ_BROKER")
	clientId := os.Getenv("MQ_CLIENT_ID")
	topic := os.Getenv("MQ_TOPIC")

	brokerUrl, err := url.Parse(broker)
	if err != nil {
		panic(err)
	}

	conn, err := mqtt.Connect(ctx, brokerUrl, topic, clientId)
	if err != nil {
		panic(err)
	}

	// Wait for the connection to come up
	if err = conn.AwaitConnection(ctx); err != nil {
		panic(err)
	}

	msgCount := 0
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	r := web.NewRouter()
	webErr := make(chan error, 1)
	go func() {
		webErr <- r.Run("localhost:8080") // will send nil or error
	}()

	for {
		select {
		case <-ticker.C:
			msgCount++
			// Publish a test message (use PublishViaQueue if you don't want to wait for a response)
			if _, err = conn.Publish(ctx, mqtt.CreateMessage(topic, msgCount)); err != nil {
				if ctx.Err() == nil {
					panic(err) // Publish will exit when context canceled or if something went wrong
				}
			}
			continue

		case <-webErr:
		case <-ctx.Done():
		}
		break
	}

	fmt.Println("signal caught - exiting")
	<-conn.Done() // Wait for a clean shutdown (cancelling the context triggered the shutdown)
}
