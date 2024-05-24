package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

// Account Creation
// 1. Account Service receive action for create account
// 2. Orchestrator create saga
// 3. Making transactions
// 4. Commit if ok, revert if errors

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting NATS: %v", err)
	}
	defer nc.Close()

	// Subscribe to start.saga to initiate saga
	nc.Subscribe("start.saga", func(m *nats.Msg) {
		log.Println("Saga started for: ", string(m.Data))

		// Step 1: Request a user ID from the users service
		msg, err := nc.Request("ms-users.getID", m.Data, 10*time.Second)
		if err != nil {
			log.Printf("Error getting user ID: %v", err)
			return
		}
		userID := string(msg.Data)
		log.Printf("Received user ID: %v", userID)

		// Step 2: Create a bank account
		err = nc.Publish("ms-accounts.createAccount", []byte(userID))
		if err != nil {
			log.Printf("Error creating bank account: %v", err)
			return
		}
		log.Printf("Bank account creation requested for user ID: %s", userID)

		// Optionally, listen for a response from the account service
		nc.Subscribe("ms-account.createAccount.response", func(m *nats.Msg) {
			log.Println("Bank account creation: ", string(m.Data))
		})
	})
	select {}
}
