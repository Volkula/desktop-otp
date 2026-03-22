//go:build !windows

package hotkey

// Start на не-Windows не регистрирует хоткей (удобно для проверки сборки без WinAPI).
func Start(modifiers, vk uint32, pressed chan<- struct{}) (stop func(), err error) {
	return func() {}, nil
}
