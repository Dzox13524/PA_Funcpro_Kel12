package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func (w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		data := map[string]string{"data":"ping"}
		json.NewEncoder(w).Encode(data)
	}) 
	http.HandleFunc("/pong", func (w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		data := map[string]string{"data":"pong"}
		json.NewEncoder(w).Encode(data)
	}) 
	http.ListenAndServe(":8080", nil);
}