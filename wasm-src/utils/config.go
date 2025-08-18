package utils

import (
	"encoding/json"
	"fmt"
	"html-cli/constants"
	"html-cli/types"
	"path/filepath"

	"github.com/spf13/afero"
)

func loadDefaults() {
	if constants.Config.Dev.Port == 0 {
		constants.Config.Dev.Port = constants.DefaultConfig.Dev.Port
	}

	if constants.Config.Build.Directory == "" {
		constants.Config.Build.Directory = constants.DefaultConfig.Build.Directory
	}
}

func LoadConfig(fs afero.Fs, rootPath string) error {
	configFile, err := afero.ReadFile(fs, filepath.Join(rootPath, "html-cli.json"))
	if err != nil {
		loadDefaults()
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
