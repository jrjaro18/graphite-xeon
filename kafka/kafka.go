package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"revx/graph"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Kafka struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
}

type UserAction struct {
	UserId    primitive.ObjectID `json:"userId"` // user id who performed the action
	PostId    primitive.ObjectID `json:"postId"` // post id on which action is performed
	Action    int8               `json:"action"` // -1 for dislike, 0 for no action, 1 for like
	TimeSpent uint8              `json:"timeSpent"` // time spent on the post in seconds
}

type PostAction struct {
	PostId   primitive.ObjectID `json:"postId"` // post id on which action is performed
	Features []string           `json:"features"` // features of the post
}

func Initialize(topic string) *Kafka {
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		Balancer: &kafka.RoundRobin{},
	})
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     topic,
		Partition: 0,
	})
	return &Kafka{
		Writer: kafkaWriter,
		Reader: kafkaReader,
	}
}

func (k *Kafka) Close() error {
	err := k.Writer.Close()
	if err != nil {
		return err
	}
	err = k.Reader.Close()
	if err != nil {
		return err
	}
	return nil
}


func (k* Kafka) UserActionConsumer(graph *graph.Graph) {
	for {
		message, err := k.Reader.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			continue
		}
		userAction := UserAction{}
		if message.Value == nil {
			continue
		}
		err = json.Unmarshal(message.Value, &userAction)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// logic to update the graph
		graph.UpdateGraphFromUserAction(userAction.UserId, userAction.PostId, userAction.Action, userAction.TimeSpent)
		// will come here
	}
}

func (k* Kafka) PostActionConsumer(graph *graph.Graph) {
	for {
		message, err := k.Reader.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			continue
		}
		if message.Value == nil {
			continue
		}
		postAction := PostAction{}
		err = json.Unmarshal(message.Value, &postAction)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// logic to update the graph
		graph.UpdateGraphFromPostAction(postAction.PostId, postAction.Features)
		// will come here
	}
}