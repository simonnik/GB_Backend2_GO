package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/mediocregopher/radix/v3"
	"gocloud.dev/pubsub"
	//_ "gocloud.dev/pubsub/kafkapubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

var (
	sub      *pubsub.Subscription
	s        *radix.Pool
	timeout  = 10 * time.Second
	connFunc = func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
			radix.DialTimeout(timeout),
		)
	}
)

func subscription(ctx context.Context) (*pubsub.Subscription, error) {
	if sub != nil {
		return sub, nil
	}
	var err error
	//sub, err = pubsub.OpenSubscription(ctx, "kafka://process?topic=rates")

	sub, err = pubsub.OpenSubscription(ctx, "rabbit://myrates")
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func main() {
	ctx := context.Background()
	for {
		s, err := subscription(ctx)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		msg, err := s.Receive(context.Background())
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		log.Println(string(msg.Body))
		err = storage().Do(radix.Cmd(nil, "LPUSH", "result", string(msg.Body)))
		if err != nil {
			log.Println(err)
		}
		if rand.Float64() < .05 {
			_ = storage().Do(radix.Cmd(nil, "LTRIM", "result", "0", "9"))
		}
		msg.Ack()
	}
}

func storage() *radix.Pool {
	if s != nil {
		return s
	}
	var err error
	s, err = radix.NewPool("tcp", "redis:6379", 1, radix.PoolConnFunc(connFunc))
	if err != nil {
		panic(err)
	}
	return s
}
