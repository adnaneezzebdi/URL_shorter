package services

import (
	"strconv"
	"strings"

	"urlShorter/internal/repositories"

	"github.com/jxskiss/base62"
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
		return "", repositories.ErrEmptyUrl
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", repositories.ErrInvalidURL
	}

	if len(url) > 2048 {
		return "", repositories.ErrInvalidURL
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {

		raw := url + ":" + salt + ":" + strconv.Itoa(attempt)

		encoded := base62.EncodeToString([]byte(raw))
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

	return "", repositories.ErrTooManyAttempts
}
