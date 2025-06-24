package mqTT

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"net/url"
)

func Connect(ctx context.Context, u *url.URL, topic string, clientId string) (*autopaho.ConnectionManager, error) {
	cliCfg := mqttConfig(u, topic, clientId)
	return autopaho.NewConnection(ctx, cliCfg) // starts a process; will reconnect until context canceled
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
			// OnPublishReceived is a slice of functions that will be called when a message is received.
			// You can write the function(s) yourself or use the supplied Router
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
