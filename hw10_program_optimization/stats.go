package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-json"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

var ErrEmptyDomainName = errors.New("empty domain name")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if len(domain) == 0 {
		return nil, ErrEmptyDomainName
	}
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var user User
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}
		if strings.HasSuffix(user.Email, "."+domain) {
			parts := strings.SplitN(user.Email, "@", 2)
			if len(parts) == 2 {
				name := strings.ToLower(parts[1])
				result[name]++
			}
		}
	}
	return result, nil
}
