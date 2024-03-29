package main

import (
	"fmt"
	"github.com/TwinProduction/go-color"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func logInfo(reference, data string) {
	fmt.Println(color.Ize(color.Green, "["+reference+"] --INF-- "+data))
	appendDataToLog("INF", reference, data)
}

func logError(reference, data string) {
	fmt.Println(color.Ize(color.Red, "["+reference+"] --ERR-- "+data))
	appendDataToLog("ERR", reference, data)
	appendDataToErrLog("ERR", reference, data)
}

func logDirectoryFileCheck(reference string) {
	dateTimeFormat := "2006-01-02 15:04:05.000"
	dir := GetDirectory()
	logDirectory := filepath.Join(dir, "log")
	_, checkPathError := os.Stat(logDirectory)
	logDirectoryExists := checkPathError == nil
	if logDirectoryExists {
		fmt.Println(color.Ize(color.Green,
			"["+reference+"] --INF-- "+"Log directory already exists "))
		return
	}
	fmt.Println(color.Ize(color.Yellow, time.Now().Format(dateTimeFormat)+" ["+reference+"] --WRN-- "+"Log directory does not exist, creating"))
	mkdirError := os.MkdirAll(logDirectory, 0777)
	if mkdirError != nil {
		fmt.Println(color.Ize(color.Red, time.Now().Format(dateTimeFormat)+" ["+reference+"] --ERR--"+"Unable to create directory for log file: "+mkdirError.Error()))
		return
	}
}

func appendDataToLog(logLevel string, reference string, data string) {
	dateTimeFormat := "2006-01-02 15:04:05.000"
	logNameDateTimeFormat := "2006-01-02"
	dir := GetDirectory()
	logDirectory := filepath.Join(dir, "log")
	logFileName := reference + " " + time.Now().Format(logNameDateTimeFormat) + ".log"
	logFullPath := strings.Join([]string{logDirectory, logFileName}, "/")
	f, err := os.OpenFile(logFullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(color.Ize(color.Yellow, time.Now().Format(dateTimeFormat)+" ["+reference+"] --WAR-- "+"Log file not present: "+err.Error()))
		return
	}
	defer f.Close()
	logData := time.Now().Format("2006-01-02 15:04:05.000   ") + reference + "   " + logLevel + "   " + data
	if _, err := f.WriteString(logData + "\r\n"); err != nil {
		fmt.Println(color.Ize(color.Red, time.Now().Format(dateTimeFormat)+" ["+reference+"] --ERR-- "+"Cannot write to file: "+err.Error()))
	}
}

func appendDataToErrLog(logLevel string, reference string, data string) {
	dateTimeFormat := "2006-01-02 15:04:05.000"
	logNameDateTimeFormat := "2006-01-02"
	dir := GetDirectory()
	logDirectory := filepath.Join(dir, "log")
	logFileName := reference + " " + time.Now().Format(logNameDateTimeFormat) + ".err"
	logFullPath := strings.Join([]string{logDirectory, logFileName}, "/")
	f, err := os.OpenFile(logFullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(color.Ize(color.Red, time.Now().Format(dateTimeFormat)+" ["+reference+"] --WAR-- "+"Log file not present: "+err.Error()))
		return
	}
	defer f.Close()
	logData := time.Now().Format("2006-01-02 15:04:05.000   ") + reference + "   " + logLevel + "   " + data
	if _, err := f.WriteString(logData + "\r\n"); err != nil {
		fmt.Println(color.Ize(color.Red, time.Now().Format(dateTimeFormat)+" ["+reference+"] --ERR-- "+"Cannot write to file: "+err.Error()))
	}
}

func deleteOldLogFiles() {
	logInfo("MAIN", "Deleting old log files")
	timer := time.Now()
	directory, err := ioutil.ReadDir("log")
	if err != nil {
		logError("MAIN", "Problem opening log directory")
		return
	}
	now := time.Now()
	logDirectory := filepath.Join(".", "log")
	for _, file := range directory {
		if fileAge := now.Sub(file.ModTime()); fileAge > deleteLogsAfter {
			logInfo("MAIN", "Deleting old log file "+file.Name()+" with age of "+fileAge.String())
			logFullPath := strings.Join([]string{logDirectory, file.Name()}, "/")
			var err = os.Remove(logFullPath)
			if err != nil {
				logError("MAIN", "Problem deleting file "+file.Name()+", "+err.Error())
			}
		}
	}
	logInfo("MAIN", "Old log files deleted, elapsed: "+time.Since(timer).String())
}

func GetDirectory() string {
	var dir string
	if runtime.GOOS == "windows" {
		executable, err := os.Executable()
		if err != nil {
			panic(err)
		}
		dir = filepath.Dir(executable)
	} else {
		dir, _ = os.Getwd()
	}
	return dir
}
