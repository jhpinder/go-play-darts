// Command go-darts
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	currentGame   countdownGame
	numPlayers    = 2
	startingScore = 301
	coll          *mongo.Collection
	ctx           = context.Background()
)

type playerScores struct {
	Scores []int
}

type countdownGame struct {
	Scores      []int `bson:"scores"`
	CurrentTurn int   `bson:"currentTurn"`
	GameID      int   `bson:"gameID"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
}

func init() {
	initializeGame()
	initMongo()
}

func main() {
	log.Info("Starting HTTP server")

	getGameStatusFromDB(currentGame.GameID)

	runHTTPServer()
}

func initMongo() {
	mongoURI := "mongodb://mongoadmin:secret@localhost:27017/admin"

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	coll = mongoClient.Database("mongo-darts").Collection("games")

	getGameStatusFromDB(currentGame.GameID)
}

func runHTTPServer() {
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

func getGameStatusFromDB(gameID int) {
	var game countdownGame
	err := coll.FindOne(ctx, bson.D{{Key: "gameID", Value: gameID}}).Decode(&game)
	if err != nil {
		log.Fatal(err)
	}
	currentGame = game
}

func updateGameStatusToDB() {
	filter := bson.D{{Key: "gameID", Value: currentGame.GameID}}
	update := bson.D{{Key: "$set", Value: currentGame}}
	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func scoreTurn(w http.ResponseWriter, r *http.Request) {
	getGameStatusFromDB(currentGame.GameID)

	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	log.Info(r)
	log.Info(vars)

	intScore, err := strconv.Atoi(vars["score"])
	if err != nil {
		log.Fatal(err)
	}

	currentGame.Scores[currentGame.CurrentTurn] -= intScore

	jsonOut, err := json.Marshal(playerScores{
		Scores: currentGame.Scores,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
	if currentGame.CurrentTurn+1 == numPlayers {
		currentGame.CurrentTurn = 0
	} else {
		currentGame.CurrentTurn++
	}
	updateGameStatusToDB()
}

func initializeGame() {
	currentGame = countdownGame{
		Scores:      make([]int, numPlayers),
		CurrentTurn: 0,
		GameID:      1000,
	}
}

func restartGame(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	currentGame.CurrentTurn = 0
	for i := range currentGame.Scores {
		currentGame.Scores[i] = startingScore
	}

	jsonOut, err := json.Marshal(playerScores{
		Scores: currentGame.Scores,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
	updateGameStatusToDB()
}

func gameStatus(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	jsonOut, err := json.Marshal(playerScores{
		Scores: currentGame.Scores,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
}
