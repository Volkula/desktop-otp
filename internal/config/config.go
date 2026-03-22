package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// File holds settings stored in config.json next to the executable.
type File struct {
	// Modifiers: WinAPI MOD_* mask: ALT=1, CTRL=2, SHIFT=4, WIN=8.
	Modifiers uint32 `json:"modifiers"`
	// VK — virtual-key code (e.g. 0x4F = 'O').
	VK uint32 `json:"vk"`
	// Language: "ru" or "en".
	Language string `json:"language,omitempty"`
	// DataDir: empty = same folder as exe; otherwise absolute path or path relative to exe dir.
	DataDir string `json:"data_dir,omitempty"`
	// MinimizeToTray: nil = default true (omit in old JSON).
	MinimizeToTray *bool `json:"minimize_to_tray,omitempty"`
	// DarkTheme: nil = default true (dark UI).
	DarkTheme *bool `json:"dark_theme,omitempty"`
	// HideAfterCopy: nil = default false — hide window after copying code from list.
	HideAfterCopy *bool `json:"hide_after_copy,omitempty"`
	// StartWithWindows: nil = default false — add exe to HKCU\...\Run (Windows only).
	StartWithWindows *bool `json:"start_with_windows,omitempty"`
}

func dir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// Path returns path to config.json beside the executable.
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
	def := defaultConfig()
	if c.Modifiers == 0 || c.VK == 0 {
		c.Modifiers = def.Modifiers
		c.VK = def.VK
	}
	if strings.TrimSpace(c.Language) == "" {
		c.Language = def.Language
	}
	return &c, nil
}

func defaultConfig() *File {
	t, f := true, false
	return &File{
		Modifiers:          ModControl | ModShift,
		VK:                 0x4F,
		Language:           "ru",
		DataDir:            "",
		MinimizeToTray:     &t,
		DarkTheme:          &t,
		HideAfterCopy:      &f,
		StartWithWindows:   &f,
	}
}

// MinimizeToTrayBool returns whether closing the window hides to tray (default true).
func (f *File) MinimizeToTrayBool() bool {
	if f.MinimizeToTray == nil {
		return true
	}
	return *f.MinimizeToTray
}

// SetMinimizeToTray sets the value to persist.
func (f *File) SetMinimizeToTray(v bool) {
	f.MinimizeToTray = &v
}

// DarkThemeBool returns whether to use the built-in dark theme (default true).
func (f *File) DarkThemeBool() bool {
	if f.DarkTheme == nil {
		return true
	}
	return *f.DarkTheme
}

// SetDarkTheme sets dark vs light theme flag.
func (f *File) SetDarkTheme(v bool) {
	f.DarkTheme = &v
}

// HideAfterCopyBool returns whether to hide the window after copying a code (default false).
func (f *File) HideAfterCopyBool() bool {
	if f.HideAfterCopy == nil {
		return false
	}
	return *f.HideAfterCopy
}

// SetHideAfterCopy sets hide-after-copy behavior.
func (f *File) SetHideAfterCopy(v bool) {
	f.HideAfterCopy = &v
}

// StartWithWindowsBool returns whether to register in Windows startup (Run key); default false.
func (f *File) StartWithWindowsBool() bool {
	if f.StartWithWindows == nil {
		return false
	}
	return *f.StartWithWindows
}

// SetStartWithWindows sets autostart for Windows (ignored on other OS at apply time).
func (f *File) SetStartWithWindows(v bool) {
	f.StartWithWindows = &v
}

// ResolveDataDir returns the absolute directory for accounts.json.
func (f *File) ResolveDataDir() (string, error) {
	exeDir, err := dir()
	if err != nil {
		return "", err
	}
	s := strings.TrimSpace(f.DataDir)
	if s == "" {
		return exeDir, nil
	}
	if filepath.IsAbs(s) {
		return filepath.Clean(s), nil
	}
	return filepath.Join(exeDir, filepath.Clean(s)), nil
}

// Save writes config.json next to the executable.
func Save(f *File) error {
	p, err := Path()
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, b, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

// Modifier masks (WinAPI RegisterHotKey).
const (
	ModAlt     uint32 = 0x0001
	ModControl uint32 = 0x0002
	ModShift   uint32 = 0x0004
	ModWin     uint32 = 0x0008
)

// Windows virtual-key codes for OEM keys (see winuser.h).
const (
	// VK_OEM_3 — US layout: ` ~ (same physical key for both characters).
	VKOEM3 uint32 = 0xC0
)

// HotkeyLabel builds a display string for the hotkey.
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
	} else if f.VK >= '0' && f.VK <= '9' {
		b = append(b, byte(f.VK))
	} else if f.VK == VKOEM3 {
		b = append(b, '`')
	} else {
		b = append(b, fmt.Sprintf("VK 0x%X", f.VK)...)
	}
	return string(b)
}

// ParseKeyString returns VK for one character: A–Z, 0–9, or ` / ~ (OEM3).
func ParseKeyString(s string) (vk uint32, ok bool) {
	s = strings.TrimSpace(s)
	if len(s) != 1 {
		return 0, false
	}
	r := rune(s[0])
	switch r {
	case '`', '~':
		return VKOEM3, true
	}
	if r >= 'a' && r <= 'z' {
		r = r - 'a' + 'A'
	}
	if r >= 'A' && r <= 'Z' {
		return uint32(r), true
	}
	if r >= '0' && r <= '9' {
		return uint32(r), true
	}
	return 0, false
}

// ModifiersFromBools builds MOD_* mask from checkboxes.
func ModifiersFromBools(alt, ctrl, shift, win bool) uint32 {
	var m uint32
	if alt {
		m |= ModAlt
	}
	if ctrl {
		m |= ModControl
	}
	if shift {
		m |= ModShift
	}
	if win {
		m |= ModWin
	}
	return m
}
