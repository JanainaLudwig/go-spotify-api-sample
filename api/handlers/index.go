package handlers

import (
	"github.com/julienschmidt/httprouter"
	"go-api-template/config"
	"go-api-template/services/spotify"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	client, err := spotify.NewClient(r.Context(), config.App.Spotify.ClientId, config.App.Spotify.ClientSecret)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	recommendations, err := client.GetRecommendations(r.Context(), "party")
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	sendResponse(w, recommendations, http.StatusOK)
}
