package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const SEARCH_URL = "https://api.boardgameatlas.com/api/search"

// Like a class
type BoardgameAtlas struct {
	// 'members'
	clientId string
}

// Game
type Game struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Price         string `json:"price"`
	YearPublished uint   `json:"year_published"`
	Description   string `json:"description"`
	Url           string `json:"official_url"`
	ImageUrl      string `json:"image_url"`
	RulesUrl      string `json:"rules_url"`
}

type SearchResult struct {
	Games []Game `json:"games"`
	Count uint   `json:"count"`
}

// 'Method' in BoardgameAtlas
func (b BoardgameAtlas) Search(ctx context.Context, query string, limit uint, skip uint) (*SearchResult, error) {

	// Create HTTP client
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, SEARCH_URL, nil)
	// Check if there is any error
	if nil != err {
		// return an error objec
		return nil, fmt.Errorf("cannot create HTTP client: %v", err)
	}

	//Get the query string object
	qs := req.URL.Query()
	// Populate the URL with query params
	qs.Add("name", query)
	qs.Add("limit", fmt.Sprintf("%d", limit))
	qs.Add("skip", strconv.Itoa(int(skip)))
	qs.Add("client_id", b.clientId)
	//Encode the query params, add it back to the request
	req.URL.RawQuery = qs.Encode()

	//fmt.Printf("URL = %s\n", req.URL.String())

	// Make the call
	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		return nil, fmt.Errorf("cannot create HTTP client for invocation: %v", err)
	}
	// HTTP status code >= 400 is error
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error HTTP status: %s", resp.Status)
	}

	var result SearchResult
	// Deserialize the JSON payload to struct
	if err := json.NewDecoder(resp.Body).Decode(&result); nil != err {
		return nil, fmt.Errorf("cannot deserialize JSON payload: %v", err)
	}

	return &result, nil
}

// Functions as a constructor
func New(clientId string) BoardgameAtlas {
	return BoardgameAtlas{clientId: clientId}
}
