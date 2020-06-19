package config

type ConfigModel struct {
	Platform_S3_URL string `json:"platform_s3_service_url"`
	Cassandra       struct {
		Host        string `json:"host"`
		Port        string `json:"port"`
		User        string `json:"user"`
		Password    string `json:"password"`
		SSLCertPath string `json:"cert_file_path"`
		Consistency string `json:"consistency"`
		Keyspace    string `json:"keyspace"`
	}
	Logger *Logger `json:"logger"`
}

type Logger struct {
	LoggerFilePath string `json:"logger_file_path"`
	LoggerFileName string `json:"logger_file_name"`
}
