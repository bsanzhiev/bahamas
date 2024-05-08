package actions

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	gatewayTypes "github.com/bsanzhiev/bahamas/ms-gateway/types"
	"github.com/bsanzhiev/bahamas/ms-users/controllers"
	"github.com/bsanzhiev/bahamas/ms-users/types"
	"log"
)

func HandleAction(action string, data interface{}, uc controllers.UserController, producer sarama.SyncProducer) {

	switch action {
	case "user_list":
		UserList(uc, producer)
	case "user_by_id":
		UserByID(data, uc, producer)
	case "user_create":
		UserCreate(data, uc, producer)
	case "user_update":
		UserUpdate(data, uc, producer)
	default:
		DefaultAction(uc, producer)
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
	// Get user ID =====================
	var userData types.UserRequestData
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return
	}

	err = json.Unmarshal(dataBytes, &userData)
	if err != nil {
		log.Printf("Error unmarshaling data to struct: %v", err)
		return
	}

	userID := userData.ID
	if userID == 0 {
		log.Printf("Invalid or missing user ID")
		return
	}

	// Generate response
	var responseData = gatewayTypes.ResponseData{}
	users, err := uc.UserByID(userID)
	if err != nil {
		responseData.Status = 404
		responseData.Message = fmt.Sprintf("Failed to get user: %v", err)
		responseData.Data = ""
		log.Printf("Failed to get user: %v", err)
	} else {
		responseData.Status = 200
		responseData.Message = "Success"
		responseData.Data = users
	}
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

func UserCreate(data interface{}, uc controllers.UserController, producer sarama.SyncProducer) {
	// Get user ID =====================
	var userData types.UserRequestData
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return
	}

	err = json.Unmarshal(dataBytes, &userData)
	if err != nil {
		log.Printf("Error unmarshaling data to struct: %v", err)
		return
	}

	// Generate response
	var responseData = gatewayTypes.ResponseData{}
	err = uc.UserCreate(userData)
	if err != nil {
		responseData.Status = 404
		responseData.Message = fmt.Sprintf("Failed to create user: %v", err)
		responseData.Data = ""
		log.Printf("Failed to get user: %v", err)
	} else {
		responseData.Status = 200
		responseData.Message = "Success"
		responseData.Data = ""
	}
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

func UserUpdate(data interface{}, uc controllers.UserController, producer sarama.SyncProducer) {
	// Get user ID =====================
	var userData types.UserRequestData
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return
	}

	err = json.Unmarshal(dataBytes, &userData)
	if err != nil {
		log.Printf("Error unmarshaling data to struct: %v", err)
		return
	}

	userID := userData.ID
	if userID == 0 {
		log.Printf("Invalid or missing user ID")
		return
	}

	// Generate response
	var responseData = gatewayTypes.ResponseData{}
	err = uc.UserUpdate(userID, userData)
	if err != nil {
		responseData.Status = 404
		responseData.Message = fmt.Sprintf("Failed to create user: %v", err)
		responseData.Data = ""
		log.Printf("Failed to get user: %v", err)
	} else {
		responseData.Status = 200
		responseData.Message = "Success"
		responseData.Data = ""
	}
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

func DefaultAction(uc controllers.UserController, producer sarama.SyncProducer) {
	// Generate response
	var responseData = gatewayTypes.ResponseData{}
	responseData.Status = 404
	responseData.Message = "Unknown action"
	responseData.Data = ""

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
