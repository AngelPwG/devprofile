package services

import (
	"devprofile/domain"
	"fmt"
)

func BuildProfile(username string) (*domain.Profile, []domain.Repository, error) {
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