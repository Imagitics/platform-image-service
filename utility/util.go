package utility

import (
	"encoding/json"
	"github.com/nik/Imagitics/platform-s3-plugin/config"
	"os"
)

// Loads the configuration from the properties  file
func LoadConfiguration(file string) (*config.ConfigModel, error) {
	var config config.ConfigModel
	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return &config, nil
}
