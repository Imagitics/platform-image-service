package config

type ConfigModel struct {
	Google_API_METAMODEL struct {
		AccessKey  string `json:"aws_access_key"`
		SecretKey  string `json:"aws_secret_key"`
		Region     string `json:"region"`
		PathPrefix string `json:"path_prefix"`
	} `json:"s3"`
	Cassandra struct {
		Host        string `json:"host"`
		Port        string `json:"port"`
		User        string `json:"user"`
		Password    string `json:"password"`
		SSLCertPath string `json:"cert_file_path"`
		Consistency string `json:"consistency"`
		Keyspace    string `json:"keyspace"`
	}
}
