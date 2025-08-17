package utils

import (
	"encoding/json"
	"fmt"
	"html-cli/constants"
	"html-cli/types"
	"os"
	"path/filepath"
)

func loadDefaults() {
	if constants.Config.Dev.Port == 0 {
		constants.Config.Dev.Port = 8080
	}

	if constants.Config.Build.Directory == "" {
		constants.Config.Build.Directory = "build"
	}
}

func LoadConfig(rootPath string) error {
	configFile, err := os.ReadFile(filepath.Join(rootPath, "html-cli.json"))
	if err != nil {
		return err
	}

	var config types.Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		return fmt.Errorf("failed to decode config: %s", err)
	}

	constants.Config = config
	loadDefaults()

	return nil
}
