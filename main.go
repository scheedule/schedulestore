package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/scheedule/schedulestore/commands"
)

func init() {
	// Default error level
	log.SetLevel(log.ErrorLevel)
}

func main() {
	commands.Execute()
}
