package kafka

import (
	"chalurania/service/log"
	"github.com/Shopify/sarama"
)

func Consumer() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_1_0_0

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Error.Printf("sarama.NewConsumer err, message: %v\n", err)
		return
	}
	defer consumer.Close()

	partition, err := consumer.ConsumePartition("test", 0, sarama.OffsetOldest)
	if err != nil {
		log.Error.Printf("try create partition_consumer err, message: %v\n", err)
		return
	}
	defer partition.Close()

	for {
		select {
		case msg := <-partition.Messages():
			log.Info.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s \n", msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value))
		case err := <-partition.Errors():
			log.Error.Printf("err :%v \n", err)
		}
	}
}