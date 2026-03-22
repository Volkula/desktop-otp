//go:build windows

package autostart

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// registry value name (stable; do not change or users get duplicate entries).
const runValueName = "DesktopOTP"

// SetEnabled adds or removes HKCU\Software\Microsoft\Windows\CurrentVersion\Run entry for this executable.
func SetEnabled(enabled bool) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.Abs(exe)
	if err != nil {
		return err
	}

	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	if !enabled {
		err := k.DeleteValue(runValueName)
		if err != nil {
			if errno, ok := err.(syscall.Errno); ok && errno == windows.ERROR_FILE_NOT_FOUND {
				return nil
			}
			return err
		}
		return nil
	}

	val := exe
	if strings.ContainsAny(exe, " \t") {
		val = `"` + exe + `"`
	}
	return k.SetStringValue(runValueName, val)
}
