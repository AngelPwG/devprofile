package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/AngelPwG/devprofile/models"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func New(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		log.Print("Error opening database: ", err)
		return nil, err
	}

	sqlStmt := `
		create table if not exists profiles (id integer primary key autoincrement, github_user text unique not null, name text, avatar_url text, bio text, followers int, following int, public_repos int, language text, pokemon text, pokemon_img text, created_at text not null, updated_at text not null);
		create table if not exists repositories (id integer primary key autoincrement, profile_id integer not null, name text not null, language text, foreign key (profile_id) references profiles(id) on delete cascade);
		create table if not exists audit_log (id integer primary key autoincrement, event text not null, resource text not null, author_ip text not null, timestamp text not null);
	`
	_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Print("Error whle creating the tables.")
		return nil, err
	}
	return &DB{conn: db}, nil
}

func (d *DB) InsertProfile(p models.Profile) error {
	_, err := d.conn.Exec("insert into profiles (github_user, name, avatar_url, bio, followers, following, public_repos, language, pokemon, pokemon_img, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", p.GithubUser, p.Name, p.AvatarURL, p.Bio, p.Followers, p.Following, p.PublicRepos, p.Language, p.Pokemon, p.PokemonImg, p.CreatedAt, p.UpdatedAt)
	return err
}

func (d *DB) GetProfiles() ([]models.Profile, error) {
	rows, err := d.conn.Query("select * from profiles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var p models.Profile
		if err := rows.Scan(&p.ID, &p.GithubUser, &p.Name, &p.AvatarURL, &p.Bio, &p.Followers, &p.Following, &p.PublicRepos, &p.Language, &p.Pokemon, &p.PokemonImg, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (d *DB) GetProfile(user string) (models.Profile, error) {
	var p models.Profile
	err := d.conn.QueryRow("select * from profiles where github_user = ?", user).Scan(&p.ID, &p.GithubUser, &p.Name, &p.AvatarURL, &p.Bio, &p.Followers, &p.Following, &p.PublicRepos, &p.Language, &p.Pokemon, &p.PokemonImg, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (d *DB) UpdateProfile(p models.Profile) error {
	_, err := d.conn.Exec("update profiles set name = ?, avatar_url = ?, bio = ?, followers = ?, following = ?, public_repos = ?, language = ?, pokemon = ?, pokemon_img = ?, updated_at = ? where github_user = ?", p.Name, p.AvatarURL, p.Bio, p.Followers, p.Following, p.PublicRepos, p.Language, p.Pokemon, p.PokemonImg, p.UpdatedAt, p.GithubUser)
	return err
}

func (d *DB) DeleteProfile(user string) error {
	_, err := d.conn.Exec("delete from profiles where github_user = ?", user)
	return err
}

func (d *DB) InsertRepositories(repos []models.Repository, id int) error {
	for _, repo := range repos {
		_, err := d.conn.Exec("insert into repositories (profile_id, name, language) values (?, ?, ?)", id, repo.Name, repo.Language)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) GetRepositories(id int) ([]models.Repository, error) {
	var repos []models.Repository
	rows, err := d.conn.Query("select * from repositories where profile_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r models.Repository
		if err := rows.Scan(&r.ID, &r.ProfileID, &r.Name, &r.Language); err != nil {
			return nil, err
		}
		repos = append(repos, r)
	}
	return repos, nil
}

func (d *DB) DeleteRepositories(id int) error {
	_, err := d.conn.Exec("delete from repositories where profile_id = ?", id)
	return err
}

func (d *DB) InsertAuditLog(event, resource, ip string) error {
	_, err := d.conn.Exec("insert into audit_log (event, resource, author_ip, timestamp) values (?, ?, ?, ?)", event, resource, ip, time.Now())
	return err
}

func (d *DB) GetAuditLogs() ([]models.AuditLog, error) {
	var logs []models.AuditLog
	rows, err := d.conn.Query("select * from audit_log")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.Event, &l.Resource, &l.AuthorIP, &l.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
