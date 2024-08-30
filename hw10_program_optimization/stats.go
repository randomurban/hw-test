package hw10programoptimization

import (
	"bufio"
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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		var user User
		if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return result, fmt.Errorf("unmarshal error: %w", err)
		}
		result[i] = user
		if i < len(result) {
			i++
		} else {
			break
		}
	}
	return result, scanner.Err()
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.HasSuffix(user.Email, "."+domain) {
			name := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[name]++
		}
	}
	return result, nil
}
