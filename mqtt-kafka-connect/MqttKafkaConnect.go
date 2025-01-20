package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

var (
	i = 0
)

type MqttKafkaConnect struct {
	config      Config
	mqttClient  mqtt.Client
	kafkaWriter *kafka.Writer
}

func NewMqttKafkaConnect(config Config) *MqttKafkaConnect {
	return &MqttKafkaConnect{
		config:      config,
		mqttClient:  CreateMqttClient(config.Mqtt),
		kafkaWriter: CreateKafkaWriter(config.Kafka),
	}
}

func (m *MqttKafkaConnect) Start() {
	m.SubscribeToMqtt()

	fmt.Printf("Listening to the incoming MQTT messages on topic: %s\n", m.config.Mqtt.Topic)

	select {}
}

func (m *MqttKafkaConnect) SubscribeToMqtt() {
	if token := m.mqttClient.Subscribe(m.config.Mqtt.Topic, byte(m.config.Mqtt.QoS), m.MessagePubHandler()); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Printf("Subscribed to topic %s\n", m.config.Mqtt.Topic)
}

func CreateMqttClient(config MqttConfig) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", config.Broker, config.Port))
	opts.SetClientID("go_mqtt_subscriber")
	// opts.SetUsername("samir")
	// opts.SetPassword("Eawa06Xmyfs9crbteLUG")
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func CreateKafkaWriter(config KafkaConfig) *kafka.Writer {
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{config.BrokerAddress},
		Topic:   config.Topic,
	})

	fmt.Printf("Created Producer\n")

	return kafkaWriter
}

func (m *MqttKafkaConnect) MessagePubHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		ctx := context.Background()

		payload := msg.Payload()
		fmt.Printf("Received message: %s from topic: %s\n", payload, msg.Topic())

		batteryTelemetry := CreateProtoFromPayload(payload)

		data, err := proto.Marshal(batteryTelemetry)
		if err != nil {
			fmt.Printf("cannot marshal proto message to binary: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Published Kafka message to topic: %s\n", m.config.Kafka.Topic)

		m.produce(ctx, data)
	}
}

func CreateProtoFromPayload(payload []byte) *BatteryTelemetry {
	parts := strings.Split(string(payload), "|")
	if len(parts) != 5 {
		fmt.Printf("Invalid payload\n")
		os.Exit(1)
	}

	values := make([]float32, 4)
	for i := 1; i < len(parts); i++ {
		value, err := strconv.ParseFloat(parts[i], 32)
		if err != nil {
			fmt.Printf("failed to parse float: %v", err)
			os.Exit(1)
		}
		values[i-1] = float32(value)
	}

	batteryTelemetry := &BatteryTelemetry{
		Id:          parts[0],
		Voltage:     values[0],
		Current:     values[1],
		Temperature: values[2],
		Charge:      values[3],
	}
	return batteryTelemetry
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func (m *MqttKafkaConnect) produce(ctx context.Context, message []byte) {
	err := m.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.Itoa(i)),
		Value: message,
	})
	if err != nil {
		panic("could not write message " + err.Error())
	}
	i++
}
