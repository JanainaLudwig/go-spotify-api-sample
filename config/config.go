package config

import (
	"path"
	"runtime"
)

var App AppConfig

type AppConfig struct {
	AppEnv string
	ApiPort string
	Spotify struct{
		ClientId string
		ClientSecret string
	}
	OpenWeatherApiKey string
}

func LoadEnv(path string) {
	load(path)

	App.AppEnv = loadString("APP_ENV", envStr("development"))
	App.ApiPort = loadString("API_PORT", envStr("8080"))

	App.Spotify.ClientId = loadString("SPOTIFY_CLIENT_ID", nil)
	App.Spotify.ClientSecret = loadString("SPOTIFY_CLIENT_SECRET", nil)

	App.OpenWeatherApiKey = loadString("OPEN_WEATHER_API_KEY", nil)
}

func RootPath() string {
	_, file, _, _ := runtime.Caller(0)

	root := path.Dir(path.Dir(file))

	return root
}