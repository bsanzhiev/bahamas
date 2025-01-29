package config

type Config struct {
	CustomerServiceURL    string
	TransactionServiceURL string
	AuthSecret            string
	HTTP                  HTTP
}

type HTTP struct {
	Port string
}

func Load() (*Config, error) {
	return &Config{
			AuthSecret:            "secret",
			CustomerServiceURL:    "localhost:50051",
			TransactionServiceURL: "localhost:50052",
			HTTP: HTTP{
				Port: "8080",
			},
		},
		nil // Add error handling
}
