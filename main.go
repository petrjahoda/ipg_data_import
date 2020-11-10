package main

import (
	"github.com/kardianos/service"
	"time"
)

const version = "2020.4.2.10"
const serviceName = "IPG Data Import Service"
const serviceDescription = "Download users and products from CSV file and imports them into Zapsi database"
const downloadInSeconds = 600
const deleteLogsAfter = 240 * time.Hour

var serviceRunning = false
var processRunning = false
var zapsiConfig = "zapsi_uzivatel:zapsi@tcp(localhost:3306)/zapsi2?charset=utf8&parseTime=True&loc=Local"

type program struct{}

func main() {
	logInfo("MAIN", serviceName+" ["+version+"] starting...")
	logInfo("MAIN", serviceDescription)
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	program := &program{}
	s, err := service.New(program, serviceConfig)
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
	err = s.Run()
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
}

func (p *program) Start(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] started")
	serviceRunning = true
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	serviceRunning = false
	if processRunning {
		logInfo("MAIN", serviceName+" ["+version+"] stopping...")
		time.Sleep(1 * time.Second)
	}
	logInfo("MAIN", serviceName+" ["+version+"] stopped")
	return nil
}

func (p *program) run() {
	logDirectoryFileCheck("MAIN")
	createConfigIfNotExists()
	loadSettingsFromConfigFile()
	for serviceRunning {
		processRunning = true
		start := time.Now()
		logInfo("MAIN", serviceName+" ["+version+"] running")
		importData()
		sleepTime := downloadInSeconds*time.Second - time.Since(start)
		logInfo("MAIN", "Sleeping for "+sleepTime.String())
		time.Sleep(sleepTime)
		deleteOldLogFiles()
		processRunning = false
	}
}
