package services

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"urlShorter/internal/repositories"

	"github.com/jxskiss/base62"
)

var (
	ErrInvalidURL      = errors.New("original url is not valid or is empty")
	ErrTooManyAttempts = errors.New("too many attempts to generate unique code")
	ErrNotFound        = errors.New("no records returned from db")
	ErrIncrementAccess = errors.New("failed to increment access counter")
)

const (
	codeLength  = 8
	maxAttempts = 100
	salt        = "shorturl"
)

type ShortUrlService struct {
	repo *repositories.Repo
}

func NewShortUrlService(repo *repositories.Repo) *ShortUrlService {
	return &ShortUrlService{repo: repo}
}

func (s *ShortUrlService) CreateShortUrl(url string) (string, error) {
	if url == "" {
		return "", ErrInvalidURL
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", ErrInvalidURL
	}

	if len(url) > 2048 {
		return "", ErrInvalidURL
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {

		raw := url + ":" + salt + ":" + strconv.Itoa(attempt)
		hash := sha256.New()
		hash.Write([]byte(raw))
		encoded := base62.EncodeToString(hash.Sum(nil))
		if len(encoded) < codeLength {
			continue
		}

		code := encoded[:codeLength]

		insertedCode, err := s.repo.InsertShortURL(code, url)
		if err == nil {
			return insertedCode, nil
		}

		if err == repositories.ErrCodeCollision {
			continue
		}

		return "", err
	}

	return "", ErrTooManyAttempts
}

func (s *ShortUrlService) ResolveShortUrl(code string) (string, error) {

	trimCode := strings.TrimSpace(code)
	if trimCode == "" {
		return "", ErrInvalidURL
	}

	res, err := s.repo.GetByCode(trimCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		} else {
			return "", ErrInvalidURL
		}
	}

	err = s.repo.IncrementAccess(trimCode, time.Now().UTC())
	if err != nil {
		return "", ErrIncrementAccess
	}
	return res.OriginalUrl, nil
}
