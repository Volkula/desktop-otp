package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// File holds optional hotkey overrides (config.json next to exe).
type File struct {
	// Modifiers: битовая маска MOD_* (Windows): ALT=1, CTRL=2, SHIFT=4, WIN=8.
	Modifiers uint32 `json:"modifiers"`
	// VK — виртуальный код клавиши (например 79 = 'O').
	VK uint32 `json:"vk"`
}

func dir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// Path returns path to config.json.
func Path() (string, error) {
	d, err := dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "config.json"), nil
}

// Load reads config.json; missing file returns defaults (no error).
func Load() (*File, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultConfig(), nil
		}
		return nil, err
	}
	var c File
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	if c.Modifiers == 0 || c.VK == 0 {
		return defaultConfig(), nil
	}
	return &c, nil
}

func defaultConfig() *File {
	// Ctrl+Shift+O — редко пересекается с IDE по умолчанию.
	return &File{Modifiers: 0x0002 | 0x0004, VK: 0x4F}
}

// Маски модификаторов (как в WinAPI RegisterHotKey).
const (
	ModAlt     uint32 = 0x0001
	ModControl uint32 = 0x0002
	ModShift   uint32 = 0x0004
	ModWin     uint32 = 0x0008
)

// HotkeyLabel — строка для подсказки в UI.
func (f *File) HotkeyLabel() string {
	var b []byte
	if f.Modifiers&ModAlt != 0 {
		b = append(b, "Alt+"...)
	}
	if f.Modifiers&ModControl != 0 {
		b = append(b, "Ctrl+"...)
	}
	if f.Modifiers&ModShift != 0 {
		b = append(b, "Shift+"...)
	}
	if f.Modifiers&ModWin != 0 {
		b = append(b, "Win+"...)
	}
	if f.VK >= 'A' && f.VK <= 'Z' {
		b = append(b, byte(f.VK))
	} else {
		b = append(b, fmt.Sprintf("VK 0x%X", f.VK)...)
	}
	return string(b)
}
