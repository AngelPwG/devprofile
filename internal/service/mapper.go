package services

import (
	"strings"

	models "github.com/AngelPwG/devprofile/internal/domain"
)

var languageToPokemon = map[string]string{
	"javascript": "Pikachu",
	"typescript": "Raichu",
	"c":          "Mew",
	"python":     "Ditto",
	"kotlin":     "Gallade",
	"swift":      "Gardevoir",
	"shell":      "Klingklang",
	"bash":       "Klingklang",
	"css":        "Kecleon",
	"assembly":   "Arceus",
	"rust":       "Kingler",
	"c++":        "Mewtwo",
	"go":         "Greedent",
	"ruby":       "Sableye",
	"html":       "Jinx",
	"php":        "Phanpy",
	"java":       "Metapod",
	"c#":         "Kakuna",
	"lua":        "Ninjask",
}

const defaultPokemon = "Unown"

func DominantLanguage(repos []models.Repository) string {
	if len(repos) == 0 {
		return ""
	}

	counts := make(map[string]int)
	for _, repo := range repos {
		if repo.Language == "" {
			continue
		}
		counts[strings.ToLower(repo.Language)]++
	}

	if len(counts) == 0 {
		return ""
	}

	dominant := ""
	max := 0
	for lang, count := range counts {
		if count > max {
			max = count
			dominant = lang
		}
	}

	return dominant
}

func LanguageToPokemon(language string) string {
	normalized := strings.ToLower(strings.TrimSpace(language))
	pokemon, ok := languageToPokemon[normalized]
	if !ok {
		return defaultPokemon
	}
	return pokemon
}
