package mqtt

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"net/url"
	"strconv"
)

type Server struct {
	Broker string
	Paho   *autopaho.ConnectionManager
}

func NewMqTT(ctx context.Context, broker, topic, clientId string) (*Server, error) {

	brokerUrl, err := url.Parse(broker)
	if err != nil {
		return nil, err
	}

	cliCfg := mqttConfig(brokerUrl, topic, clientId)
	cm, _ := autopaho.NewConnection(ctx, cliCfg)

	/*
		cm.AddOnPublishReceived(
			func(pr autopaho.PublishReceived) (bool, error) {
				fmt.Printf("received (two) message on topic %s; body: %s (retain: %t)\n", pr.Packet.Topic, pr.Packet.Payload, pr.Packet.Retain)
				return true, nil
			},
		)
	*/

	return &Server{
		Broker: broker,
		Paho:   cm,
	}, nil
}

func mqttConfig(u *url.URL, topic string, clientId string) autopaho.ClientConfig {
	cliCfg := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		KeepAlive:                     20,
		CleanStartOnInitialConnection: false,
		SessionExpiryInterval:         60,
		OnConnectionUp:                onConnect(topic),
		OnConnectError:                errorConnect(),
		ClientConfig: paho.ClientConfig{
			ClientID: clientId,
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				func(pr paho.PublishReceived) (bool, error) {
					fmt.Printf("received message on topic %s; body: %s (retain: %t)\n", pr.Packet.Topic, pr.Packet.Payload, pr.Packet.Retain)
					return true, nil
				}},
			OnClientError:      clientError(),
			OnServerDisconnect: clientDisconnect(),
		},
	}
	return cliCfg
}

func clientDisconnect() func(d *paho.Disconnect) {
	return func(d *paho.Disconnect) {
		if d.Properties != nil {
			fmt.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
		} else {
			fmt.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
		}
	}
}

func clientError() func(err error) {
	return func(err error) { fmt.Printf("client error: %s\n", err) }
}

func errorConnect() func(err error) {
	return func(err error) { fmt.Printf("error whilst attempting connection: %s\n", err) }
}

func onConnect(topic string) func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
	return func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
		fmt.Println("mqtt connection up")
		// Subscribing in the OnConnectionUp callback is recommended (ensures the subscription is reestablished if
		// the connection drops)
		if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
			Subscriptions: []paho.SubscribeOptions{
				{Topic: topic, QoS: 1},
			},
		}); err != nil {
			fmt.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
		}
		fmt.Println("mqtt subscription made")
	}
}

func CreateMessage(topic string, msgCount int) *paho.Publish {
	return &paho.Publish{
		QoS:     1,
		Topic:   topic,
		Payload: []byte("TestMessage: " + strconv.Itoa(msgCount)),
	}
}
