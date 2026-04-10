package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type pokeResponse struct {
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

func GetPokemonSprite(name string) (string, error) {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)
	if err != nil {
		return "", fmt.Errorf("pokeapi: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("pokeapi: status %d", resp.StatusCode)
	}

	var result pokeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("pokeapi: decode: %w", err)
	}

	if result.Sprites.FrontDefault == "" {
		return "", fmt.Errorf("pokeapi: sprite not found for %s", name)
	}

	return result.Sprites.FrontDefault, nil
}