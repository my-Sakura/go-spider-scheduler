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

	nc.Subscribe("tasks", func(m *nats.Msg) {

		fmt.Printf("Received a message: %s\n", string(m.Data))

	})
	select {}
}
