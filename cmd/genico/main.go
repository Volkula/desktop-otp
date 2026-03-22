// genico builds build/windows/app.ico from the PNG embedded in internal/app/icon.go (Windows Explorer uses the .ico embedded via rsrc, not Fyne at runtime).
package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"regexp"

	ico "github.com/fyne-io/image/ico"
	"github.com/nfnt/resize"
)

func main() {
	outPath := flag.String("out", "build/windows/app.ico", "output .ico path")
	iconGo := flag.String("icongo", "internal/app/icon.go", "path to icon.go with twemojiKeyPNG")
	flag.Parse()

	src, err := os.ReadFile(*iconGo)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`twemojiKeyPNG = "([^"]+)"`)
	m := re.FindSubmatch(src)
	if len(m) < 2 {
		panic("twemojiKeyPNG not found")
	}
	raw, err := base64.StdEncoding.DecodeString(string(m[1]))
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		panic(err)
	}
	// 256×256: ICO directory uses 0 for “256” per Windows convention; uint8(256)==0.
	sized := resize.Resize(256, 256, img, resize.Lanczos3)
	var buf bytes.Buffer
	if err := ico.Encode(&buf, sized); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Dir(*outPath), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(*outPath, buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}
