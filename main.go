package main

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const TOPIC = "mqtt/sample"

var msgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Topic: %s\n", msg.Topic())
	fmt.Printf("Message: %s\n", msg.Payload())
}

func main() {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://localhost:1883").
		SetClientID("mqtt-sample").
		SetUsername("admin").
		SetPassword("public")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); !token.WaitTimeout(15 * time.Second) {
		// token.Wait() && token.Error != nil {
		fmt.Printf("failed to Connect to broker. Error: %v\n", token.Error())
		panic(token.Error())
	}

	if token := client.Subscribe(TOPIC, 0, msgHandler); !token.WaitTimeout(15 * time.Second) {
		fmt.Printf("failed to subscribe to Topic: %s, Error: %v\n", TOPIC, token.Error())
		os.Exit(1)
	}

	if token := client.Publish(TOPIC, 0, false, "mqtt sample message --- end"); !token.WaitTimeout(15 * time.Second) {
		fmt.Printf("failed to publish to Topic: %s, Error: %v\n", TOPIC, token.Error())
		os.Exit(1)
	}

	if token := client.Unsubscribe(TOPIC); !token.WaitTimeout(15 * time.Second) {
		fmt.Printf("failed to unsubscribe to Topic: %s, Error: %v\n", TOPIC, token.Error())
		os.Exit(1)
	}

	client.Disconnect(250)

	time.Sleep(1 * time.Second)
}
