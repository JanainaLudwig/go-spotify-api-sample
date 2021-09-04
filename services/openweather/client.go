package openweather

import (
	"context"
	"encoding/json"
	"go-api-template/core"
	"io"
	"net/http"
	"strings"
)

const (
	baseUri = "https://api.openweathermap.org/data/2.5"
)

type Client struct {
	key string
	httpClient http.Client
}

func NewClient(key string) *Client {
	return &Client{key: key}
}

type ResponseError struct {
	Message string `json:"message"`
	Status int `json:"status"`
}

func (c *Client) GetRequest(ctx context.Context, uri string, data interface{}) error {
	url := baseUri+uri

	if strings.Contains(url, "?") {
		url = url + "&appid=" + c.key
	} else {
		url = url + "?appid=" + c.key
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return core.NewAppError(core.AppError{
			Message: err.Error(),
		})
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return core.NewAppError(core.AppError{
			Message: err.Error(),
		})
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return core.NewAppError(core.AppError{
			Message: err.Error(),
		})
	}

	if response.StatusCode >= 400 || response.StatusCode < 200 {
		responseError := ResponseError{}
		err = json.Unmarshal(resBody, &responseError)
		if err != nil {
			return core.NewAppError(core.AppError{
				Message: err.Error(),
			})
		}

		return core.NewAppError(core.AppError{
			Message: "open weather error: " + responseError.Message,
			Status: responseError.Status,
		})
	}

	err = json.Unmarshal(resBody, data)
	if err != nil {
		return core.NewAppError(core.AppError{
			Message: err.Error(),
		})
	}

	return nil
}
