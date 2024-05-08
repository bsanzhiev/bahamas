package actions

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	gatewayTypes "github.com/bsanzhiev/bahamas/ms-gateway/types"
	"github.com/bsanzhiev/bahamas/ms-users/controllers"
	"log"
)

func HandleAction(action string, data interface{}, uc controllers.UserController, producer sarama.SyncProducer) {
	switch action {
	case "user_list":
		UserList(uc, producer)
	case "user_by_id":
		UserByID(data, uc, producer)
	default:
		fmt.Println("Unknown action")

	}
}

func UserList(uc controllers.UserController, producer sarama.SyncProducer) {
	// Generate response
	var responseData = gatewayTypes.ResponseData{}
	users, err := uc.GetUsers()
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
		return
	}
	producerMsg := &sarama.ProducerMessage{
		Topic: responseTopic,
		Value: sarama.ByteEncoder(responseJSON),
	}

	if _, _, err := producer.SendMessage(producerMsg); err != nil {
		log.Printf("Failed to send response message: %v", err)
	}
}

func UserByID(data interface{}, uc controllers.UserController, producer sarama.SyncProducer) {
	// Generate response
	var responseData = gatewayTypes.ResponseData{}
	users, err := uc.UserByID(data)
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
		return
	}
	producerMsg := &sarama.ProducerMessage{
		Topic: responseTopic,
		Value: sarama.ByteEncoder(responseJSON),
	}

	if _, _, err := producer.SendMessage(producerMsg); err != nil {
		log.Printf("Failed to send response message: %v", err)
	}
}
