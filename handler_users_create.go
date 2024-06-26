package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type returnUser struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	//create user
	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	respondWithJSON(w, http.StatusCreated, returnUser{
		Email: user.Email,
		ID:    user.ID,
	})
}
