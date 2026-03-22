package app

import (
	"strings"
	"time"

	"desctop-otp/internal/config"
	"desctop-otp/internal/hotkey"
	"desctop-otp/internal/otpauth"
	"desctop-otp/internal/store"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/pquerna/otp/totp"
)

// Run запускает UI, глобальный хоткей и автосохранение accounts.json рядом с exe.
func Run(cfg *config.File) {
	d, err := store.Load()
	if err != nil {
		panic(err)
	}

	a := fyneapp.NewWithID("desctop-otp")
	w := a.NewWindow("TOTP")
	w.Resize(fyne.NewSize(520, 480))

	hint := widget.NewLabel("Горячая клавиша: " + cfg.HotkeyLabel() + " — показать окно. Файл accounts.json лежит рядом с программой.")
	hint.Wrapping = fyne.TextWrapWord

	nameEnt := widget.NewEntry()
	nameEnt.SetPlaceHolder("Название (например GitHub)")
	secEnt := widget.NewEntry()
	secEnt.SetPlaceHolder("Секрет Base32")
	secEnt.Password = true

	importEnt := widget.NewEntry()
	importEnt.SetPlaceHolder("otpauth://totp/...")

	listBox := container.NewVBox()
	var codeLabels []*widget.Label

	var rebuildList func()
	rebuildList = func() {
		listBox.RemoveAll()
		codeLabels = codeLabels[:0]

		for i := range d.Accounts {
			i := i
			acc := d.Accounts[i]
			secret := strings.TrimSpace(strings.ToUpper(acc.Secret))

			codeLbl := widget.NewLabel("------")
			codeLbl.TextStyle = fyne.TextStyle{Monospace: true}

			copyBtn := widget.NewButton("Копировать", func() {
				c, err := totp.GenerateCode(secret, time.Now())
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				w.Clipboard().SetContent(c)
			})

			delBtn := widget.NewButton("Удалить", func() {
				d.Accounts = append(d.Accounts[:i], d.Accounts[i+1:]...)
				if err := store.Save(d); err != nil {
					dialog.ShowError(err, w)
					return
				}
				rebuildList()
			})

			row := container.NewHBox(
				widget.NewLabel(acc.Name),
				layout.NewSpacer(),
				codeLbl,
				copyBtn,
				delBtn,
			)
			listBox.Add(row)
			codeLabels = append(codeLabels, codeLbl)
		}
	}

	addBtn := widget.NewButton("Добавить", func() {
		name := strings.TrimSpace(nameEnt.Text)
		sec := strings.TrimSpace(secEnt.Text)
		if name == "" {
			dialog.ShowInformation("Пусто", "Укажите название.", w)
			return
		}
		if err := store.ValidateSecret(sec); err != nil {
			dialog.ShowError(err, w)
			return
		}
		d.Accounts = append(d.Accounts, store.Account{Name: name, Secret: strings.ToUpper(sec)})
		if err := store.Save(d); err != nil {
			dialog.ShowError(err, w)
			return
		}
		nameEnt.SetText("")
		secEnt.SetText("")
		rebuildList()
	})

	importBtn := widget.NewButton("Импорт URI", func() {
		raw := strings.TrimSpace(importEnt.Text)
		if raw == "" {
			dialog.ShowInformation("Пусто", "Вставьте otpauth:// ссылку.", w)
			return
		}
		n, sec, err := otpauth.ParseURI(raw)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if err := store.ValidateSecret(sec); err != nil {
			dialog.ShowError(err, w)
			return
		}
		d.Accounts = append(d.Accounts, store.Account{Name: n, Secret: strings.ToUpper(strings.TrimSpace(sec))})
		if err := store.Save(d); err != nil {
			dialog.ShowError(err, w)
			return
		}
		importEnt.SetText("")
		rebuildList()
	})

	pasteBtn := widget.NewButton("Из буфера", func() {
		raw := w.Clipboard().Content()
		if raw == "" {
			dialog.ShowInformation("Буфер", "Буфер обмена пуст.", w)
			return
		}
		importEnt.SetText(raw)
		n, sec, err := otpauth.ParseURI(raw)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if err := store.ValidateSecret(sec); err != nil {
			dialog.ShowError(err, w)
			return
		}
		d.Accounts = append(d.Accounts, store.Account{Name: n, Secret: strings.ToUpper(strings.TrimSpace(sec))})
		if err := store.Save(d); err != nil {
			dialog.ShowError(err, w)
			return
		}
		rebuildList()
	})

	form := container.NewVBox(
		widget.NewLabel("Новый аккаунт"),
		nameEnt,
		secEnt,
		addBtn,
		widget.NewSeparator(),
		widget.NewLabel("Импорт из Google Authenticator и др."),
		importEnt,
		container.NewHBox(importBtn, pasteBtn),
	)

	scroll := container.NewVScroll(listBox)
	scroll.SetMinSize(fyne.NewSize(500, 220))

	root := container.NewBorder(
		hint, form, nil, nil,
		scroll,
	)
	w.SetContent(root)

	rebuildList()

	go func() {
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()
		for range t.C {
			fyne.Do(func() {
				for i := range d.Accounts {
					if i >= len(codeLabels) {
						return
					}
					sec := strings.TrimSpace(strings.ToUpper(d.Accounts[i].Secret))
					c, err := totp.GenerateCode(sec, time.Now())
					if err != nil {
						codeLabels[i].SetText("ошибка")
						continue
					}
					codeLabels[i].SetText(c)
				}
			})
		}
	}()

	if desk, ok := a.(desktop.App); ok {
		menu := fyne.NewMenu("TOTP",
			fyne.NewMenuItem("Показать", func() {
				w.Show()
				w.RequestFocus()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Выход", func() { a.Quit() }),
		)
		desk.SetSystemTrayMenu(menu)
		w.SetCloseIntercept(func() {
			w.Hide()
		})
	}

	pressed := make(chan struct{}, 8)
	stopHK, hkErr := hotkey.Start(cfg.Modifiers, cfg.VK, pressed)
	if hkErr != nil {
		dialog.ShowError(hkErr, w)
	}

	go func() {
		for range pressed {
			fyne.Do(func() {
				w.Show()
				w.RequestFocus()
			})
		}
	}()

	defer func() {
		if stopHK != nil {
			stopHK()
		}
	}()

	if len(d.Accounts) > 0 {
		w.Hide()
	} else {
		w.Show()
	}
	a.Run()
}
