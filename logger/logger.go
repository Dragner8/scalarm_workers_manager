package logger

import "log"

var verbose bool = false

func SetVerbosity(verboseMode bool) {
	verbose = verboseMode
}

func Info(format string, a ...interface{}) {
	log.Printf("[I] "+format, a...)
}

func Debug(format string, a ...interface{}) {
	if verbose {
		log.Printf("[D] "+format, a...)
	}
}

func Fatal(msg string) {
	log.Fatal("[F] " + msg)
}
