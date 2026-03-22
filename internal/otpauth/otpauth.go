package otpauth

import (
	"errors"
	"net/url"
	"strings"
)

// ParseURI разбирает otpauth://totp/... или otpauth://hotp/....
func ParseURI(raw string) (name string, secret string, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", errors.New("пустая строка")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", err
	}
	if u.Scheme != "otpauth" {
		return "", "", errors.New("нужна ссылка otpauth://")
	}
	secret = strings.TrimSpace(u.Query().Get("secret"))
	if secret == "" {
		return "", "", errors.New("в ссылке нет параметра secret")
	}
	path := strings.TrimPrefix(u.Path, "/")
	if i := strings.Index(path, ":"); i >= 0 {
		name = strings.TrimSpace(path[i+1:])
	} else {
		name = strings.TrimSpace(path)
	}
	if name == "" {
		name = "TOTP"
	}
	return name, secret, nil
}
