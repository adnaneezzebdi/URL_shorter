package services

import (
	"strconv"
	"strings"
	"urlShorter/internal/repositories"

	"github.com/jxskiss/base62"
	"golang.org/x/crypto/bcrypt"
)

func CreateShortUrl(url string) (string, error) {
	if url == "" {
		return "", repositories.ErrEmptyUrl
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", repositories.ErrInvalidURL
	}
	if len(url) > 2048 {
		return "", repositories.ErrInvalidURL
	}

	var salt, attempt = "shorturl", 0
	for {
		s := url + ":" + salt + ":" + strconv.Itoa(attempt)
		attempt++
		hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost) //da finire
		if err == nil {
			return base62.EncodeToString(hash)[:8], nil
		}
		if attempt > 100 {
			return "", repositories.ErrTooManyAttempts
		}
	}

}
