package main

import (
	"context"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

const (
	BrokerAddress = "localhost:19092"
	Topic         = "kafka-bess"
)

func main() {
	consume(context.Background())
}

func consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{BrokerAddress},
		Topic:   Topic,
		GroupID: "my-group", // for deduplication
	})

	fmt.Printf("start consuming the messages from topic: %s\n", Topic)

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}

		batteryTelemetry := BatteryTelemetry{}
		err = proto.Unmarshal(msg.Value, &batteryTelemetry)
		if err != nil {
			fmt.Printf("cannot marshal proto message to binary: %v", err)
			os.Exit(1)
		}

		fmt.Printf("received battery telemetry: %v\n", &batteryTelemetry)
	}
}
