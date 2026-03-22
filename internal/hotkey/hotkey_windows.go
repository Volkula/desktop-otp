//go:build windows

package hotkey

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	wmHotkey  = 0x0312
	pmRemove  = 0x0001
	hotkeyWin = 1
)

// Start registers a global hotkey (hwnd=0 → сообщения потоку) и шлёт сигнал в pressed.
func Start(modifiers, vk uint32, pressed chan<- struct{}) (stop func(), err error) {
	errCh := make(chan error, 1)
	stopCh := make(chan struct{})
	done := make(chan struct{})

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer close(done)

		user32 := windows.NewLazyDLL("user32.dll")
		procRegisterHotKey := user32.NewProc("RegisterHotKey")
		procUnregisterHotKey := user32.NewProc("UnregisterHotKey")
		procPeekMessage := user32.NewProc("PeekMessageW")

		r1, _, _ := procRegisterHotKey.Call(
			0,
			uintptr(hotkeyWin),
			uintptr(modifiers),
			uintptr(vk),
		)
		if r1 == 0 {
			errCh <- fmt.Errorf("RegisterHotKey: %w", windows.GetLastError())
			return
		}
		defer func() {
			procUnregisterHotKey.Call(0, uintptr(hotkeyWin))
		}()

		errCh <- nil

		var msg struct {
			Hwnd    uintptr
			Message uint32
			Padding uint32
			WParam  uintptr
			LParam  uintptr
			Time    uint32
			Pt      struct{ X, Y int32 }
		}

		for {
			select {
			case <-stopCh:
				return
			default:
			}

			r, _, _ := procPeekMessage.Call(
				uintptr(unsafe.Pointer(&msg)),
				0,
				0,
				0,
				pmRemove,
			)
			if r != 0 {
				if msg.Message == wmHotkey {
					select {
					case pressed <- struct{}{}:
					default:
					}
				}
			} else {
				time.Sleep(25 * time.Millisecond)
			}
		}
	}()

	if e := <-errCh; e != nil {
		return nil, e
	}

	return func() {
		close(stopCh)
		<-done
	}, nil
}
