package service

import (
	"context"
	"encoding/json"
	"fmt"

	p2pHost "github.com/TheLazarusNetwork/erebrus-gateway/app/p2p-Node/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

const DiscoveryServiceTag = "erebrus"

type status struct {
	Status string
}

func Init(h host.Host, ctx context.Context) {
	ps := NewService(h, ctx)
	dht, err := p2pHost.NewDHT(ctx, h, nil)
	if err != nil {
		panic(err)
	}

	// Setup global peer discovery over DiscoveryServiceTag.
	p2pHost.Discover(ctx, h, dht)
	SubscribeTopics(ps, h, ctx)
}
func NewService(h host.Host, ctx context.Context) *pubsub.PubSub {
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}
	return ps
}

var Status_data []*Status
var StatusData map[string]*Status

func SubscribeTopics(ps *pubsub.PubSub, h host.Host, ctx context.Context) {

	topicString := "status"
	topic, err := ps.Join(DiscoveryServiceTag + "/" + topicString)
	if err != nil {
		panic(err)
	}
	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			// Block until we recieve a new message.
			msg, err := sub.Next(ctx)
			if err != nil {
				panic(err)
			}
			if msg.ReceivedFrom == h.ID() {
				continue
			}
			status := new(Status)
			if err := json.Unmarshal(msg.Data, status); err != nil {
				panic(err)
			}

			//Status_data = append(Status_data, status)
			StatusData[status.PublicKey] = status
			if err := topic.Publish(ctx, []byte("Gateway recieved the status")); err != nil {
				panic(err)
			}
			topic.EventHandler()
		}
	}()

	topicString2 := "client"
	topic2, err := ps.Join(DiscoveryServiceTag + "/" + topicString2)
	if err != nil {
		panic(err)
	}

	sub2, err := topic2.Subscribe()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			// Block until we recieve a new message.
			msg, err := sub2.Next(ctx)
			if err != nil {
				panic(err)
			}
			if msg.ReceivedFrom == h.ID() {
				continue
			}
			fmt.Printf("[%s] , status isz: %s", msg.ReceivedFrom, string(msg.Data))
			if err := topic2.Publish(ctx, []byte("heres a reply from client")); err != nil {
				panic(err)
			}
		}
	}()

}
