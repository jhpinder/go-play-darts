// Command go-darts
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	startingScore = 301
	gameColl      *mongo.Collection
	playerColl    *mongo.Collection
	ctx           = context.Background()
)

func dartBoardScoreMapping(shot int) int {
	internalMap := [83]int{
		0,
		20, 60, 20, 40,
		1, 3, 1, 2,
		18, 54, 18, 36,
		4, 12, 4, 8,
		13, 39, 13, 26,
		6, 18, 6, 12,
		10, 30, 10, 20,
		15, 45, 15, 30,
		2, 6, 2, 4,
		17, 51, 17, 34,
		3, 9, 3, 6,
		19, 57, 19, 38,
		7, 21, 7, 14,
		16, 48, 16, 32,
		8, 24, 8, 16,
		11, 33, 11, 22,
		14, 42, 14, 28,
		9, 27, 9, 18,
		12, 36, 12, 24,
		5, 15, 5, 10,
		50, 25,
	}
	return internalMap[shot]
}

type countdownGame struct {
	Scores             map[string]int   `bson:"scores"`
	GameID             string           `bson:"gameID"`
	Shots              map[string][]int `bson:"shots"`
	OrderedPlayers     []string         `bson:"orderedPlayers"`
	CurrentPlayerIndex int              `bson:"currentPlayerIndex"`
}

type player struct {
	PlayerID   string `bson:"playerID"`
	PlayerName string `bson:"playerName"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
}

func init() {
	initMongo()
}

func main() {
	log.Info("Starting HTTP server")

	runHTTPServer()
}

func initMongo() {
	mongoURI := "mongodb://mongoadmin:secret@localhost:27017/admin"

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	gameColl = mongoClient.Database("darts").Collection("game")
	playerColl = mongoClient.Database("darts").Collection("player")
}

func runHTTPServer() {
	router := mux.NewRouter()
	router.HandleFunc("/game/new", newGame).Methods("POST")
	router.HandleFunc("/game/{gameID}", gameStatus).Methods("GET")
	router.HandleFunc("/game/{gameID}", scoreThrow).Methods("POST")

	router.HandleFunc("/player/new", newPlayer).Methods("POST")
	router.HandleFunc("/player/{playerID}", getPlayerName).Methods("GET")
	router.HandleFunc("/player/{playerID}", updatePlayerName).Methods("POST")

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

func getGameStatusFromDB(gameID string) (game countdownGame) {
	err := gameColl.FindOne(ctx, bson.D{{Key: "gameID", Value: gameID}}).Decode(&game)
	if err != nil {
		log.Debug("Could not find game in DB")
	}
	return
}

func updateGameStatusToDB(currentGame countdownGame) {
	filter := bson.D{{Key: "gameID", Value: currentGame.GameID}}
	update := bson.D{{Key: "$set", Value: currentGame}}
	_, err := gameColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func scoreThrow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	log.Info(r)
	log.Info(vars)

	gameID := vars["gameID"]

	shot := int(0)
	err := json.NewDecoder(r.Body).Decode(&shot)
	if err != nil {
		log.Fatal()
	}

	currentGame := getGameStatusFromDB(gameID)
	player := currentGame.OrderedPlayers[currentGame.CurrentPlayerIndex]
	shots := currentGame.Shots[player]

	// add a shot
	shots = append(shots, shot)
	currentGame.Shots[player] = shots

	// update score
	currentGame.Scores[player] -= dartBoardScoreMapping(shot)

	// if number of shots by current player mod 3 is 0, next player
	if len(shots)%3 == 0 {
		// increment the index of the current player
		currentGame.CurrentPlayerIndex++
		if currentGame.CurrentPlayerIndex >= len(currentGame.OrderedPlayers) {
			currentGame.CurrentPlayerIndex = 0
		}
	}

	updateGameStatusToDB(currentGame)

	jsonOut, err := json.Marshal(currentGame)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
}

func newGame(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	requestBody := []string{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Fatal()
	}

	doc := countdownGame{
		Scores:             map[string]int{},
		GameID:             uuid.NewString(),
		Shots:              map[string][]int{},
		OrderedPlayers:     []string{},
		CurrentPlayerIndex: 0,
	}

	for i := range requestBody {
		doc.Scores[requestBody[i]] = startingScore
		doc.Shots[requestBody[i]] = []int{}
	}

	doc.OrderedPlayers = requestBody
	log.Debug(requestBody)

	_, err = gameColl.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}

	jsonOut, err := json.Marshal(doc)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
}

func gameStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debugf("%+v", vars)

	gameID := vars["gameID"]

	currentGame := getGameStatusFromDB(gameID)

	w.WriteHeader(http.StatusOK)
	jsonOut, err := json.Marshal(currentGame)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
}

func newPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug(r)
	log.Debugf("%+v", vars)

	incomingRequest := player{}
	err := json.NewDecoder(r.Body).Decode(&incomingRequest)
	if err != nil {
		log.Fatal()
	}

	playerToAdd := player{
		PlayerID:   uuid.NewString(),
		PlayerName: incomingRequest.PlayerName,
	}

	_, err = playerColl.InsertOne(ctx, playerToAdd)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)

	jsonOut, err := json.Marshal(playerToAdd)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
}

func getPlayerName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug(r)
	log.Debugf("%+v", vars)

	playerID, ok := vars["playerID"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := player{}

	err := playerColl.FindOne(ctx, bson.D{{Key: "playerID", Value: playerID}}).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)

	jsonOut, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(w, string(jsonOut))
	if err != nil {
		log.Fatal(err)
	}
}

func updatePlayerName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug(r)
	log.Debugf("%+v", vars)

	playerID, ok := vars["playerID"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	incomingRequest := player{}
	err := json.NewDecoder(r.Body).Decode(&incomingRequest)
	if err != nil {
		log.Fatal()
	}

	filter := bson.D{{Key: "playerID", Value: playerID}}
	update := bson.D{{
		Key: "$set",
		Value: bson.D{{
			Key:   "playerName",
			Value: incomingRequest.PlayerName,
		}},
	}}

	response, err := playerColl.UpdateOne(ctx, filter, update)
	if err != nil || response.MatchedCount != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = fmt.Fprintln(w, "Unable to update player name")
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = fmt.Fprintln(w, "Updated player name")
	if err != nil {
		log.Fatal(err)
	}
}
