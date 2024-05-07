package actions

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	gatewayTypes "github.com/bsanzhiev/bahamas/ms-gateway/types"
	"github.com/bsanzhiev/bahamas/ms-users/controllers"
	"log"
)

// Define user controller
var userController = controllers.UserController{
	Ctx: context.Background(),
	// TODO: How to pass or get DBPool
	DBPool: DBPool,
}

func UserList() {
	// Generate response
	var responseData = gatewayTypes.ResponseData{}
	users, err := userController.GetUsers()
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}
	responseData.Status = 200
	responseData.Message = "Success"
	responseData.Data = users

	// Send response to Kafka topic ('users_responses')
	responseTopic := "users_responses"
	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		log.Printf("Failed to marshal response data: %v", err)
		continue
	}
	producerMsg := sarama.ProducerMessage{
		Topic: responseTopic,
		Value: sarama.ByteEncoder(responseJSON),
	}

	if _, _, err := producer.SendMessage(&producerMsg); err != nil {
		log.Printf("Failed to send response message: %v", err)
		continue
	}
}
