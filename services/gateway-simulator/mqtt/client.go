package mqtt

import (
	"fmt"
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client paho.Client
	qos    byte
}

func NewClient(broker, clientID string, qos byte) (*Client, error) {
	opts := paho.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetConnectionLostHandler(func(c paho.Client, err error) {
			log.Printf("[MQTT] Connection lost: %v", err)
		}).
		SetOnConnectHandler(func(c paho.Client) {
			log.Printf("[MQTT] Connected: %s", clientID)
		})

	client := paho.NewClient(opts)

	token := client.Connect()
	if !token.WaitTimeout(10 * time.Second) {
		return nil, fmt.Errorf("connection timeout")
	}
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	return &Client{client: client, qos: qos}, nil
}

func (c *Client) Publish(topic string, payload []byte) error {
	token := c.client.Publish(topic, c.qos, false, payload)
	token.Wait()
	return token.Error()
}

func (c *Client) Disconnect() {
	c.client.Disconnect(1000)
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}
