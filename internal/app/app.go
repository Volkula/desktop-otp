package app

import (
	"errors"
	"strings"
	"time"

	"desctop-otp/internal/config"
	"desctop-otp/internal/hotkey"
	"desctop-otp/internal/locale"
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

// Run starts the UI, global hotkey, and tray.
func Run(cfg *config.File) {
	if err := store.InitFromConfig(cfg); err != nil {
		panic(err)
	}
	d, err := store.Load()
	if err != nil {
		panic(err)
	}

	loc := locale.Bundle{Lang: locale.Parse(cfg.Language)}

	a := fyneapp.NewWithID("desctop-otp")
	applyAppIcon(a)
	w := a.NewWindow(loc.T("window_title"))
	w.Resize(fyne.NewSize(560, 520))

	hint := widget.NewLabel("")
	hint.Wrapping = fyne.TextWrapWord

	sectionNew := widget.NewLabel("")
	sectionImport := widget.NewLabel("")

	nameEnt := widget.NewEntry()
	secEnt := widget.NewEntry()
	secEnt.Password = true
	importEnt := widget.NewEntry()

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

			copyBtn := widget.NewButton(loc.T("btn_copy"), func() {
				c, err := totp.GenerateCode(secret, time.Now())
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				w.Clipboard().SetContent(c)
			})

			delBtn := widget.NewButton(loc.T("btn_delete"), func() {
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

	addBtn := widget.NewButton(loc.T("btn_add"), func() {
		name := strings.TrimSpace(nameEnt.Text)
		sec := strings.TrimSpace(secEnt.Text)
		if name == "" {
			dialog.ShowInformation(loc.T("dlg_empty"), loc.T("dlg_need_name"), w)
			return
		}
		if err := store.ValidateSecret(sec); err != nil {
			showSecretErr(w, loc, err)
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

	importBtn := widget.NewButton(loc.T("btn_import"), func() {
		raw := strings.TrimSpace(importEnt.Text)
		if raw == "" {
			dialog.ShowInformation(loc.T("dlg_empty"), loc.T("dlg_need_uri"), w)
			return
		}
		n, sec, err := otpauth.ParseURI(raw)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if err := store.ValidateSecret(sec); err != nil {
			showSecretErr(w, loc, err)
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

	pasteBtn := widget.NewButton(loc.T("btn_paste"), func() {
		raw := w.Clipboard().Content()
		if raw == "" {
			dialog.ShowInformation(loc.T("dlg_clipboard_title"), loc.T("dlg_clipboard_empty"), w)
			return
		}
		importEnt.SetText(raw)
		n, sec, err := otpauth.ParseURI(raw)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if err := store.ValidateSecret(sec); err != nil {
			showSecretErr(w, loc, err)
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
		sectionNew,
		nameEnt,
		secEnt,
		addBtn,
		widget.NewSeparator(),
		sectionImport,
		importEnt,
		container.NewHBox(importBtn, pasteBtn),
	)

	lblDataDir := widget.NewLabel("")
	dataDirEnt := widget.NewEntry()
	dataDirEnt.SetText(cfg.DataDir)

	trayCheck := widget.NewCheck("", func(bool) {})
	trayCheck.SetChecked(cfg.MinimizeToTrayBool())

	lblHotkey := widget.NewLabel("")
	modAlt := widget.NewCheck("", nil)
	modCtrl := widget.NewCheck("", nil)
	modShift := widget.NewCheck("", nil)
	modWin := widget.NewCheck("", nil)
	modAlt.SetChecked(cfg.Modifiers&config.ModAlt != 0)
	modCtrl.SetChecked(cfg.Modifiers&config.ModControl != 0)
	modShift.SetChecked(cfg.Modifiers&config.ModShift != 0)
	modWin.SetChecked(cfg.Modifiers&config.ModWin != 0)

	keyEnt := widget.NewEntry()
	switch vk := cfg.VK; {
	case (vk >= 'A' && vk <= 'Z') || (vk >= '0' && vk <= '9'):
		keyEnt.SetText(string(rune(vk)))
	case vk == config.VKOEM3:
		keyEnt.SetText("`")
	}

	browseBtn := widget.NewButton("", func() {})
	lblLang := widget.NewLabel("")
	langSelect := widget.NewSelect([]string{"Русский", "English"}, nil)

	saveSettingsBtn := widget.NewButton("", func() {})

	var rebuildTray func()
	var refreshHint func()
	var applyLocale func()
	var restartHotkey func()
	var applyCloseBehavior func()
	var stopHK func()
	pressed := make(chan struct{}, 8)

	refreshHint = func() {
		ap, err := store.AccountsPath()
		if err != nil {
			ap = "?"
		}
		hint.SetText(loc.Format("hint_hotkey", cfg.HotkeyLabel(), ap))
	}

	rebuildTray = func() {
		if desk, ok := a.(desktop.App); ok {
			desk.SetSystemTrayMenu(fyne.NewMenu(loc.T("tray_menu"),
				fyne.NewMenuItem(loc.T("tray_show"), func() {
					w.Show()
					w.RequestFocus()
				}),
				fyne.NewMenuItemSeparator(),
				fyne.NewMenuItem(loc.T("tray_exit"), func() { a.Quit() }),
			))
		}
	}

	restartHotkey = func() {
		if stopHK != nil {
			stopHK()
			stopHK = nil
		}
		var hkErr error
		stopHK, hkErr = hotkey.Start(cfg.Modifiers, cfg.VK, pressed)
		if hkErr != nil {
			dialog.ShowError(hkErr, w)
		}
	}

	applyCloseBehavior = func() {
		if cfg.MinimizeToTrayBool() {
			w.SetCloseIntercept(func() {
				w.Hide()
			})
		} else {
			w.SetCloseIntercept(func() {
				a.Quit()
			})
		}
	}

	saveSettingsBtn.OnTapped = func() {
		vk, ok := config.ParseKeyString(keyEnt.Text)
		if !ok {
			dialog.ShowInformation(loc.T("dlg_empty"), loc.T("settings_err_hotkey"), w)
			return
		}
		mods := config.ModifiersFromBools(modAlt.Checked, modCtrl.Checked, modShift.Checked, modWin.Checked)
		if mods == 0 {
			dialog.ShowInformation(loc.T("dlg_empty"), loc.T("settings_err_hotkey"), w)
			return
		}
		cfg.Modifiers = mods
		cfg.VK = vk
		cfg.DataDir = strings.TrimSpace(dataDirEnt.Text)
		cfg.SetMinimizeToTray(trayCheck.Checked)

		if err := config.Save(cfg); err != nil {
			dialog.ShowError(errors.New(loc.T("settings_err_save")), w)
			return
		}
		if err := store.InitFromConfig(cfg); err != nil {
			dialog.ShowError(errors.New(loc.T("settings_err_data_dir")), w)
			return
		}
		nd, err := store.Load()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		*d = *nd
		rebuildList()
		restartHotkey()
		applyCloseBehavior()
		refreshHint()
		dialog.ShowInformation(loc.T("dlg_empty"), loc.T("settings_saved"), w)
	}

	browseBtn.OnTapped = func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			if p := uri.Path(); p != "" {
				dataDirEnt.SetText(p)
			}
		}, w)
	}

	applyLocale = func() {
		w.SetTitle(loc.T("window_title"))
		sectionNew.SetText(loc.T("section_new"))
		sectionImport.SetText(loc.T("section_import"))
		nameEnt.SetPlaceHolder(loc.T("ph_name"))
		secEnt.SetPlaceHolder(loc.T("ph_secret"))
		importEnt.SetPlaceHolder(loc.T("ph_uri"))
		addBtn.SetText(loc.T("btn_add"))
		importBtn.SetText(loc.T("btn_import"))
		pasteBtn.SetText(loc.T("btn_paste"))
		lblDataDir.SetText(loc.T("settings_data_dir"))
		browseBtn.SetText(loc.T("settings_browse"))
		trayCheck.SetText(loc.T("settings_minimize_tray"))
		lblHotkey.SetText(loc.T("settings_hotkey"))
		modAlt.SetText(loc.T("settings_mod_alt"))
		modCtrl.SetText(loc.T("settings_mod_ctrl"))
		modShift.SetText(loc.T("settings_mod_shift"))
		modWin.SetText(loc.T("settings_mod_win"))
		keyEnt.SetPlaceHolder(loc.T("settings_key"))
		lblLang.SetText(loc.T("settings_lang"))
		saveSettingsBtn.SetText(loc.T("settings_save"))
		if tabs := w.Content(); tabs != nil {
			if at, ok := tabs.(*container.AppTabs); ok && len(at.Items) > 1 {
				at.Items[0].Text = loc.T("tab_totp")
				at.Items[1].Text = loc.T("tab_settings")
				at.Refresh()
			}
		}
		refreshHint()
		rebuildTray()
		rebuildList()
	}

	settingsForm := container.NewVBox(
		lblDataDir,
		dataDirEnt,
		browseBtn,
		trayCheck,
		widget.NewSeparator(),
		lblHotkey,
		container.NewHBox(modAlt, modCtrl, modShift, modWin),
		keyEnt,
		widget.NewSeparator(),
		lblLang,
		langSelect,
		saveSettingsBtn,
	)

	scroll := container.NewVScroll(listBox)
	scroll.SetMinSize(fyne.NewSize(520, 220))

	totpContent := container.NewBorder(
		hint, form, nil, nil,
		scroll,
	)

	settingsContent := container.NewVScroll(settingsForm)
	settingsContent.SetMinSize(fyne.NewSize(520, 360))

	tabs := container.NewAppTabs(
		container.NewTabItem(loc.T("tab_totp"), totpContent),
		container.NewTabItem(loc.T("tab_settings"), settingsContent),
	)
	w.SetContent(tabs)

	langSelect.OnChanged = func(s string) {
		if s == "Русский" {
			cfg.Language = "ru"
		} else {
			cfg.Language = "en"
		}
		loc.Lang = locale.Parse(cfg.Language)
		if err := config.Save(cfg); err != nil {
			dialog.ShowError(err, w)
			return
		}
		applyLocale()
	}
	if locale.Parse(cfg.Language) == locale.EN {
		langSelect.SetSelected("English")
	} else {
		langSelect.SetSelected("Русский")
	}

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
						codeLabels[i].SetText(loc.T("code_error"))
						continue
					}
					codeLabels[i].SetText(c)
				}
			})
		}
	}()

	if desk, ok := a.(desktop.App); ok {
		_ = desk
		rebuildTray()
		applyCloseBehavior()
	}

	restartHotkey()
	defer func() {
		if stopHK != nil {
			stopHK()
		}
	}()

	go func() {
		for range pressed {
			fyne.Do(func() {
				w.Show()
				w.RequestFocus()
			})
		}
	}()

	applyLocale()

	if len(d.Accounts) > 0 {
		w.Hide()
	} else {
		w.Show()
	}
	a.Run()
}

func showSecretErr(w fyne.Window, loc locale.Bundle, err error) {
	switch {
	case errors.Is(err, store.ErrEmptySecret):
		dialog.ShowInformation(loc.T("dlg_empty"), loc.T("err_empty_secret"), w)
	case errors.Is(err, store.ErrInvalidSecret):
		dialog.ShowInformation(loc.T("dlg_empty"), loc.T("err_invalid_secret"), w)
	default:
		dialog.ShowError(err, w)
	}
}
