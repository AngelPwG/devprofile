package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	models "github.com/AngelPwG/devprofile/internal/domain"
)

type ghRepo struct {
	Name     string `json:"name"`
	Language struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`
}

type ghUser struct {
	Name      string `json:"name"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatarUrl"`
	Bio       string `json:"bio"`
	Followers struct {
		TotalCount int `json:"totalCount"`
	} `json:"followers"`
	Following struct {
		TotalCount int `json:"totalCount"`
	} `json:"following"`
	Repositories struct {
		TotalCount int      `json:"totalCount"`
		Nodes      []ghRepo `json:"nodes"`
	} `json:"repositories"`
	PinnedItems struct {
		Nodes []ghRepo `json:"nodes"`
	} `json:"pinnedItems"`
	ContributionsCollection struct {
		TotalCommitContributions int `json:"totalCommitContributions"`
	} `json:"contributionsCollection"`
}

type ghResponse struct {
	Data struct {
		User ghUser `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

const reposQuery = `
query($login: String!) {
  user(login: $login) {
    name
    login
    avatarUrl
    bio
    followers { totalCount }
    following  { totalCount }
    repositories(first: 5, orderBy: {field: CREATED_AT, direction: DESC}) {
      totalCount
      nodes {
        name
        primaryLanguage { name }
      }
    }
    pinnedItems(first: 6, types: [REPOSITORY]) {
      nodes {
        ... on Repository {
          name
          primaryLanguage { name }
        }
      }
    }
    contributionsCollection {
      totalCommitContributions
    }
  }
}`

func GetRepos(username string) (*models.Profile, []models.Repository, error) {
	body := map[string]interface{}{
		"query": reposQuery,
		"variables": map[string]string{
			"login": username,
		},
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, nil, fmt.Errorf("github: marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, fmt.Errorf("github: create request: %w", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("github: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("github: status %d", resp.StatusCode)
	}

	var ghResp ghResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, nil, fmt.Errorf("github: decode response: %w", err)
	}

	if len(ghResp.Errors) > 0 {
		msg := ghResp.Errors[0].Message
		if strings.Contains(msg, "Could not resolve to a User") {
			return nil, nil, fmt.Errorf("user not found")
		}		
		return nil, nil, fmt.Errorf("github: api error: %s", msg)
	}

	u := ghResp.Data.User

	repos := make([]models.Repository, 0)
	if len(u.PinnedItems.Nodes) > 0 {
		for _, node := range u.PinnedItems.Nodes {
			repos = append(repos, models.Repository{
				Name:     node.Name,
				Language: node.Language.Name,
			})
		}
	} else {
		limit := len(u.Repositories.Nodes)
		if limit > 5 {
			limit = 5
		}
		for _, node := range u.Repositories.Nodes[:limit] {
			repos = append(repos, models.Repository{
				Name:     node.Name,
				Language: node.Language.Name,
			})
		}
	}

	profile := &models.Profile{
		GithubUser:  u.Login,
		Name:        u.Name,
		AvatarURL:   u.AvatarURL,
		Bio:         u.Bio,
		Followers:   u.Followers.TotalCount,
		Following:   u.Following.TotalCount,
		PublicRepos: u.Repositories.TotalCount,
	}

	return profile, repos, nil
}
