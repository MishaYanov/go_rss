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
		FeedID uuid.UUID `json:"feed_id"`
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

func (apiCfg *apiConfig) HandleDeleteUserFeedFollows(w http.ResponseWriter, r *http.Request, user database.User){
	FeedFollowId := r.PathValue("feedFollowId")
	if FeedFollowId == "" {
		respondWithError(w, 400, "No feed follow selected to delete")
		return
	}

	id, err := uuid.Parse(FeedFollowId)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing UUID: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID: id,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 500, "failed to delete follow")
		return
	}

	respondWithJSON(w, 200, "Follow deleted successfully")
}
