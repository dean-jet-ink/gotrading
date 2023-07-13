package utils

import (
	"io"
	"log"
	"os"
)

func SetLogging(logFileName string) {
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("file=logFile error=%s", err.Error())
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiWriter)
}
