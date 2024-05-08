package types

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type UserRequestData struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type IncomingData struct {
	Service string                 `json:"service"` // Имя сервиса
	Action  string                 `json:"action"`  // Имя операции
	Data    map[string]interface{} `json:"data"`    // Объект запроса
}

type OutgoingData struct {
	Status  int         `json:"status"`  // Код ответа
	Message string      `json:"message"` // Сообщение об ошибке или результате
	Data    interface{} `json:"data"`    // Объект ответа
}

type UserCreate struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
}
