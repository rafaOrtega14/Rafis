package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"sync"
)

type Data struct {
	Value interface{} `json:"value"`
	Key string `json:"key"`
}

var dataMap = make(map[string]interface{})
var result = make(chan interface{})
var lock = sync.RWMutex{}

func getData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"];
	go readData(key)
	data := <- result
	json.NewEncoder(w).Encode(data)
}
func insertData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data Data
	_ = json.NewDecoder(r.Body).Decode(&data)
	go writeData(data)
	json.NewEncoder(w).Encode(&data)
}
func readData(key string) {
	result <- dataMap[key];
}
func writeData(data Data) {
	lock.Lock()
	defer lock.Unlock();
	dataMap[data.Key] = data.Value
}
func main() {
    r := mux.NewRouter()
    r.HandleFunc("/getData/{key}", getData).Methods("GET");
	r.HandleFunc("/insert", insertData).Methods("POST");

	fmt.Printf("Server listening to port 80")
    http.ListenAndServe(":80", r)
}