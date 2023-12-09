package config

import (
	"encoding/json"
	"os"
)

type LogSetting struct {
	FolderPath string
	FileName   string
}

type WebServerSetting struct {
	Title    string
	Host     string
	BasePath string
}

type AppConfig struct {
	Log         LogSetting
	Web         WebServerSetting
	SupabaseUrl string
	SupabaseKey string
}

func NewAppConfig(filePath string) (result *AppConfig, err error) {
	var (
		appConfig     *AppConfig
		appConfigFile *os.File
	)

	if appConfigFile, err = os.Open(filePath); err == nil {
		jsonParser := json.NewDecoder(appConfigFile)
		jsonParser.Decode(&appConfig)
	}

	defer appConfigFile.Close()

	return appConfig, err
}
