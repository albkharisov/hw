package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string `json:"email"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var user User
	result := make(DomainStat)
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// no need to clear 'user' variable, because we have 1 field (TODO: test it)
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
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
