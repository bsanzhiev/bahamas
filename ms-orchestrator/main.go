package main

import (
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	nc.Subscribe("start.saga", func(m *nats.Msg) {
		log.Println("Saga started for: ", string(m.Data))
		// Implement saga logic
		err := nc.Publish("ms-users.create", m.Data)
		if err != nil {
			return
		}
	})
	select {}
}
