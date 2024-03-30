package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"permnews/internal/storage/filters"
	"permnews/internal/storage/models"
	"permnews/internal/storage/repository"
)

type EventRepo interface {
	Create(event *models.Event) error
	GetAll(filter *filters.EventFilter) ([]*models.Event, error)
}

type Consumer struct {
	Urls         []string
	Topic        string
	partConsumer *sarama.PartitionConsumer
	repo         EventRepo
}

func NewConsumer(urls []string, topic string, storage repository.Storage) *Consumer {
	c, err := sarama.NewConsumer(urls, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	partConsumer, err := c.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
		return nil
	}

	return &Consumer{
		Urls:         urls,
		Topic:        topic,
		partConsumer: &partConsumer,
		repo:         repository.NewEventsRepository(storage),
	}
}

func (c *Consumer) StartListening() {
	go func() {
		defer (*c.partConsumer).Close()
		for {
			select {
			case e := <-(*c.partConsumer).Errors():
				log.Printf("Af: %v", e.Error())
				return
			// Чтение сообщения из Kafka
			case msg, ok := <-(*c.partConsumer).Messages():
				if !ok {
					log.Fatal("Channel closed, exiting goroutine")
					return
				}
				event := make([]models.Event, 0)
				err := json.Unmarshal(msg.Value, &event)
				if err != nil {
					log.Fatal(err, string(msg.Value))
					return
				}
				log.Printf("%v", len(event))
				for _, ev := range event {
					all, err := c.repo.GetAll(&filters.EventFilter{Source: ev.Source, Title: ev.Title})
					if err != nil {
						log.Println(err)
						continue
					}
					if len(all) > 0 {
						flag := false
						for _, e := range all {
							if e.Notice == ev.Notice {
								flag = true
							}
						}
						if flag {
							log.Println("event already exists")
							break
						}
					}
					err = c.repo.Create(&ev)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}()
}
