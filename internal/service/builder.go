package services

import (
	"fmt"

	models "github.com/AngelPwG/devprofile/internal/domain"
)

func BuildProfile(username string) (*models.Profile, []models.Repository, error) {
	profile, repos, err := GetRepos(username)
	if err != nil {
		return nil, nil, fmt.Errorf("builder: github: %w", err)
	}

	language := DominantLanguage(repos)
	profile.Language = language

	pokemonName := LanguageToPokemon(language)
	profile.Pokemon = pokemonName

	spriteURL, err := GetPokemonSprite(pokemonName)
	if err != nil {
		return nil, nil, fmt.Errorf("builder: pokeapi: %w", err)
	}
	profile.PokemonImg = spriteURL

	return profile, repos, nil
}
