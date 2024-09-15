package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal json: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string){
	if code > 499 {
		log.Printf("Responding with 5XX error: %v", msg)
	}
	type errRes struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errRes{
		Error: msg,
	})
}