package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string `json:"email"`
}

type DomainStat map[string]int

var (
	ErrDomain        = errors.New("Uncorrect domain name")
	ErrMalformedJSON = errors.New("Malformed JSON")
)

var re = regexp.MustCompile("\\A\\w+\\z")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if ok := re.MatchString(domain); !ok {
		return nil, ErrDomain
	}

	var user User
	result := make(DomainStat)
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		user.Email = user.Email[:0]
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, ErrMalformedJSON
		}
		if i := strings.Index(user.Email, "."+domain); i > -1 {
			index := strings.ToLower(user.Email[strings.IndexRune(user.Email, '@')+1:])
			result[index] = result[index] + 1
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading error: %w", err)
	}

	return result, nil
}
