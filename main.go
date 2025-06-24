package main

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.golang/paho"
	"goTT/mqTT"
	"goTT/webserve"
	"net/url"
	"os"
	"os/signal"
	"strconv"
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

	conn, err := mqTT.Connect(ctx, brokerUrl, topic, clientId)
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

	webErr := webserve.Serve()

	for {
		select {
		case <-ticker.C:
			msgCount++
			// Publish a test message (use PublishViaQueue if you don't want to wait for a response)
			if _, err = conn.Publish(ctx, &paho.Publish{
				QoS:     1,
				Topic:   topic,
				Payload: []byte("TestMessage: " + strconv.Itoa(msgCount)),
			}); err != nil {
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
