package db

import (
	"database/sql"
	"log"
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
