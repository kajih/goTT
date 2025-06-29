package mqtt

import (
	"errors"
	"fmt"
	mqconst "github.com/eclipse/paho.golang/packets"
	"github.com/eclipse/paho.golang/paho"
	"golang.org/x/net/context"
	"log"
	"net"
	"net/url"
)

type Server struct {
	Broker     string
	ClinetId   string
	Client     *paho.Client
	Connection net.Conn
}

func NewMqTT(broker, topic, clientId string) (*Server, error) {

	brokerUrl, err := url.Parse(broker)

	if err != nil {
		log.Fatal("Broker URL parse error")
	}

	conn, err := net.Dial(brokerUrl.Scheme, brokerUrl.Host)
	if err != nil {
		fmt.Printf("Failed to connect to broker: %v\n", err)
		return nil, err
	}

	router := paho.NewStandardRouter()
	router.RegisterHandler(topic, func(m *paho.Publish) {
		fmt.Printf("Received on %s: %s\n", m.Topic, m.Payload)
	})

	/* router.SetDefaultHandler(func(m *paho.Publish) {
		fmt.Printf("Received on %s: %s\n", m.Topic, m.Payload)
	}) */

	client := paho.NewClient(paho.ClientConfig{
		Router: router,
		Conn:   conn,
	})

	server := &Server{
		Broker:     broker,
		ClinetId:   clientId,
		Client:     client,
		Connection: conn,
	}

	return server, nil
}

// DisConnect
func (s *Server) DisConnect() {
	_ = s.Client.Disconnect(&paho.Disconnect{})
	_ = s.Connection.Close()
	fmt.Println("Disconnected.")
}

// Connect
func (s *Server) Connect(ctx context.Context) error {
	connAck, err := s.Client.Connect(ctx, &paho.Connect{
		ClientID:   s.ClinetId,
		KeepAlive:  20,
		CleanStart: true,
	})

	if err != nil || connAck.ReasonCode != mqconst.ConnackSuccess {
		fmt.Printf("Failed to connect (reason %d): %v\n", connAck.ReasonCode, err)
		return err
	}
	fmt.Println("Connected (MQTT 5).")
	return nil
}

// Publish
func (s *Server) Publish(ctx context.Context, topic string, payload []byte) error {

	if s.Client == nil {
		return errors.New("No Paho client initialized")
	}

	pubResp, err := s.Client.Publish(ctx, &paho.Publish{
		Topic:   topic,
		Payload: payload,
		QoS:     0,
	})

	if err != nil {
		fmt.Printf("Failed to publish: %v\n", err)
		return err
	} else {
		fmt.Printf("Published to topic %s: %s\n", topic, payload)
	}

	if pubResp != nil {
		fmt.Printf("Published to topic %s: %s\n", topic, pubResp)
	}
	return nil
}

// Subscribe
func (s *Server) Subscribe(ctx context.Context, topic string) error {
	subResp, err := s.Client.Subscribe(ctx, &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{Topic: topic, QoS: 0},
		},
	})
	if err != nil || subResp.Reasons[0] != 0 {
		fmt.Printf("Subscribe failed: %v, reason: %v\n", err, subResp.Reasons[0])
		return err
	}
	fmt.Printf("Subscribed to %s.\n", topic)
	return nil
}
