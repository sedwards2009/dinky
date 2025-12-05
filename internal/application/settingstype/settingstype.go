package settingstype

type Settings struct {
	ShowWhitespace         bool   `json:"showWhitespace"`
	ShowLineNumbers        bool   `json:"showLineNumbers"`
	ShowMatchBracket       bool   `json:"showMatchBracket"`
	SoftWrap               bool   `json:"softWrap"`
	TabSize                int    `json:"tabSize"`
	TabCharacter           string `json:"tabCharacter"` // "tab" or "space"
	ColorScheme            string `json:"colorScheme"`
	ShowTrailingWhitespace bool   `json:"showTrailingWhitespace"`
}

func DefaultSettings() Settings {
	return Settings{
		ShowWhitespace:         false,
		ShowLineNumbers:        true,
		ShowMatchBracket:       true,
		SoftWrap:               false,
		TabSize:                4,
		TabCharacter:           "space",
		ColorScheme:            "default",
		ShowTrailingWhitespace: true,
	}
}
