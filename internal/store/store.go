package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
)

// Account holds one TOTP secret (Base32) and a display name.
type Account struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

// Data is the on-disk format next to the executable.
type Data struct {
	Accounts []Account `json:"accounts"`
}

func dataDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// AccountsPath returns the path to accounts.json beside the binary.
func AccountsPath() (string, error) {
	dir, err := dataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "accounts.json"), nil
}

// Load reads accounts.json; missing file yields empty data (no error).
func Load() (*Data, error) {
	p, err := AccountsPath()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Data{Accounts: nil}, nil
		}
		return nil, err
	}
	var d Data
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, fmt.Errorf("accounts.json: %w", err)
	}
	for i := range d.Accounts {
		d.Accounts[i].Name = strings.TrimSpace(d.Accounts[i].Name)
		d.Accounts[i].Secret = strings.TrimSpace(strings.ToUpper(d.Accounts[i].Secret))
	}
	return &d, nil
}

// Save writes accounts.json atomically.
func Save(d *Data) error {
	p, err := AccountsPath()
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, b, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

// ValidateSecret checks Base32 and that TOTP can be generated.
func ValidateSecret(secret string) error {
	secret = strings.TrimSpace(strings.ToUpper(secret))
	if secret == "" {
		return errors.New("пустой секрет")
	}
	_, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return fmt.Errorf("некорректный Base32 или секрет: %w", err)
	}
	return nil
}
