package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	. "./config"
	. "./dao"
	. "./models"
	"strconv"
)

const (
	channelcodeCollection = "channelcodes"
)

var config = Config{}
var registryDB = MongoDB{}

// Registry statusInfo
type statusInfo struct {
    Name			string
    Description		string
    UtcCurrentTime	time.Time
    CurrentTime		time.Time
    Port			int
}
// GET a registry service status ( Port : 2201 )
func getStatus(w http.ResponseWriter, r *http.Request) {
	var status statusInfo;
	status.Name = "registry";
	status.Description = "Nanda Registry Service";
	status.UtcCurrentTime = time.Now().UTC();
	status.CurrentTime = time.Now();
	status.Port = 2201;
	respondWithJson(w, http.StatusOK, status)
}

// GET a channels
func getAllChannels(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query();
	var skip, limit int = 0, 1000
	var sort  = queryParams.Get("sort")
	if sort == ""{
		sort = "_id"
	}
	if queryParams.Get("sort")!= "" && queryParams.Get("sortType") == "DESC"{
		sort = "-"+sort
	}
	limit,_ = strconv.Atoi(queryParams.Get("limit"))
	skip,_ = strconv.Atoi(queryParams.Get("skip"))
	var query = bson.M{}
	allChannels, err := registryDB.FindAll(channelcodeCollection,query,skip,limit,sort)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, allChannels)
}

// POST a channel
func createChannel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var channelCodes ChannelcodeModel
	if err := json.NewDecoder(r.Body).Decode(&channelCodes); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	channelCodes.ID = bson.NewObjectId()
	if err := registryDB.Insert(channelcodeCollection,channelCodes); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, channelCodes)
}

// RespondWithError to handle error
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

// RespondWithJson to handle JSON data
func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	registryDB.Server = config.Server
	registryDB.Database = config.Database
	registryDB.Connect()
}

// Define HTTP request routes
func main() {
	registryRouter := mux.NewRouter()
	// Start : API
	registryRouter.HandleFunc("/registry/status", getStatus).Methods("GET")
	
	registryRouter.HandleFunc("/registry/getallchannels", getAllChannels).Methods("GET")
	registryRouter.HandleFunc("/registry/createchannel", createChannel).Methods("POST")
	// End : API
	if err := http.ListenAndServe(":2201", registryRouter); err != nil {
		log.Fatal(err)
	}
}
