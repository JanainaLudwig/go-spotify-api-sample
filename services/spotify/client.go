package spotify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"go-api-template/core"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Token string
	httpClient http.Client
}

const (
	authBaseUri   = "https://accounts.spotify.com/api"
	apiBaseUri    = "https://api.spotify.com"
	authorization = "Authorization"
)

func NewClient(ctx context.Context, clientId, clientSecret string) (*Client, error) {
	data := url.Values{}
	data.Add("grant_type", "client_credentials")

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, authBaseUri+"/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, core.NewAppError(core.AppError{
			Message: err.Error(),
		})
	}

	code := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
	request.Header.Add(authorization, "Basic "+code)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client := Client{
		httpClient: http.Client{
			Timeout: 10 * time.Second,
		},
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, core.NewAppError(core.AppError{
			Message: err.Error(),
		})
	}

	res := struct {
		AccessToken string `json:"access_token"`
	}{}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return nil, core.NewAppError(core.AppError{
			Message: err.Error(),
		})
	}

	client.Token = res.AccessToken

	return &client, nil
}

type RequestError struct {
	Error struct{
		Status int `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Client) GetRequest(ctx context.Context, uri string, data interface{}) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, apiBaseUri+uri, nil)
	if err != nil {
		return core.NewAppError(core.AppError{
			Message: err.Error(),
		})
	}

	request.Header.Add(authorization, "Bearer " + c.Token)

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
		responseError := RequestError{}
		err = json.Unmarshal(resBody, &responseError)
		if err != nil {
			return core.NewAppError(core.AppError{
				Message: err.Error(),
			})
		}

		return core.NewAppError(core.AppError{
			Message: "spotify api error: " + responseError.Error.Message,
			Status: responseError.Error.Status,
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
