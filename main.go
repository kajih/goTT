package main

import (
	"context"
	"fmt"
	"goTT/mqtt"
	"goTT/web"
	"goTT/web/sse"
	"os"
	"time"
)

func main() {

	broker := os.Getenv("MQ_BROKER")
	clientId := os.Getenv("MQ_CLIENT_ID")
	topic := os.Getenv("MQ_TOPIC")

	/*
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	*/

	mqttConnection, err := mqtt.NewMqTT(broker, topic, clientId)
	if err != nil {
		panic(err)
	}

	// Server loop
	msgCount := 0
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	r := web.NewRouter()
	webErr := make(chan error, 1)
	go func() {
		webErr <- r.Run("localhost:8080") // will send nil or error
	}()

	// Start the sse change engine
	go sse.SimulateStateChanges()

	/*
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
	*/

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = mqttConnection.Connect(ctx)

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = mqttConnection.Subscribe(ctx, topic)

	for {
		select {
		case <-ticker.C:
			msgCount++

			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
			payload := []byte(fmt.Sprintf("Message %d", msgCount))
			if err = mqttConnection.Publish(ctx, topic, payload); err != nil {
				if ctx.Err() == nil {
					panic(err) // Publish will exit when context canceled or if something went wrong
				}
			}
			continue

		case <-webErr:
		case <-ctx.Done():
			mqttConnection.DisConnect()
		}
		break
	}

	fmt.Println("signal caught - exiting")
	//<-conn.Paho.Done() // Wait for a clean shutdown (cancelling the context triggered the shutdown)

	// Wait a bit to receive messages
	//time.Sleep(10 * time.Second)
	//mqttConnection.DisConnect()

}
