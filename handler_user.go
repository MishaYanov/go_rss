package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MishaYanov/rssagg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	respondWithJSON(w, 201, dbUserToUser(user))
}

func (apiCfg *apiConfig) HandleGetUser(w http.ResponseWriter, r *http.Request, user database.User){
	respondWithJSON(w, 200, dbUserToUser(user))
}

func (apiCfg *apiConfig) HandleFollowedPosts(w http.ResponseWriter, r *http.Request, user database.User){
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: 10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error fetching posts: %v", err))
		return
	}

	respondWithJSON(w, 200, dbPostsToPosts(posts))
}
