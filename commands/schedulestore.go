// Package commands hold all the commands and subcommands for this application.
// CLI flags, environment vars, etc. will be extracted and settings will be set.
package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/scheedule/schedulestore/api"
	"github.com/scheedule/schedulestore/db"
	"github.com/spf13/cobra"
	"net/http"
)

// Main command of the program
var ScheduleStoreCmd = &cobra.Command{
	Use:   "schedulestore",
	Short: "Schedule key value store.",
	Long:  "Serve storage and retrieval endpoint for schedules.",
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()

		// Create DB Object
		mydb := db.NewDB(db_host, db_port, database, collection)
		mydb.Init()

		// API Object
		myapi := &api.Api{mydb}

		http.HandleFunc("/lookup", myapi.HandleLookup)
		http.HandleFunc("/put", myapi.HandlePut)
		log.Info("Serving on port:", serve_port)
		http.ListenAndServe(":"+serve_port, nil)
	},
}

var Verbose bool
var serve_port, db_host, db_port, database, collection string

// Initialize Flags
func init() {
	ScheduleStoreCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	ScheduleStoreCmd.Flags().StringVarP(
		&db_host, "db_host", "", "localhost", "Hostname of DB to insert into and retrieve from.")

	ScheduleStoreCmd.Flags().StringVarP(
		&db_port, "db_port", "", "27017", "Port to access DB on.")

	ScheduleStoreCmd.Flags().StringVarP(
		&serve_port, "serve_port", "", "5000", "Port to serve endpoint on.")

	ScheduleStoreCmd.Flags().StringVarP(
		&database, "db_name", "", "test", "Database name.")

	ScheduleStoreCmd.Flags().StringVarP(
		&collection, "db_collection", "", "schedules", "Collection in database for schedules.")
}

// Initialize configuration settings
func InitializeConfig() {
	if Verbose {
		log.SetLevel(log.DebugLevel)
	}
}

// Execute schedulestore command
func Execute() {
	if err := ScheduleStoreCmd.Execute(); err != nil {
		panic(err)
	}
}
