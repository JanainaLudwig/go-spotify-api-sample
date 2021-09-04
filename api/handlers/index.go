package handlers

import (
	"github.com/julienschmidt/httprouter"
	"go-api-template/config"
	"go-api-template/services/openweather"
	"go-api-template/services/spotify"
	"net/http"
)

type weatherResponse struct {
	Mood string  `json:"mood"`
	Temp float64 `json:"temp"`
}
type trackResponse struct {
	Name string `json:"name"`
}
type weatherTrackResponse struct {
	Weather weatherResponse `json:"weather"`
	Tracks []trackResponse `json:"tracks"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	city := r.URL.Query().Get("city")
	lat := r.URL.Query().Get("lat")
	long := r.URL.Query().Get("long")

	input := openweather.GetWeatherInput{
		City: city,
	}

	if lat != "" && long != "" {
		input.Coordinates = &openweather.CoordinatesInput{
			Lat:  lat,
			Long: long,
		}
	}

	weatherClient := openweather.NewClient(config.App.OpenWeatherApiKey)
	weather, err := weatherClient.GetWeather(r.Context(), input)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	weatherTackResponse := weatherTrackResponse{
		Weather: weatherResponse{
			Mood: weatherMood(weather.Main.Temp),
			Temp: weather.Main.Temp,
		},
	}

	spotifyClient, err := spotify.NewClient(r.Context(), config.App.Spotify.ClientId, config.App.Spotify.ClientSecret)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	recommendations, err := spotifyClient.GetRecommendations(r.Context(), weatherTackResponse.Weather.Mood)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	for _, track := range recommendations.Tracks {
		weatherTackResponse.Tracks = append(weatherTackResponse.Tracks, trackResponse{
			Name: track.Name,
		})
	}

	sendResponse(w, weatherTackResponse, http.StatusOK)
}

func weatherMood(temp float64) string {
	if temp > 30 {
		return "party"
	}

	if temp >= 15 {
		return "pop"
	}

	if temp > 10 {
		return "rock"
	}

	return "classical"
}
