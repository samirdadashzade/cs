package main

type Config struct {
	Mqtt  MqttConfig
	Kafka KafkaConfig
}

type MqttConfig struct {
	QoS    int
	Topic  string
	Broker string
	Port   string
}

type KafkaConfig struct {
	BrokerAddress string
	Topic         string
}
