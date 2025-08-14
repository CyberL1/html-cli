package utils

import (
	"encoding/json"
	"fmt"
	"html-cli/constants"
	"html-cli/types"
	"os"
)

func LoadConfig() error {
	configFile, err := os.Open("html-cli.json")
	if err != nil {
		return nil
	}

	defer configFile.Close()
	decoder := json.NewDecoder(configFile)

	var config types.Config
	if err := decoder.Decode(&config); err != nil {
		return fmt.Errorf("failed to decode html config: %s", err)
	}

	constants.Config = config
	return nil
}
