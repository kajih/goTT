package sse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

/*
func Handler(context *gin.Context) {
	context.String(200, "Hello, World!")
}
*/

type Subscriber struct {
	channel chan string
}

var (
	subscribers = make(map[*Subscriber]bool)
	mutex       sync.Mutex
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/subscribe", handler)
}

// Handler SSE endpoint using Gin
func handler(c *gin.Context) {
	// Setup headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // dev CORS

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "Streaming unsupported")
		return
	}

	// Create and register subscriber
	subscriber := &Subscriber{
		channel: make(chan string),
	}

	mutex.Lock()
	subscribers[subscriber] = true
	mutex.Unlock()

	// Cleanup on disconnect
	defer func() {
		mutex.Lock()
		delete(subscribers, subscriber)
		mutex.Unlock()
		close(subscriber.channel)
	}()

	// Send initial ping (to ensure browser sees connection is live)
	fmt.Fprintf(c.Writer, "data: connected\n\n")
	flusher.Flush()

	// Stream updates
	for {
		select {
		case msg, ok := <-subscriber.channel:
			if !ok {
				return
			}
			fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
			flusher.Flush()
		case <-c.Request.Context().Done():
			return // client disconnected
		}
	}
}

// Push message to all connected clients
func broadcastUpdate(message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for sub := range subscribers {
		select {
		case sub.channel <- message:
		default:
			// Drop dead connections
			delete(subscribers, sub)
			close(sub.channel)
		}
	}
}

// SimulateStateChanges Simulate background state changes
func SimulateStateChanges() {
	for {
		time.Sleep(time.Duration(rand.Intn(6)+4) * time.Second)
		state := fmt.Sprintf("updated at %s", time.Now().Format(time.RFC3339))
		log.Println("Broadcasting:", state)
		broadcastUpdate(state)
	}
}
