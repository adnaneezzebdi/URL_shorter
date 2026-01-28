package repositories

import (
	"database/sql"
	"errors"
	"log/slog"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) InsertShortURL(code string, originalUrl string) (string, error) {
	const query = `INSERT INTO short_urls (code, original_url) VALUES($1,$2)`
	if originalUrl == "" {
		slog.Error("originalUrl non può essere vuoto", "original_url", originalUrl)
		return "", errors.New("originalUrl non può essere vuoto")
	}

	_, err := r.db.Exec(query, code, originalUrl)
	if err != nil {
		slog.Error("errore query insert url", "error", err, "user_card_id", code, "original_url", originalUrl)
		return "", err
	}

	if err == 23505 {
		return "", nil
	}
	return code, nil
}

func GetByCode() {

}

func IncrementAccess() {

}

func GetStats() {

}
