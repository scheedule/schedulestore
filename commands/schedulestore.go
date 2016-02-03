// Package commands hold all the commands and subcommands for this application.
// CLI flags, environment vars, etc. will be extracted and settings will be set.
package commands

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/scheedule/schedulestore/api"
	"github.com/scheedule/schedulestore/db"
)

// Main command of the program
var schedulestoreCmd = &cobra.Command{
	Use:   "schedulestore",
	Short: "Schedule key value store.",
	Long:  "Serve storage and retrieval endpoint for schedules.",
	Run: func(cmd *cobra.Command, args []string) {
		initializeConfig()

		// Create DB Object
		myDB := db.New(dbHost, dbPort, database, collection)
		err := myDB.Init()
		if err != nil {
			log.Fatal("DB failure: ", err)
		}

		// API Object
		myAPI := api.New(myDB)

		r := mux.NewRouter()

		r.HandleFunc("/", myAPI.HandleGET).Methods("GET")
		r.HandleFunc("/", myAPI.HandlePUT).Methods("PUT")
		r.HandleFunc("/", myAPI.HandleDELETE).Methods("DELETE")

		if !public {
			http.ListenAndServe(":"+servePort, PrivateMiddleware(r))
		} else {
			http.ListenAndServe(":"+servePort, r)
		}

		log.Info("Serving on port:", servePort)

	},
}

func PublicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("user_id", r.RemoteAddr)
		h.ServeHTTP(w, r)
	})
}

func PrivateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID := r.Header.Get("user_id")
		if userID == "" {
			log.Warn("user_id not set, rejecting.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

var verbose, public bool
var servePort, dbHost, dbPort, database, collection string

// Initialize Flags
func init() {
	schedulestoreCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	schedulestoreCmd.Flags().BoolVarP(&public, "public", "p", false, "For debugging, use ip address as user_id")

	schedulestoreCmd.Flags().StringVarP(
		&dbHost, "db_host", "", "localhost", "Hostname of DB to insert into and retrieve from.")

	schedulestoreCmd.Flags().StringVarP(
		&dbPort, "db_port", "", "27017", "Port to access DB on.")

	schedulestoreCmd.Flags().StringVarP(
		&servePort, "serve_port", "", "5000", "Port to serve endpoint on.")

	schedulestoreCmd.Flags().StringVarP(
		&database, "db_name", "", "test", "Database name.")

	schedulestoreCmd.Flags().StringVarP(
		&collection, "db_collection", "", "schedules", "Collection in database for schedules.")
}

// Initialize configuration settings
func initializeConfig() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

// Execute schedulestore command
func Execute() {
	if err := schedulestoreCmd.Execute(); err != nil {
		panic(err)
	}
}
