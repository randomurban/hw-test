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
	br := bufio.NewReader(r)
	for i := 0; i < len(result); i++ {
		line, errRead := br.ReadString('\n')
		if errRead != nil {
			if errRead == io.EOF {
				break
			} else {
				return result, fmt.Errorf("read error: %w", errRead)
			}
		}
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.HasSuffix(user.Email, "."+domain) {
			name := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[name] += 1
		}
	}
	return result, nil
}
