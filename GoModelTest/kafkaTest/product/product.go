package main

import (
	"github.com/Shopify/sarama"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// log.SetFlags(log.LUTC)
	SyncProducerExample()
	// AsyncProducerExample()
}

func SyncProducerExample() {
	config := sarama.NewConfig() // 1
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Version = sarama.V0_10_2_1
	producer, err := sarama.NewSyncProducer([]string{"192.168.2.163:9092"}, config)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println(err)
		}
	}()

	msg := &sarama.ProducerMessage{Topic: "iot", Value: sarama.StringEncoder("testing 123")}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("FAILED to send message: %s\n", err)
	} else {
		log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	}
}

func AsyncProducerExample() {
	config := sarama.NewConfig() // 1
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Version = sarama.V0_10_2_1
	client, err := sarama.NewClient([]string{"192.168.2.163:9092"}, config) // 2
	if err != nil {
		panic(err)
	}
	defer client.Close()
	producer, err := sarama.NewAsyncProducerFromClient(client) // 3
	if err != nil {
		panic(err)
	}
	defer producer.AsyncClose()

	message := &sarama.ProducerMessage{Topic: "iot", Value: sarama.StringEncoder("testing 123")} // 4
	producer.Input() <- message
	log.Println("over")
}
