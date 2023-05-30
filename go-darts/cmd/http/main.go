// Command go-darts
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
}

func main() {
	log.Info("Logging some darts info")
	router := mux.NewRouter()
	router.HandleFunc("/{score}", homeHandler).Methods("GET")

	server := http.Server{
		Handler:     router,
		Addr:        ":8080",
		ReadTimeout: time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	log.Info(r)
	log.Info(vars)

	data := map[string]int{
		"playerOne": 300,
		"playerTwo": 301,
	}

	jsonOut, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
}
