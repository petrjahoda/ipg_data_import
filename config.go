package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	DatabaseType string
	IpAddress    string
	DatabaseName string
	Port         string
	Login        string
	Password     string
}

func createConfigIfNotExists() {
	configDirectory := filepath.Join(".", "config")
	configFileName := "config.json"
	configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
	if _, checkPathError := os.Stat(configFullPath); checkPathError == nil {
		logInfo("MAIN", "Config file already exists")
	} else if os.IsNotExist(checkPathError) {
		logInfo("MAIN", "Config file does not exist, creating")
		mkdirError := os.MkdirAll(configDirectory, 0777)
		if mkdirError != nil {
			logError("MAIN", "Unable to create directory for config file: "+mkdirError.Error())
		} else {
			logInfo("MAIN", "Directory for config file created")
			data := Config{
				DatabaseType: "mysql",
				IpAddress:    "localhost",
				DatabaseName: "zapsi2",
				Port:         "3306",
				Login:        "zapsi_uzivatel",
				Password:     "zapsi",
			}
			file, _ := json.MarshalIndent(data, "", "  ")
			writingError := ioutil.WriteFile(configFullPath, file, 0666)
			logInfo("MAIN", "Writing data to JSON file")
			if writingError != nil {
				logError("MAIN", "Unable to write data to config file: "+writingError.Error())
			} else {
				logInfo("MAIN", "Data written to config file")
			}
		}
	} else {
		logError("MAIN", "Config file does not exist")
	}
}

func loadSettingsFromConfigFile() {
	configDirectory := filepath.Join(".", "config")
	configFileName := "config.json"
	configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
	ConfigFile := Config{}
	for len(ConfigFile.DatabaseName) == 0 {
		readFile, err := ioutil.ReadFile(configFullPath)
		if err != nil {
			logError("MAIN", "Problem reading config file")
			var err = os.Remove(configFullPath)
			if err != nil {
				logError("MAIN", "Problem deleting file "+configFullPath+", "+err.Error())
				break
			}
			createConfigIfNotExists()
		}
		err = json.Unmarshal(readFile, &ConfigFile)
		if err != nil {
			logError("MAIN", "Problem parsing config file, deleting config file")
			var err = os.Remove(configFullPath)
			if err != nil {
				logError("MAIN", "Problem deleting file "+configFullPath+", "+err.Error())
				break
			}
			createConfigIfNotExists()
		}
	}
	zapsiConfig = ConfigFile.Login + ":" + ConfigFile.Password + "@tcp(" + ConfigFile.IpAddress + ":" + ConfigFile.Port + ")/" + ConfigFile.DatabaseName + "?charset=utf8&parseTime=True&loc=Local"
}
