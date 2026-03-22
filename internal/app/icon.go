package app

import (
	"encoding/base64"

	"fyne.io/fyne/v2"
)

// Twemoji key (U+1F511), CC-BY 4.0 — https://twemoji.twitter.com/
const twemojiKeyPNG = "iVBORw0KGgoAAAANSUhEUgAAAEgAAABICAMAAABiM0N1AAAAM1BMVEVHcEzBaU/BaU/BaU/BaU/BaU/BaU/BaU/BaU/BaU/BaU/BaU/BaU/BaU/BaU/BaU8TO/0hAAAAEXRSTlMAMGCAr7+PcECf7//fEM8gUNW42NsAAAFLSURBVHgB7dYHgqQgEIXhR1LgUdj3v+zkPIhld23e7wC/gYgznA8xPVjWjOvlUvlB9LiKb/yqBsFZrnOkepwTuCcK9CRyX3P6TudMdTYdfanwSBMorBxrXuAKnyQc2yqHuuCR55MVhxLHtk8TwooOZI4lvNj4JOBA5FjAK6peaeOO8iVEf92QseHFhc8iphL3rF9nq2CGnJe2zlcZE44TbQmF7wImLtQrmAjUS/9Dt4Uy9RZMbNQLmKlUy5iJVMOUp1bElFQqXTBXqNMw5yp1vFGnGXWYjTqLUadbdeRP6qBRKdisseTO32Jr+Z7xwOlQdRAf+a4vDsfioPPEXcKjNQs0yrijZ9iRdRVtRzUP/3fGmlEHVh00ow6iQcfFFP7GeYj02417N+pgNeqgG3Vg1clGHWxGHaAadRCMOpBm0wFcV3R01kT2ILjBPVocWyZYgdJoAAAAAElFTkSuQmCC"

var appIconResource fyne.Resource

func init() {
	b, err := base64.StdEncoding.DecodeString(twemojiKeyPNG)
	if err == nil && len(b) > 0 {
		appIconResource = fyne.NewStaticResource("app_icon.png", b)
	}
}

func applyAppIcon(a fyne.App) {
	if appIconResource != nil {
		a.SetIcon(appIconResource)
	}
}
