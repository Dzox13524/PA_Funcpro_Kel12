package response

import (
	"encoding/json"
	"net/http"
)
type APIresponse struct {
	Status string      `json:"status"` 
	Author string      `json:"author"` 
	Data   interface{} `json:"data"`   
}

func WriteJSON(w http.ResponseWriter, statusCode int, status string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIresponse {
		Status: status,
		Author: "Kelompok 12",
		Data: data,
	}

	 json.NewEncoder(w).Encode(response)
}