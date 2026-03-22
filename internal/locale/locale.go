package locale

import (
	"fmt"
	"strings"
)

// Lang is a BCP-47 style tag; only ru and en are supported.
type Lang string

const (
	RU Lang = "ru"
	EN Lang = "en"
)

// Parse normalizes a language string from config.
func Parse(s string) Lang {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "en", "english":
		return EN
	default:
		return RU
	}
}

// Bundle holds the active language for lookups.
type Bundle struct {
	Lang Lang
}

func (b Bundle) T(key string) string {
	m := table[key]
	if m == nil {
		return key
	}
	s := m[b.Lang]
	if s == "" {
		s = m[RU]
	}
	return s
}

// Format is fmt.Sprintf on T(key).
func (b Bundle) Format(key string, args ...any) string {
	return fmt.Sprintf(b.T(key), args...)
}

var table = map[string]map[Lang]string{
	"window_title": {
		RU: "TOTP",
		EN: "TOTP",
	},
	"tab_totp": {
		RU: "Коды",
		EN: "Codes",
	},
	"tab_settings": {
		RU: "Настройки",
		EN: "Settings",
	},
	"hint_hotkey": {
		RU: "Горячая клавиша: %s — показать окно. Данные: %s",
		EN: "Hotkey %s — show window. Data file: %s",
	},
	"section_new": {
		RU: "Новый аккаунт",
		EN: "New account",
	},
	"ph_name": {
		RU: "Название (например GitHub)",
		EN: "Name (e.g. GitHub)",
	},
	"ph_secret": {
		RU: "Секрет Base32",
		EN: "Base32 secret",
	},
	"ph_uri": {
		RU: "otpauth://totp/...",
		EN: "otpauth://totp/...",
	},
	"btn_add": {
		RU: "Добавить",
		EN: "Add",
	},
	"section_import": {
		RU: "Импорт из Google Authenticator и др.",
		EN: "Import (Google Authenticator, etc.)",
	},
	"btn_import": {
		RU: "Импорт URI",
		EN: "Import URI",
	},
	"btn_paste": {
		RU: "Из буфера",
		EN: "From clipboard",
	},
	"btn_copy": {
		RU: "Копировать",
		EN: "Copy",
	},
	"btn_delete": {
		RU: "Удалить",
		EN: "Delete",
	},
	"dlg_empty": {
		RU: "Пусто",
		EN: "Empty",
	},
	"dlg_need_name": {
		RU: "Укажите название.",
		EN: "Enter a name.",
	},
	"dlg_need_uri": {
		RU: "Вставьте otpauth:// ссылку.",
		EN: "Paste an otpauth:// URI.",
	},
	"dlg_clipboard_empty": {
		RU: "Буфер обмена пуст.",
		EN: "Clipboard is empty.",
	},
	"dlg_clipboard_title": {
		RU: "Буфер",
		EN: "Clipboard",
	},
	"code_error": {
		RU: "ошибка",
		EN: "error",
	},
	"tray_menu": {
		RU: "TOTP",
		EN: "TOTP",
	},
	"tray_show": {
		RU: "Показать",
		EN: "Show",
	},
	"tray_exit": {
		RU: "Выход",
		EN: "Quit",
	},
	"settings_data_dir": {
		RU: "Папка с accounts.json (пусто = рядом с exe)",
		EN: "Folder for accounts.json (empty = next to exe)",
	},
	"settings_browse": {
		RU: "Обзор…",
		EN: "Browse…",
	},
	"settings_minimize_tray": {
		RU: "Сворачивать в трей при закрытии окна",
		EN: "Minimize to tray when closing the window",
	},
	"settings_hotkey": {
		RU: "Глобальная горячая клавиша",
		EN: "Global hotkey",
	},
	"settings_mod_alt": {
		RU: "Alt",
		EN: "Alt",
	},
	"settings_mod_ctrl": {
		RU: "Ctrl",
		EN: "Ctrl",
	},
	"settings_mod_shift": {
		RU: "Shift",
		EN: "Shift",
	},
	"settings_mod_win": {
		RU: "Win",
		EN: "Win",
	},
	"settings_key": {
		RU: "Клавиша (A–Z или 0–9)",
		EN: "Key (A–Z or 0–9)",
	},
	"settings_lang": {
		RU: "Язык интерфейса",
		EN: "Interface language",
	},
	"settings_save": {
		RU: "Сохранить настройки",
		EN: "Save settings",
	},
	"settings_saved": {
		RU: "Настройки сохранены.",
		EN: "Settings saved.",
	},
	"settings_err_hotkey": {
		RU: "Нужна одна буква A–Z или цифра 0–9 и хотя бы один модификатор.",
		EN: "Use one letter A–Z or digit 0–9 and at least one modifier.",
	},
	"settings_err_save": {
		RU: "Не удалось сохранить config.json",
		EN: "Could not save config.json",
	},
	"settings_err_data_dir": {
		RU: "Некорректная папка данных",
		EN: "Invalid data folder",
	},
	"lang_ru": {
		RU: "Русский",
		EN: "Russian",
	},
	"lang_en": {
		RU: "English",
		EN: "English",
	},
	"err_empty_secret": {
		RU: "Пустой секрет",
		EN: "Empty secret",
	},
	"err_invalid_secret": {
		RU: "Некорректный Base32 или секрет",
		EN: "Invalid Base32 or secret",
	},
}
