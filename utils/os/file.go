package os

import (
	"encoding/json"
	"fmt"
	"os"
)

var configPath = "/selected_windows.json"

func SaveSelectedWindows(selectedWindows []string) error {
	data, err := json.Marshal(selectedWindows)
	if err != nil {
		return fmt.Errorf("failed to encode selectedWindows: %w", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	err = os.WriteFile(currentDir+configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func LoadSelectedWindows() (selectedWindows []string, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return selectedWindows, fmt.Errorf("failed to get current directory: %w", err)
	}

	data, err := os.ReadFile(currentDir + configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return selectedWindows, nil
		}

		return selectedWindows, err
	}

	err = json.Unmarshal(data, &selectedWindows)
	if err != nil {
		println("failed to unmarshal selectedWindows: ", err)

		return selectedWindows, err
	}

	return selectedWindows, nil
}
