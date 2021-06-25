package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
)

var wg sync.WaitGroup

func init() {
	log.SetFlags(log.Lshortfile)
	log.SetFlags(log.Ldate)
	log.SetFlags(log.Ltime)
}

func main() {
	config := sarama.NewConfig() // 1
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Version = sarama.V0_10_2_1
	consumer, err := sarama.NewConsumer([]string{"192.168.2.163:9092"}, config)
	if err != nil {
		log.Println("consumer connect error: ", err)
		return
	}
	fmt.Println("consumer connect success...")
	defer consumer.Close()

	partitions, err := consumer.Partitions("iot")
	if err != nil {
		fmt.Println("get partitions failed, err: ", err)
		return
	}

	for _, p := range partitions {
		partitionConsumer, err := consumer.ConsumePartition("iot", p, sarama.OffsetOldest)
		if err != nil {
			fmt.Println("partitionConsumer err:", err)
			continue
		}
		wg.Add(1)
		go func() {
			for m := range partitionConsumer.Messages() {
				fmt.Printf("key: %s,\ndata: %s,\noffset: %d\n", string(m.Key), string(m.Value), m.Offset)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
