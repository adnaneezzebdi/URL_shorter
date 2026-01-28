package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

var (
	ErrCodeCollision   = errors.New("short code collision")
	ErrInvalidURL      = errors.New("original url cannot be empty")
	ErrEmptyCode       = errors.New("code cannot be empty")
	ErrEmptyUrl        = errors.New("url cannot be empty")
	ErrTooManyAttempts = errors.New("too many attempts to generate unique code")
)

func (r *Repo) InsertShortURL(code, originalURL string) (string, error) {
	if originalURL == "" {
		return "", ErrInvalidURL
	}

	const insertQuery = `
		INSERT INTO short_urls (code, original_url)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(insertQuery, code, originalURL)
	if err == nil {
		return code, nil
	}

	pgErr, ok := err.(*pq.Error)
	if !ok {
		return "", err
	}

	if pgErr.Code != "23505" {
		return "", err
	}

	switch pgErr.Constraint {

	case "short_urls_original_url_key":
		// idempotenza: stessa URL → ritorna code esistente
		existingCode, err := r.getCodeByOriginalURL(originalURL)
		if err != nil {
			return "", err
		}
		return existingCode, nil

	case "short_urls_code_key":
		// collisione vera
		return "", ErrCodeCollision
	}

	return "", err
}

func (r *Repo) getCodeByOriginalURL(originalURL string) (string, error) {
	const query = `
		SELECT code
		FROM short_urls
		WHERE original_url = $1
	`

	var code string
	err := r.db.QueryRow(query, originalURL).Scan(&code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		return "", err
	}

	return code, nil
}

type ShortUrl struct {
	Code         string
	OriginalUrl  string
	CreatedAt    time.Time
	AccessCount  int64
	LastAccessAt time.Time
}

func (r *Repo) GetByCode(code string) (*ShortUrl, error) {
	if code == "" {
		return nil, ErrEmptyCode
	}
	const query = `SELECT code, 
	original_url, 
	created_at, 
	access_count, 
	last_access_at FROM short_urls WHERE code = $1`

	var ResolvedCode ShortUrl
	Now := time.Now()
	ResolvedCode.LastAccessAt = Now

	err := r.db.QueryRow(query, code).Scan(&ResolvedCode.Code, &ResolvedCode.OriginalUrl, &ResolvedCode.CreatedAt, &ResolvedCode.AccessCount, &ResolvedCode.LastAccessAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return &ShortUrl{
		Code:         ResolvedCode.Code,
		OriginalUrl:  ResolvedCode.OriginalUrl,
		CreatedAt:    ResolvedCode.CreatedAt,
		AccessCount:  ResolvedCode.AccessCount,
		LastAccessAt: ResolvedCode.LastAccessAt,
	}, nil
}

func (r *Repo) IncrementAccess(code string, accessTime time.Time) error {
	if code == "" {
		return ErrEmptyCode
	}

	const query = `
		UPDATE short_urls
		SET access_count = access_count + 1,
			last_access_at = $2
		WHERE code = $1
	`
	psql, err := r.db.Exec(query, code, accessTime)
	if err != nil {
		return err
	}
	if rwAff, _ := psql.RowsAffected(); rwAff == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (r *Repo) GetStats(code string) (*ShortUrl, error) {
	if code == "" {
		return nil, ErrEmptyCode
	}
	const query = `SELECT code, 
	original_url, 
	created_at, 
	access_count, 
	last_access_at FROM short_urls WHERE code = $1`

	var ResolvedCode ShortUrl
	Now := time.Now()
	ResolvedCode.LastAccessAt = Now

	err := r.db.QueryRow(query, code).Scan(&ResolvedCode.Code, &ResolvedCode.OriginalUrl, &ResolvedCode.CreatedAt, &ResolvedCode.AccessCount, &ResolvedCode.LastAccessAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return &ShortUrl{
		Code:         ResolvedCode.Code,
		OriginalUrl:  ResolvedCode.OriginalUrl,
		CreatedAt:    ResolvedCode.CreatedAt,
		AccessCount:  ResolvedCode.AccessCount,
		LastAccessAt: ResolvedCode.LastAccessAt,
	}, nil
}
