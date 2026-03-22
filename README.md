# Desktop OTP

Portable Windows TOTP generator: pre-saved secrets, global hotkey, system tray.

---

## English

### Features

- TOTP codes (RFC 6238, default: SHA-1, 30 s, 6 digits) for saved accounts
- **Global hotkey** to show the window (default: **Ctrl+Shift+O**)
- **System tray**: hide to tray, show / exit from the menu
- Data in **`accounts.json`** next to the executable (easy to copy the whole folder)
- Add manually (name + Base32 secret) or import **`otpauth://`** links (including from clipboard)

### Requirements (local build)

- [Go](https://go.dev/dl/) 1.22+
- C compiler for **CGO** (Fyne UI), e.g. [MinGW-w64](https://www.mingw-w64.org/) (`gcc` on `PATH`)

### Build (Windows)

```powershell
go mod tidy
go build -ldflags="-s -w -H windowsgui" -o otp.exe ./cmd/otp
```

`-H windowsgui` hides the console window.

### Portable usage

Copy `otp.exe` to a folder. On first run (or with an empty account list) the window opens; after you add accounts, the app can start **minimized to the tray**. Secrets are stored in **`accounts.json`** in the same folder as the executable.

### Hotkey (`config.json`)

Optional file next to `otp.exe`:

```json
{
  "modifiers": 6,
  "vk": 79
}
```

| Field        | Meaning |
|-------------|---------|
| `modifiers` | Bit mask: Alt = 1, Ctrl = 2, Shift = 4, Win = 8 (e.g. `6` = Ctrl+Shift) |
| `vk`        | Virtual-key code (e.g. `79` = **O**) |

If the file is missing, defaults are used (Ctrl+Shift+O).

### Security note

`accounts.json` holds **shared secrets** in readable form. Protect the folder (encryption, permissions) if needed.

### CI builds

On push to the default branch, [GitHub Actions](.github/workflows/build.yml) builds a Windows `otp.exe` and uploads it as a workflow artifact.

---

## Русский

### Возможности

- Генерация **TOTP** (RFC 6238, по умолчанию SHA-1, 30 с, 6 цифр) по заранее сохранённым аккаунтам
- **Глобальная горячая клавиша** для показа окна (по умолчанию **Ctrl+Shift+O**)
- **Трей**: сворачивание в область уведомлений, пункты «Показать» / «Выход»
- Данные в **`accounts.json`** рядом с exe — папку с программой легко копировать на другой ПК
- Ручное добавление (имя + секрет Base32) или импорт ссылок **`otpauth://`** (в т.ч. из буфера обмена)

### Требования (сборка у себя)

- [Go](https://go.dev/dl/) 1.22+
- Компилятор C для **CGO** (интерфейс Fyne), например [MinGW-w64](https://www.mingw-w64.org/) (`gcc` в `PATH`)

### Сборка (Windows)

```powershell
go mod tidy
go build -ldflags="-s -w -H windowsgui" -o otp.exe ./cmd/otp
```

Флаг `-H windowsgui` отключает консольное окно.

### Переносимость

Положите `otp.exe` в любую папку. При первом запуске или пустом списке аккаунтов окно показывается сразу; после добавления аккаунтов приложение может стартовать **в трее**. Секреты хранятся в **`accounts.json`** рядом с исполняемым файлом.

### Горячая клавиша (`config.json`)

Необязательный файл рядом с `otp.exe`:

```json
{
  "modifiers": 6,
  "vk": 79
}
```

| Поле        | Значение |
|------------|----------|
| `modifiers` | Битовая маска: Alt = 1, Ctrl = 2, Shift = 4, Win = 8 (например `6` = Ctrl+Shift) |
| `vk`        | Код виртуальной клавиши (например `79` = **O**) |

Если файла нет — используются значения по умолчанию (Ctrl+Shift+O).

### Безопасность

В **`accounts.json`** лежат **секреты** в открытом виде. При необходимости ограничьте доступ к папке или шифруйте носитель.

### Сборка в CI

После push в репозиторий [GitHub Actions](.github/workflows/build.yml) собирает Windows-бинарник `otp.exe` и прикладывает его как артефакт workflow.
