package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	models "go-url-shortener/models"
	utils "go-url-shortener/utils"
)

var urlMap sync.Map

// Handler to create a short URL
func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var shortCode string
	if req.CustomUrl == "" {
		shortCode = utils.GenerateShortURL()
	} else {
		if _, ok := urlMap.Load(req.CustomUrl); ok {
			url := fmt.Sprintf("http://%s/%s", r.Host, req.CustomUrl)
			http.Error(w, "Shortened URL " + url + " already exists!", http.StatusConflict)
			return
		}

		if len(req.CustomUrl) != 6 {
			http.Error(w, "Shortened URL must be 6 characters", http.StatusBadRequest)
			return
		}

		shortCode = req.CustomUrl
	}

	urlMap.Store(shortCode, req.URL)

	resp := models.ShortenResponse{
		ShortURL: fmt.Sprintf("http://%s/%s", r.Host, shortCode),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp) // Assume error handling here as per previous discussion
}

// Handler to redirect to the original URL
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:]

	if url, ok := urlMap.Load(shortCode); ok {
		http.Redirect(w, r, url.(string), http.StatusFound)
		return
	}

	http.NotFound(w, r)
}