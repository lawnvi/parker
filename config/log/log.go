package log

import (
	"fmt"
	"log"
	"os"
)

var saveLevel = 1

const (
	DebugLv = 1
	InfoLv = 2
	WarningLv = 3
	ErrorLv = 4
	NoWrite = 10
)

func Init(level int, path string) {
	saveLevel = level
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile|log.Ldate|log.Ltime)
}

//打印
func printLog(level, where, detail string, save bool) {
	if detail == "" {
		return
	}
	log.SetPrefix("[" + level + "] ")
	detail = "<" + where + "> " + detail
	if save {
		log.Println(detail)
	}
	//fmt.Println(detail)
}

func D(tag, msg string) {
	printLog("d", tag, msg, DebugLv >= saveLevel)
}

func I(tag, msg string) {
	printLog("i", tag, msg, InfoLv >= saveLevel)
}

func W(tag, msg string) {
	printLog("w", tag, msg, WarningLv >= saveLevel)
}

func E(tag, msg string) {
	printLog("e", tag, msg, ErrorLv >= saveLevel)
}