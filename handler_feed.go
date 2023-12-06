package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sysleec/Blog-Aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Cant create feed: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedToFeed(feed))
}

func (apiCfg apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't get feeds: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedsToFeeds(feeds))
}
