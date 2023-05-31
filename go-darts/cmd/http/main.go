// Command go-darts
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var (
	scores        []int
	currentTurn   int
	numPlayers    = 2
	startingScore = 301
)

type playerScores struct {
	Scores []int
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
}

func main() {
	initializeGame()
	log.Info("Logging some darts info")
	router := mux.NewRouter()
	router.HandleFunc("/", gameStatus).Methods("GET")
	router.HandleFunc("/restart", restartGame).Methods("GET")
	router.HandleFunc("/{score}", scoreTurn).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://192.168.50.244:4200"},
		AllowCredentials: true,
	})

	server := http.Server{
		Handler:     c.Handler(router),
		Addr:        ":9090",
		ReadTimeout: time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func scoreTurn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	log.Info(r)
	log.Info(vars)

	intScore, err := strconv.Atoi(vars["score"])
	if err != nil {
		log.Fatal(err)
	}

	scores[currentTurn] -= intScore

	jsonOut, err := json.Marshal(playerScores{
		Scores: scores,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
	if currentTurn+1 == numPlayers {
		currentTurn = 0
	} else {
		currentTurn++
	}
}

func initializeGame() {
	scores = make([]int, numPlayers)
	currentTurn = 0

	for i := range scores {
		scores[i] = startingScore
	}
}

func restartGame(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	initializeGame()
	jsonOut, err := json.Marshal(playerScores{
		Scores: scores,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
}

func gameStatus(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	jsonOut, err := json.Marshal(playerScores{
		Scores: scores,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
}
