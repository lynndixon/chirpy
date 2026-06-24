package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIDString := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder != "desc" {
		sortOrder = "asc"
	}
	if authorIDString != "" {
		authorID, err := uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}

		chirps, err := cfg.db.GetChirpsByAuthor(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps by author", err)
			return
		}

		respChirps := []Chirp{}
		for _, chirp := range chirps {
			respChirp := Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			}
			respChirps = append(respChirps, respChirp)
		}
		sort.Slice(respChirps, func(i, j int) bool {
			if sortOrder == "desc" {
				return respChirps[i].CreatedAt.After(respChirps[j].CreatedAt)
			}
			return respChirps[i].CreatedAt.Before(respChirps[j].CreatedAt)
		})

		respondWithJSON(w, http.StatusOK, respChirps)
		return
	}

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	respChirps := []Chirp{}

	for _, chirp := range chirps {
		respChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		respChirps = append(respChirps, respChirp)
	}
	sort.Slice(respChirps, func(i, j int) bool {
		if sortOrder == "desc" {
			return respChirps[i].CreatedAt.After(respChirps[j].CreatedAt)
		}
		return respChirps[i].CreatedAt.Before(respChirps[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, respChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}
	respChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, respChirp)
}
