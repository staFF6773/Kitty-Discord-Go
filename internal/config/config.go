package config

type Configuration struct {
	Token string `json:"token"`

	Logging struct {
		Method []string `json:"method"`
		File   string   `json:"file"`
		Level  string   `json:"level"`
	} `json:"logging"`
}
