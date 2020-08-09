package config

type ConfigModel struct {
	Platform_S3_URL string     `json:"platform_s3_service_url"`
	AWS             *AWS       `json:"aws"`
	Dynamodb        *Dynamodb  `json:"dynamodb"`
	Logger          *Logger    `json:"logger"`
	Messaging       *Messaging `json:"messaging"`
	Cassandra       struct {
		Host        string `json:"host"`
		Port        string `json:"port"`
		User        string `json:"user"`
		Password    string `json:"password"`
		SSLCertPath string `json:"cert_file_path"`
		Consistency string `json:"consistency"`
		Keyspace    string `json:"keyspace"`
	}
}

type Logger struct {
	LoggerFilePath string `json:"logger_file_path"`
	LoggerFileName string `json:"logger_file_name"`
}

type Messaging struct {
	Region string `json:"region"`
}

type AWS struct {
	AccessKey  string `json:"aws_access_key"`
	SecretKey  string `json:"aws_secret_key"`
	Region     string `json:"region"`
	PathPrefix string `json:"path_prefix"`
}

type Dynamodb struct {
	Endpoint string `json:"endpoint"`
}
