/*
Copyright 2024 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	daprd "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
)

func main() {
	client, err := daprd.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(">>Created client\n")

	// Another method of streaming subscriptions, this time for the topic "sendorder".
	// The given subscription handler is called when a message is received.
	// The  returned `stop` function is used to stop the subscription and close the connection.
	stop, err := client.SubscribeWithHandler(context.Background(),
		daprd.SubscriptionOptions{
			PubsubName: "messages",
			Topic:      "sendorder",
		},
		eventHandler,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Streaming subscription for topic "neworder" on pubsub component "messages".
	// The returned `sub` object is used to receive messages.
	// `sub` must be closed once it's no longer needed.
	sub, err := client.Subscribe(context.Background(), daprd.SubscriptionOptions{
		PubsubName: "messages",
		Topic:      "neworder",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(">>Created subscription\n")

	for i := 0; i < 3; i++ {
		msg, err := sub.Receive()
		if err != nil {
			log.Fatalf("error receiving message: %v", err)
		}
		log.Printf(">>Received message\n")
		log.Printf("event - PubsubName: %s, Topic: %s, ID: %s, Data: %s\n", msg.PubsubName, msg.Topic, msg.ID, msg.RawData)

		if err := msg.Success(); err != nil {
			log.Fatalf("error sending message success: %v", err)
		}
	}

	time.Sleep(time.Second * 5)

	if err := errors.Join(stop(), sub.Close()); err != nil {
		log.Fatal(err)
	}
}

func eventHandler(e *common.TopicEvent) common.SubscriptionResponseStatus {
	log.Printf(">>Received message\n")
	log.Printf("event - PubsubName: %s, Topic: %s, ID: %s, Data: %s\n", e.PubsubName, e.Topic, e.ID, e.Data)
	return common.SubscriptionResponseStatusSuccess
}