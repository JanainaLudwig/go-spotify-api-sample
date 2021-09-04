package openweather

import (
	"context"
	"fmt"
	"go-api-template/core"
	"net/http"
)

type WeatherResponse struct {
	Main struct{
		Temp float64 `json:"temp"`
	} `json:"main"`
}

type GetWeatherInput struct {
	City string
	Coordinates *CoordinatesInput
}
type CoordinatesInput struct {
	Lat string
	Long string
}

func (g *GetWeatherInput) Validate() error {
	if g.City == "" && g.Coordinates == nil {
		return core.NewAppError(core.AppError{
			Message: "Invalid params. Please provide city or coordinates",
			Code: "invalid-params",
			Status: http.StatusBadRequest,
		})
	}

	if g.Coordinates != nil && (g.Coordinates.Lat == "" || g.Coordinates.Long == "") {
		return core.NewAppError(core.AppError{
			Message: "Invalid params. Please provide city or coordinates",
			Code: "invalid-coordinates",
			Status: http.StatusBadRequest,
		})
	}

	return nil
}

func (c * Client) GetWeather(ctx context.Context, input GetWeatherInput) (*WeatherResponse, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	uri := "/weather?units=metric"

	if input.City != "" {
		uri += "&q=" + input.City
	} else {
		uri = fmt.Sprintf("%s&lat=%s&lon=%s", uri, input.Coordinates.Lat, input.Coordinates.Long)
	}

	data := WeatherResponse{}
	err = c.GetRequest(ctx, uri, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
