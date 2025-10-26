package application

import (
	"dinky/internal/application/settingstype"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sedwards2009/femto"
	"github.com/sedwards2009/femto/runtime"
)

func LoadUserSettings() settingstype.Settings {
	settingsPath := userSettingsPath()
	return loadSettings(settingsPath)
}

func userSettingsPath() string {
	userSettingsDir := userSettingsDirPath()
	return filepath.Join(userSettingsDir, "settings.json")
}

func userSettingsDirPath() string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(userConfigDir, "dinky")
}

// Load the settings JSON from the settings path file and return them.
func loadSettings(settingsPath string) settingstype.Settings {
	// Read the settings file
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return settingstype.DefaultSettings()
	}
	settings := settingstype.DefaultSettings()
	if err := json.Unmarshal(data, &settings); err != nil {
		return settingstype.DefaultSettings()
	}

	settings.TabSize = cleanTabSize(settings.TabSize)
	settings.TabCharacter = cleanTabCharacter(settings.TabCharacter)
	settings.ColorScheme = cleanColorSchemeName(settings.ColorScheme)

	return settings
}

// Function to save the settings struct to a JSON file
func SaveSettings(settings settingstype.Settings) error {
	// Implementation to marshal Settings struct to JSON and write to a file
	userConfigDir := userSettingsDirPath()
	if userConfigDir == "" {
		return nil
	}

	// Ensure the user config directory exists
	err := os.MkdirAll(userConfigDir, os.ModePerm)
	if err != nil {
		return err
	}

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	settingsPath := userSettingsPath()
	return os.WriteFile(settingsPath, data, 0644)
}

func cleanTabSize(tabSize int) int {
	switch tabSize {
	case 2, 4, 8, 16:
		return tabSize
	default:
		return 4
	}
}

func cleanTabCharacter(tabCharacter string) string {
	if tabCharacter == "tab" || tabCharacter == "space" {
		return tabCharacter
	}
	return "space"
}

func cleanColorSchemeName(colorScheme string) string {
	if colorScheme == "" {
		return "default"
	}
	colorFiles := runtime.Files.ListRuntimeFiles(femto.RTColorscheme)
	for _, file := range colorFiles {
		if file.Name() == colorScheme {
			return colorScheme
		}
	}
	return "default"
}
