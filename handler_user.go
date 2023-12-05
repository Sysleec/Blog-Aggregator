package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sysleec/Blog-Aggregator/internal/auth"
	"github.com/Sysleec/Blog-Aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Cant create user: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (apiCfg apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth error: %s", err))
		return
	}

	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't get user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
