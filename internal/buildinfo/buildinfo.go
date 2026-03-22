package buildinfo

import "strings"

// Injected at link time, e.g.:
//   go build -ldflags "-X desctop-otp/internal/buildinfo.Version=1.0.0 ..."
var (
	Version   = "dev"
	BuildDate = ""
	Branch    = ""
)

// Value returns s trimmed, or an em dash when empty (for UI).
func Value(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "—"
	}
	return s
}
