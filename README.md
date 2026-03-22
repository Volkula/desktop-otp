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
- **About** tab: version, git branch, and build date (when set at compile time)

### Requirements (local build)

- [Go](https://go.dev/dl/) 1.22+
- C compiler for **CGO** (Fyne UI), e.g. [MinGW-w64](https://www.mingw-w64.org/) (`gcc` on `PATH`)

### Build (Windows)

```powershell
go mod tidy
go build -ldflags="-s -w -H windowsgui" -o otp.exe ./cmd/otp
```

`-H windowsgui` hides the console window.

The icon shown in **File Explorer** for `otp.exe` comes from a Windows resource (`cmd/otp/rsrc.syso`), not from Fyne. It is generated from `internal/app/icon.go` together with `build/windows/app.ico`. After changing the embedded PNG, run `go install github.com/akavel/rsrc@latest`, then `go generate ./cmd/otp` from the repo root (or run the `go generate` lines in `cmd/otp/main.go` manually).

To fill the **About** tab (optional), pass link-time strings, for example:

```powershell
$v = "dev"; $d = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ"); $b = (git rev-parse --abbrev-ref HEAD)
go build -ldflags "-s -w -H windowsgui -X desctop-otp/internal/buildinfo.Version=$v -X desctop-otp/internal/buildinfo.BuildDate=$d -X desctop-otp/internal/buildinfo.Branch=$b" -o otp.exe ./cmd/otp
```

### Build (Linux)

Fyne needs **CGO** and system headers/libs (GTK, X11, OpenGL, GLFW). On Debian/Ubuntu:

```bash
sudo apt-get update
sudo apt-get install -y build-essential pkg-config \
  xorg-dev libgl1-mesa-dev libglfw3-dev libgtk-3-dev
sudo ldconfig
```

Then:

```bash
export CGO_ENABLED=1
go mod tidy
go build -trimpath -ldflags="-s -w" -o desktop-otp ./cmd/otp
```

Optional **About** metadata:

```bash
VERSION="$(git describe --tags --always 2>/dev/null || echo dev)"
BUILD_DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
BRANCH="$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)"
go build -trimpath -ldflags "-s -w -X desctop-otp/internal/buildinfo.Version=${VERSION} -X desctop-otp/internal/buildinfo.BuildDate=${BUILD_DATE} -X desctop-otp/internal/buildinfo.Branch=${BRANCH}" -o desktop-otp ./cmd/otp
```

Run `./desktop-otp` from a folder where you want `accounts.json` (or set the data directory in Settings). Prebuilt `.deb` / `.rpm` packages may be published with project releases when available.

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

### Antivirus / SmartScreen (Windows)

Unsigned `otp.exe` builds are sometimes flagged as suspicious by **Microsoft Defender** or other antivirus software. The program uses a **global hotkey**, **system tray**, and (if enabled) a **Run** registry entry for autostart — behavior that resembles persistence tools, so **heuristic detection** can produce **false positives** for this open-source app.

**What you can do:** restore the file from **Protection history** and choose **Allow** on device; add a **folder exclusion** for the directory where you keep the app (Windows Security → Virus & threat protection → Manage settings → Exclusions). You can **report a false positive** via [Microsoft’s file submission](https://www.microsoft.com/wdsi/filesubmission). Long term, an **Authenticode code-signing certificate** (paid) reduces SmartScreen warnings and improves reputation.

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
- Вкладка **О программе**: версия, ветка git и дата сборки (если заданы при компиляции)

### Требования (сборка у себя)

- [Go](https://go.dev/dl/) 1.22+
- Компилятор C для **CGO** (интерфейс Fyne), например [MinGW-w64](https://www.mingw-w64.org/) (`gcc` в `PATH`)

### Сборка (Windows)

```powershell
go mod tidy
go build -ldflags="-s -w -H windowsgui" -o otp.exe ./cmd/otp
```

Флаг `-H windowsgui` отключает консольное окно.

Иконка **в проводнике** у `otp.exe` берётся из ресурса Windows (`cmd/otp/rsrc.syso`), а не из Fyne. Файлы `rsrc.syso` и `build/windows/app.ico` собираются из `internal/app/icon.go`. После смены PNG в `icon.go` установите `github.com/akavel/rsrc`, затем выполните `go generate ./cmd/otp` из корня репозитория.

Чтобы во вкладке **О программе** отображались версия, ветка и дата (по желанию), задайте `-X` при сборке:

```powershell
$v = "dev"; $d = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ"); $b = (git rev-parse --abbrev-ref HEAD)
go build -ldflags "-s -w -H windowsgui -X desctop-otp/internal/buildinfo.Version=$v -X desctop-otp/internal/buildinfo.BuildDate=$d -X desctop-otp/internal/buildinfo.Branch=$b" -o otp.exe ./cmd/otp
```

### Сборка (Linux)

Для Fyne нужны **CGO** и системные библиотеки (GTK, X11, OpenGL, GLFW). В Debian/Ubuntu:

```bash
sudo apt-get update
sudo apt-get install -y build-essential pkg-config \
  xorg-dev libgl1-mesa-dev libglfw3-dev libgtk-3-dev
sudo ldconfig
```

Далее:

```bash
export CGO_ENABLED=1
go mod tidy
go build -trimpath -ldflags="-s -w" -o desktop-otp ./cmd/otp
```

Метаданные для вкладки **О программе** (по желанию):

```bash
VERSION="$(git describe --tags --always 2>/dev/null || echo dev)"
BUILD_DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
BRANCH="$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)"
go build -trimpath -ldflags "-s -w -X desctop-otp/internal/buildinfo.Version=${VERSION} -X desctop-otp/internal/buildinfo.BuildDate=${BUILD_DATE} -X desctop-otp/internal/buildinfo.Branch=${BRANCH}" -o desktop-otp ./cmd/otp
```

Запускайте `./desktop-otp` из каталога, где должен лежать `accounts.json` (или укажите папку данных в настройках). Готовые **.deb** / **.rpm** могут быть в релизах на GitHub при публикации.

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

### Антивирус и SmartScreen (Windows)

Неподписанный **`otp.exe`** иногда попадает под эвристику **Защитника Windows** или другого антивируса. У программы есть **глобальная горячая клавиша**, **трей** и при желании запись в **автозагрузке** (реестр) — такое сочетание похоже на «персистентность», поэтому возможны **ложные срабатывания**.

**Что сделать:** откройте **Журнал защиты** → восстановите файл и нажмите **Разрешить на устройстве**; при необходимости добавьте **исключение по папке**, где лежит программа (Параметры → Конфиденциальность и защита → Безопасность Windows → Защита от вирусов и угроз → Управление настройками → Исключения). Ложное срабатывание можно [отправить в Microsoft](https://www.microsoft.com/wdsi/filesubmission). На будущее снижает предупреждения **подпись кода** (Authenticode, сертификат платный).

### Сборка в CI

После push в репозиторий [GitHub Actions](.github/workflows/build.yml) собирает Windows-бинарник `otp.exe` и прикладывает его как артефакт workflow.
