package main

func main() {
	config := Config{
		Mqtt: MqttConfig{
			QoS:    1,
			Topic:  "mqtt/bess",
			Broker: "localhost",
			Port:   "1883",
		},
		Kafka: KafkaConfig{
			BrokerAddress: "localhost:19092",
			Topic:         "kafka-bess",
		},
	}

	mqttKafkaConnector := NewMqttKafkaConnect(config)
	mqttKafkaConnector.Start()
}
