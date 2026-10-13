package pokeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/BashCtl/pokedex-cli/internal/pokecache"
)

func newTestClient(server *httptest.Server) *Client {
	return &Client{
		httpClient: *server.Client(),
		cache:      pokecache.NewCache(time.Hour),
		baseURL:    server.URL,
	}
}

func TestListLocationAreas(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/location-area" {
			t.Errorf("unexpected request path: %s", r.URL.Path)
		}
		fmt.Fprint(w, `{"count":2,"next":"next-page","previous":null,"results":[{"name":"canalave-city-area","url":"area-url"}]}`)
	}))
	defer server.Close()

	locations, err := newTestClient(server).ListLocationAreas(nil)
	if err != nil {
		t.Fatalf("ListLocationAreas returned error: %v", err)
	}
	if locations.Count != 2 {
		t.Errorf("Count = %d, want 2", locations.Count)
	}
	if locations.Next == nil || *locations.Next != "next-page" {
		t.Errorf("Next = %v, want next-page", locations.Next)
	}
	if locations.Previous != nil {
		t.Errorf("Previous = %v, want nil", locations.Previous)
	}
	if len(locations.Results) != 1 || locations.Results[0].Name != "canalave-city-area" {
		t.Errorf("Results = %+v, want canalave-city-area", locations.Results)
	}
}

func TestGetLocationArea(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/location-area/canalave-city-area" {
			t.Errorf("unexpected request path: %s", r.URL.Path)
		}
		fmt.Fprint(w, `{"name":"canalave-city-area","pokemon_encounters":[{"pokemon":{"name":"tentacool","url":"pokemon-url"}}]}`)
	}))
	defer server.Close()

	area, err := newTestClient(server).GetLocationArea("canalave-city-area")
	if err != nil {
		t.Fatalf("GetLocationArea returned error: %v", err)
	}
	if area.Name != "canalave-city-area" {
		t.Errorf("Name = %q, want canalave-city-area", area.Name)
	}
	if len(area.PokemonEncounters) != 1 || area.PokemonEncounters[0].Pokemon.Name != "tentacool" {
		t.Errorf("PokemonEncounters = %+v, want one encounter for tentacool", area.PokemonEncounters)
	}
}

func TestGetPokemon(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pokemon/pikachu" {
			t.Errorf("unexpected request path: %s", r.URL.Path)
		}
		fmt.Fprint(w, `{"name":"pikachu","base_experience":112,"height":4,"weight":60,"stats":[{"base_stat":35,"stat":{"name":"hp"}}],"types":[{"slot":1,"type":{"name":"electric"}}]}`)
	}))
	defer server.Close()

	pokemon, err := newTestClient(server).GetPokemon("pikachu")
	if err != nil {
		t.Fatalf("GetPokemon returned error: %v", err)
	}
	if pokemon.Name != "pikachu" || pokemon.BaseExperience != 112 || pokemon.Height != 4 || pokemon.Weight != 60 {
		t.Errorf("unexpected Pokemon fields: %+v", pokemon)
	}
	if len(pokemon.Stats) != 1 || pokemon.Stats[0].Stat.Name != "hp" || pokemon.Stats[0].BaseStat != 35 {
		t.Errorf("unexpected Pokemon stats: %+v", pokemon.Stats)
	}
	if len(pokemon.Types) != 1 || pokemon.Types[0].Type.Name != "electric" || pokemon.Types[0].Slot != 1 {
		t.Errorf("unexpected Pokemon types: %+v", pokemon.Types)
	}
}

func TestGetPokemonUsesCache(t *testing.T) {
	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		fmt.Fprint(w, `{"name":"pikachu","base_experience":112}`)
	}))
	defer server.Close()

	client := newTestClient(server)
	for i := 0; i < 2; i++ {
		if _, err := client.GetPokemon("pikachu"); err != nil {
			t.Fatalf("GetPokemon call %d returned error: %v", i+1, err)
		}
	}
	if got := requestCount.Load(); got != 1 {
		t.Errorf("request count = %d, want 1", got)
	}
}

func TestGetPokemonReturnsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	if _, err := newTestClient(server).GetPokemon("missingno"); err == nil {
		t.Fatal("GetPokemon returned nil error for HTTP 404")
	}
}
