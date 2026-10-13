package pokeapi

import (
	"encoding/json"
	"fmt"
	"github.com/BashCtl/pokedex-cli/internal/pokecache"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

func NewClient(timeout time.Duration) Client {
	return Client{httpClient: http.Client{Timeout: timeout},
		cache: pokecache.NewCache(timeout)}
}

type LocationAreasResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) ListLocationAreas(pageURL *string) (LocationAreasResp, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	res, err := c.httpClient.Get(url)
	if err != nil {
		return LocationAreasResp{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return LocationAreasResp{}, fmt.Errorf("bad status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreasResp{}, err
	}

	var locations LocationAreasResp
	if err := json.Unmarshal(body, &locations); err != nil {
		return LocationAreasResp{}, err
	}
	return locations, nil
}
