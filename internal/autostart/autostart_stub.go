//go:build !windows

package autostart

// SetEnabled updates autostart for the current executable; no-op outside Windows.
func SetEnabled(enabled bool) error {
	return nil
}
