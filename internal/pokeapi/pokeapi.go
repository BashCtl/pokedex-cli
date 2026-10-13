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
func (c *Client) getBody(url string) ([]byte, error) {
	if body, ok := c.cache.Get(url); ok {
		return body, nil
	}

	res, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	c.cache.Add(url, body)
	return body, nil
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

type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
}

type LocationAreaDetails struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

func (c *Client) ListLocationAreas(pageURL *string) (LocationAreasResp, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	body, err := c.getBody(url)
	if err != nil {
		return LocationAreasResp{}, err
	}

	var locations LocationAreasResp
	if err := json.Unmarshal(body, &locations); err != nil {
		return LocationAreasResp{}, err
	}
	return locations, nil
}

func (c *Client) GetLocationArea(name string) (LocationAreaDetails, error) {
	url := fmt.Sprintf("%s/location-area/%s", baseURL, name)
	body, err := c.getBody(url)
	if err != nil {
		return LocationAreaDetails{}, err
	}

	var payload LocationAreaDetails
	if err := json.Unmarshal(body, &payload); err != nil {
		return LocationAreaDetails{}, err
	}

	return payload, nil
}

type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type PokemonType struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type Pokemon struct {
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
}

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", baseURL, name)
	body, err := c.getBody(url)
	if err != nil {
		return Pokemon{}, err
	}

	var payload Pokemon
	if err := json.Unmarshal(body, &payload); err != nil {
		return Pokemon{}, err
	}

	return payload, nil
}
