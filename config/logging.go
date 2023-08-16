package config

import (
	"io"
	"log"
	"os"
)

func SetLogger(logFileName string) {
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to set up logger: %s", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiWriter)
}
