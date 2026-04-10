package models

type Profile struct {
	ID          int    `json:"id"`
	GithubUser  string `json:"github_user"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	PublicRepos int    `json:"public_repos"`
	Language    string `json:"language"`
	Pokemon     string `json:"pokemon"`
	PokemonImg  string `json:"pokemon_img"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
type Repository struct {
	ID        int    `json:"id"`
	ProfileID int    `json:"profile_id"`
	Name      string `json:"name"`
	Language  string `json:"language"`
}
type AuditLog struct {
	ID        int    `json:"id"`
	Event     string `json:"event"`
	Resource  string `json:"resource"`
	AuthorIP  string `json:"author_ip"`
	Timestamp string `json:"timestamp"`
}
