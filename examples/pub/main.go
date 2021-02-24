package main

import (
	"fmt"

	nats "github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return
	}

	for i := 0; i < 10000; i++ {
		msg := fmt.Sprintf("hello sub ----------- %d", i)
		err = nc.Publish("tasks", []byte(msg))
	}

	nc.Flush()
}
