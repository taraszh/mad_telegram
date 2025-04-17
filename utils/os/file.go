package os

import (
	"encoding/json"
	"fmt"
	"os"
)

func saveSelectedWindows(filePath string, selectedWindows []string) error {
	data, err := json.Marshal(selectedWindows)
	if err != nil {
		return fmt.Errorf("failed to marshal selectedWindows: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func loadSelectedWindows(filePath string) (selectedWindows []string, err error) {
	data, err := os.ReadFile(filePath)
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
