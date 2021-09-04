package spotify

import "context"

type RecommendationsResponse struct {
	Tracks []struct{
		Href string `json:"href"`
		Name string `json:"name"`
	} `json:"tracks"`
}

func (c * Client) GetRecommendations(ctx context.Context, category string) (*RecommendationsResponse, error) {
	uri := "/v1/recommendations?seed_genres=" + category

	res := RecommendationsResponse{}
	err := c.GetRequest(ctx, uri, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
