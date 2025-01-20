package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	MqttQoS    = 1
	MqttTopic  = "mqtt/bess"
	MqttBroker = "localhost"
	MqttPort   = 1883
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func publish(client mqtt.Client) {
	num := 1000
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("%d|%d|%d|%d|%d", i+1, i+2, i+3, i+4, i+5)
		token := client.Publish(MqttTopic, MqttQoS, false, text)
		fmt.Printf("Published message: %s to topic: %s\n", text, MqttTopic)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", MqttBroker, MqttPort))
	opts.SetClientID("go_mqtt_publisher")
	// opts.SetUsername("samir")
	// opts.SetPassword("Eawa06Xmyfs9crbteLUG")
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	publish(client)

	client.Disconnect(250)
}
