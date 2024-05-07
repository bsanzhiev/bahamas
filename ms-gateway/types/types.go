package types

// RequestData Тип для запросов
type RequestData struct {
	Service string                 `json:"service"` // Имя сервиса
	Action  string                 `json:"action"`  // Имя операции
	Data    map[string]interface{} `json:"data"`    // Объект запроса
}

// ResponseData Тип для ответов
type ResponseData struct {
	Status  int         `json:"status"`  // Код ответа
	Message string      `json:"message"` // Сообщение об ошибке или результате
	Data    interface{} `json:"data"`    // Объект ответа
}
