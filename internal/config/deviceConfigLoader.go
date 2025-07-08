package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func DeviceConfigLoader(configJSONPath string) (config DevicesConfigJSON, err error) {
	abspath, errAbsPath := filepath.Abs(configJSONPath)
	if errAbsPath != nil {
		err = fmt.Errorf("error while trying to get absolute path of %s : %w", configJSONPath, errAbsPath)
		return
	}
	bytes, errReadFile := os.ReadFile(abspath)
	if errReadFile != nil {
		err = fmt.Errorf("error while reading specified config json file at %s : %w", abspath, errReadFile)
		return
	}
	var tmp DevicesConfigJSON
	errUnmarshal := json.Unmarshal(bytes, &tmp)
	if errUnmarshal != nil {
		err = fmt.Errorf("error while json unmarshal for config json at %s : %w", abspath, errUnmarshal)
		return
	}

	config = tmp
	err = nil
	return
}
