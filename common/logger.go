package common

import (
	"log"
	"os"
)

var logger *log.Logger

func Loginit() {
	f, err := os.OpenFile("book.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	// defer f.Close()
	logger = log.New(f, "", log.LstdFlags)
	logger.Println("Log init")
}

func Log(ns string, log interface{}) {
	logger.Println(ns, log)
}
