package main

import (
	"log"

	gateway "github.com/bsanzhiev/bahamas/ms-gateway"
)

func main() {
	gateway.StartGateway()
	log.Println("Server started")
}
