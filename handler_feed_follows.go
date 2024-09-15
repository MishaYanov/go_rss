package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MishaYanov/rssagg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) HandleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct {
		UserID uuid.UUID `json:"user_id"`
		FeedID uuid.UUID `json:"fedd_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	feedFollow, err := apiCfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: params.UserID,
		FeedID: params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating Feed: %v", err))
		return
	}

	respondWithJSON(w, 200, dbFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) HandleGetUserFeedFollows(w http.ResponseWriter, r *http.Request, user database.User){
	feedFollows, err := apiCfg.DB.GetUserFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error fetching Feeds: %v", err))
		return
	}
	respondWithJSON(w, 200, dbFeedFollowsToFeedFollows(feedFollows))
}