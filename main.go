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
var lock = sync.RWMutex{}

func getData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"];
	data := readData(key)
	json.NewEncoder(w).Encode(data)
}
func insertData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data Data
	_ = json.NewDecoder(r.Body).Decode(&data)
	writeData(data)
	json.NewEncoder(w).Encode(&data)
}
func readData(key string) interface{} {
	lock.Lock()
	defer lock.Unlock();
	return dataMap[key];
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

	fmt.Printf("Server listening to port 5000")
    http.ListenAndServe(":5000", r)
}
