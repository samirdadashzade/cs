# Project Setup and Execution Guide

## Docker Compose Setup
This project uses Docker Compose to set up and run the necessary services. The `docker-compose.yml` file includes the following services:
- **Kafka Broker**: A Kafka broker for message streaming.
- **Zookeeper**: A Zookeeper instance for managing the Kafka broker.
- **Mosquitto**: An MQTT broker for handling MQTT messages.

### Commands to Manage Docker Compose
- Start the services:
```sh
    docker-compose up -d
```

- Stop the services:
```sh
    docker-compose down
```

## Project Descriptions

### Mosquitto
Mosquitto is an open-source MQTT broker that implements the MQTT protocol. It is used to handle MQTT messages in this project.

### MQTT Publisher
The MQTT Publisher is a project that publishes messages to the Mosquitto MQTT broker. It simulates the generation of telemetry data and sends it to the MQTT broker.

### MQTT-Kafka Connect
The MQTT-Kafka Connect project bridges the MQTT broker and the Kafka broker. It subscribes to the MQTT topics and republishes the messages to Kafka topics.

### Kafka Consumer
The Kafka Consumer project consumes messages from the Kafka broker. It reads the telemetry data from the Kafka topics and processes it accordingly.

## Running the Projects Locally

### Step 1: Run Docker Compose
Start the Docker Compose services by running the following command:
```sh
    docker-compose up -d
```

### Step 2: Build Each Project
Navigate to each project directory and build the executables.

#### Build Kafka Consumer
```sh
    cd kafka-consumer
    go build -o kafka-consumer
```

#### Build MQTT-Kafka Connect
```sh
    cd mqtt-kafka-connect
    go build -o mqtt-kafka-connect
```

#### Build MQTT Publisher
```sh
    cd mqtt-publisher
    go build -o mqtt-publisher
```

### Step 3: Run the Executables in Order

#### Run Kafka Consumer
```sh
    ./kafka-consumer
```

#### Run MQTT-Kafka Connect
```sh
    ./mqtt-kafka-connect
```

#### Run MQTT Publisher
```sh
    ./mqtt-publisher
```

## Additional Commands

### Create Kafka Topics
```sh
    docker exec -it kafka-broker-1 kafka-topics --create --topic my-topic --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1
```

### Verify the Creation of Kafka Topics
```sh
    docker exec -it kafka-broker-1 kafka-topics --bootstrap-server localhost:9092 --list
```

### Start the Consumer for the Topic
```sh
    docker exec -it kafka-broker-1 kafka-console-consumer --topic my-topic --bootstrap-server localhost:9092
```

### Start the Producer for the Topic and Publish Messages
```sh
    docker exec -it kafka-broker-1 kafka-console-producer --topic my-topic --bootstrap-server localhost:9092
```

### Protoc File Generator
```sh
    protoc --go_out=. --go_opt=paths=source_relative --proto_path=. battery_telemetry.proto
```