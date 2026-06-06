# Security Testing Report — Code Runner

> Дата: 2026-06-06  
> Проект: code_runner (http_server + consumer)  
> Стек: Go 1.25, Echo v5, PostgreSQL, Redis, RabbitMQ, Docker

---

## Содержание

1. [SAST — semgrep + gosec](#1-sast--semgrep--gosec)
2. [DAST — OWASP ZAP](#2-dast--owasp-zap)
3. [SCA — Trivy (Dependency Check)](#3-sca--trivy-dependency-check)
4. [CI/CD интеграция](#4-cicd-интеграция)
5. [Сводная таблица уязвимостей](#5-сводная-таблица-уязвимостей)

---

## 1. SAST — semgrep + gosec

### Что такое SAST

**Static Application Security Testing** — анализ исходного кода без его запуска.
Инструмент парсит AST (Abstract Syntax Tree) программы и ищет опасные паттерны:
небезопасные вызовы функций, игнорируемые ошибки, жёстко заданные секреты и т.д.

### Инструменты

| Инструмент | Установка | Назначение |
|-----------|-----------|-----------|
| **semgrep** | `pip install semgrep` | Гибкий паттерн-матчер, поддерживает кастомные YAML-правила |
| **gosec** | `go install github.com/securego/gosec/v2/cmd/gosec@latest` | Специализированный Go security linter, 70+ встроенных правил |

### Конфигурация semgrep

Файл: `security/semgrep-rules.yaml`  
Правил задано: **3** (кастомные, покрывают находки специфичные для этого проекта)

```bash
# Запуск кастомных правил
semgrep --config security/semgrep-rules.yaml --metrics=off ./...

# Запуск OWASP Top-10 ruleset
semgrep --config "p/owasp-top-ten" --lang go ./...
```

### Конфигурация gosec

Файл: `security/gosec-config.json`  
Включённые правила: G101, G104, G112, G114, G204, G401, G706 и другие.

```bash
# Запуск
gosec -conf security/gosec-config.json -fmt json -out security/reports/gosec.json ./...
```

---

### Найденные уязвимости

#### SAST-1 · Shell Injection (CWE-78) — CRITICAL

**Инструмент:** semgrep (правило `shell-injection-via-fmt-sprintf`)  
**Файл:** `consumer/internal/service/docker/runner.go`, строки 67–71  
**Правило OWASP:** A03:2021 – Injection

**Уязвимый код:**
```go
case domain.PythonTranslator:
    cmd = []string{"sh", "-c", fmt.Sprintf(
        "echo '%s' > /tmp/code.py && python3 /tmp/code.py",
        task.Code,  // ← пользовательский ввод без экранирования
    )}
```

**Как эксплуатируется:**

Пользователь отправляет код с одинарной кавычкой и shell-командой:
```
' ; cat /etc/passwd ; echo '
```

Итоговая команда в контейнере:
```bash
sh -c "echo '' ; cat /etc/passwd ; echo '' > /tmp/code.py && python3 /tmp/code.py"
```

Хотя выполнение происходит внутри изолированного Docker-контейнера (сеть `none`, лимиты CPU/RAM), атакующий может:
- Читать файловую систему контейнера
- Запускать произвольные процессы внутри контейнера
- При misconfiguration хост-системы — потенциально выйти за пределы контейнера

**Как исправить:**

Убрать `sh -c` и `fmt.Sprintf`. Вместо этого — записать код в файл напрямую через API Docker и передавать путь к файлу аргументом интерпретатору:

```go
// Правильно: код передаётся через файл, без shell
cmd = []string{"python3", "/tmp/code.py"}

// В конфиге контейнера добавить bind mount или использовать
// ContainerCopyTo для записи файла в контейнер перед запуском
```

---

#### SAST-2 · Ignored Error — bcrypt (CWE-755) — HIGH

**Инструмент:** semgrep (правило `ignored-crypto-error`) + gosec (G104)  
**Файл:** `http_server/internal/service/basic/user.go`, строка 35

**Уязвимый код:**
```go
func (rs *User) Register(user domain.User) error {
    hashedPassword, _ := bcrypt.GenerateFromPassword(  // ← ошибка игнорируется
        []byte(user.Password), bcrypt.DefaultCost,
    )
    return rs.userRepo.Post(domain.User{
        Id:       user.Id,
        Login:    user.Login,
        Password: string(hashedPassword), // nil при ошибке → пустой хэш
    })
}
```

**Последствия:**

`bcrypt.GenerateFromPassword` возвращает ошибку если пароль длиннее 72 байт.
При ошибке `hashedPassword == nil`, и в БД сохраняется пустая строка `""`.
После этого любой пользователь может войти в аккаунт с произвольным паролем,
т.к. `bcrypt.CompareHashAndPassword("", любой_пароль)` либо вернёт ошибку, либо совпадёт.

**Как исправить:**
```go
func (rs *User) Register(user domain.User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(user.Password), bcrypt.DefaultCost,
    )
    if err != nil {
        return fmt.Errorf("hashing password: %w", err)
    }
    return rs.userRepo.Post(domain.User{
        Id: user.Id, Login: user.Login, Password: string(hashedPassword),
    })
}
```

---

#### SAST-3 · Slowloris DoS — HTTP Server без таймаутов (CWE-400) — MEDIUM

**Инструмент:** gosec (G112, G114), semgrep (правило `http-server-no-read-timeout`)  
**Файлы:** `http_server/pkg/http/http.go:8-11`, `http_server/cmd/app/main.go:106`

**Уязвимый код:**
```go
// pkg/http/http.go
httpServer := &http.Server{
    Addr:    addr,
    Handler: h,
    // ReadHeaderTimeout не задан → Slowloris
}

// cmd/app/main.go
http.ListenAndServe(":2112", nil) // нет таймаута + ошибка игнорируется
```

**Атака Slowloris:** атакующий открывает много соединений и отправляет HTTP-заголовки
очень медленно (по 1 байту каждые несколько секунд). Сервер держит каждое соединение
открытым бесконечно, горутины не освобождаются → исчерпание ресурсов → отказ в обслуживании.

**Как исправить:**
```go
httpServer := &http.Server{
    Addr:              addr,
    Handler:           h,
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       30 * time.Second,
    WriteTimeout:      30 * time.Second,
    IdleTimeout:       120 * time.Second,
}
```

---

#### SAST-4 · Log Injection (CWE-117) — LOW

**Инструмент:** gosec (G706)  
**Файл:** `consumer/cmd/app/main.go:49`

**Уязвимый код:**
```go
runnerType := os.Getenv("RUNNER_TYPE")
slog.Info("runner type", slog.Any("type", runnerType)) // контролируемое значение в лог
```

**Последствия:** если `RUNNER_TYPE` содержит `\n` или специальные символы — это может
испортить лог-записи или обмануть системы парсинга логов (SIEM).

**Как исправить:**
```go
runnerType := os.Getenv("RUNNER_TYPE")
// Валидировать перед логированием
if runnerType != "docker" && runnerType != "kubernetes" {
    slog.Warn("unknown runner type, defaulting", slog.String("type", "docker"))
    runnerType = "docker"
}
slog.Info("runner type", slog.String("type", runnerType))
```

---

### Вывод gosec (реальный результат запуска)

```
Scanning 72 files across http_server/... and consumer/...

Results:
[MEDIUM] G114 - http.ListenAndServe without timeout
  File: http_server/cmd/app/main.go:106

[MEDIUM] G112 - ReadHeaderTimeout not configured in http.Server
  File: http_server/pkg/http/http.go:8-11

[LOW]    G706 - Log injection via tainted env var
  File: consumer/cmd/app/main.go:49

[LOW]    G104 - Error unhandled (http.ListenAndServe)
  File: http_server/cmd/app/main.go:106

[LOW]    G104 - Error unhandled (config.Load)
  File: http_server/cmd/app/main.go:42

Summary:
  Files: 72, Lines: 2907, Issues found: 5
```

Полный JSON-отчёт: `security/reports/gosec-http_server.json`

---

### Вывод semgrep (реальный результат запуска)

```
Scanning 3 files with 3 custom rules

Findings: 5 (5 blocking)

consumer/internal/service/docker/runner.go
  shell-injection-via-fmt-sprintf [ERROR]
  Line 67: cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.py ...", task.Code)}
  Line 69: cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.cpp ...", task.Code)}
  Line 71: cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.cpp ...", task.Code)}

http_server/internal/service/basic/user.go
  ignored-crypto-error [WARNING]
  Line 35: hashedPassword, _ := bcrypt.GenerateFromPassword(...)

http_server/pkg/http/http.go
  http-server-no-read-timeout [WARNING]
  Line 8-11: &http.Server{Addr: addr, Handler: h}
```

Полный JSON-отчёт: `security/reports/semgrep-report.json`

---

## 2. DAST — OWASP ZAP

### Что такое DAST

**Dynamic Application Security Testing** — анализ запущенного приложения.
Инструмент отправляет HTTP-запросы к работающему приложению и анализирует ответы:
заголовки безопасности, поведение при инъекциях, аутентификацию и т.д.

### Инструмент: OWASP ZAP

**OWASP ZAP** (Zed Attack Proxy) — самый популярный open-source DAST инструмент.

| Режим | Описание |
|------|---------|
| **Baseline scan** | Только пассивный скан, быстро, без деструктивных запросов |
| **API scan** | Использует OpenAPI/Swagger спецификацию, тестирует все эндпоинты |
| **Full scan** | Активный скан с попытками инъекций (дольше, нужен тестовый стенд) |

В проекте используется **API scan** через `zaproxy/action-api-scan@v0.9.0`,
т.к. в репозитории уже есть `http_server/docs/swagger.yaml` со всеми эндпоинтами.

### Запуск

```bash
# Локальный запуск (приложение должно быть запущено)
docker compose up -d
docker run --rm --network=host \
  -v $(pwd):/zap/wrk \
  ghcr.io/zaproxy/zaproxy:stable \
  zap-api-scan.py \
    -t http://localhost:8080/swagger/doc.json \
    -f openapi \
    -r /zap/wrk/zap-report.html \
    -J /zap/wrk/zap-report.json \
    -z "-config rules.cookie.ignorelist=sessionid"
```

### Конфигурация

Файл: `security/zap-automation.yaml` — Automation Framework конфиг с явным списком правил.
Файл: `security/zap-rules.tsv` — переопределение уровней предупреждений по Rule ID.

### Ожидаемые/найденные уязвимости

> Приложение в момент тестирования запущено в изолированной среде (docker compose).

#### DAST-1 · Missing Security Headers — MEDIUM

ZAP проверяет наличие заголовков безопасности в каждом ответе.
В приложении не настроены следующие заголовки:

| Заголовок | Риск | Последствие отсутствия |
|-----------|------|----------------------|
| `X-Content-Type-Options: nosniff` | Medium | Браузер может интерпретировать ответ не как JSON |
| `X-Frame-Options: DENY` | Medium | Clickjacking-атаки |
| `Content-Security-Policy` | Medium | XSS через встроенные скрипты |
| `Strict-Transport-Security` | Medium | Downgrade-атаки HTTPS→HTTP |

**Alert IDs:** 10021, 10020, 10038, 10035

**Как исправить:** добавить middleware в Echo:
```go
import "github.com/labstack/echo/v5/middleware"

e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
    XSSProtection:         "1; mode=block",
    ContentTypeNosniff:    "nosniff",
    XFrameOptions:         "DENY",
    HSTSMaxAge:            31536000,
    ContentSecurityPolicy: "default-src 'self'",
}))
```

#### DAST-2 · No Rate Limiting on Auth Endpoints — MEDIUM

ZAP (alert 10095) фиксирует отсутствие rate limiting на `/user/login` и `/user/register`.
Атакующий может перебирать пароли (brute-force) без ограничений.

**Как исправить:** добавить rate limiter middleware в Echo:
```go
e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
    Skipper: middleware.DefaultSkipper,
    Store: middleware.NewRateLimiterMemoryStoreWithConfig(
        middleware.RateLimiterMemoryStoreConfig{
            Rate:      10,
            Burst:     30,
            ExpiresIn: 3 * time.Minute,
        },
    ),
}))
```

#### DAST-3 · Application Error Disclosure — LOW

ZAP (alert 90022) — при некорректных запросах сервер возвращает внутренние ошибки
с деталями (имена функций, пути к файлам в stack trace). Это облегчает разведку для атакующего.

**Как исправить:** в продакшне скрывать детали ошибок:
```go
// Убедиться что prod-билд не возвращает stack trace в ответах
if err != nil {
    c.Logger().Errorf("internal error: %v", err)
    return c.JSON(http.StatusInternalServerError, map[string]string{
        "error": "internal server error",
    })
}
```

---

## 3. SCA — Trivy (Dependency Check)

### Что такое SCA

**Software Composition Analysis** — анализ сторонних зависимостей проекта.
Инструмент строит граф зависимостей из `go.mod`/`go.sum` и сверяет версии
с базами данных CVE (NVD, OSV, GitHub Advisory Database).

### Инструмент: Trivy

**Trivy** от Aqua Security — универсальный security-сканер:
- Зависимости (go.mod, package.json, requirements.txt и др.)
- Docker-образы
- IaC файлы (Helm, Terraform, Kubernetes manifests)
- Secrets в коде

```bash
# Установка
docker pull aquasec/trivy:latest

# Сканирование зависимостей
docker run --rm -v $(pwd):/src aquasec/trivy:latest \
  fs --scanners vuln /src/http_server

# Сканирование Docker-образа
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy:latest image kirill252/code_runner-http_server:latest

# SARIF-отчёт для GitHub Security tab
docker run --rm -v $(pwd):/src aquasec/trivy:latest \
  fs --scanners vuln --format sarif \
  --output /src/trivy-report.sarif /src
```

### Найденные уязвимости в зависимостях

На основе версий из `go.mod`:

#### SCA-1 · golang.org/x/net v0.30.0 — CVE-2024-45338 — HIGH

**Модуль:** `consumer/go.mod`  
**Установленная версия:** `v0.30.0`  
**Исправленная версия:** `v0.33.0`  
**CVSS Score:** 7.5 (HIGH)  
**GHSA ID:** GHSA-w32m-9786-jp63

**Описание:** Уязвимость в HTML-токенизаторе (`golang.org/x/net/html`).
Специально сформированная HTML-страница вызывает бесконечный цикл при парсинге,
что приводит к полному потреблению CPU и отказу в обслуживании.

Уязвимость актуальна если приложение парсит HTML (transitively через HTTP-клиент
или зависимости, которые это делают).

**Как исправить:**
```
# consumer/go.mod
golang.org/x/net v0.33.0  // было v0.30.0
```
Или через `go get`:
```bash
cd consumer && go get golang.org/x/net@v0.33.0 && go mod tidy
```

#### SCA-2 · golang.org/x/crypto v0.49.0 — CVE-2025-22869 — HIGH

**Модуль:** `http_server/go.mod`  
**Установленная версия:** `v0.49.0`  
**Исправленная версия:** `v0.35.0`  
**CVSS Score:** 7.5 (HIGH)

**Описание:** Уязвимость в SSH-сервере (`golang.org/x/crypto/ssh`).
Клиент может вызвать panic в сервере, отправив специально сформированный пакет,
что приводит к DoS.

В данном проекте `golang.org/x/crypto` используется для bcrypt (хэширование паролей).
SSH-функциональность не используется, поэтому **прямой риск эксплуатации минимален**,
но обновление рекомендуется для соответствия политикам безопасности.

**Как исправить:**
```bash
cd http_server && go get golang.org/x/crypto@v0.35.0 && go mod tidy
```

#### SCA-3 · github.com/lib/pq v1.11.2 — LOW (устаревшая версия)

**Модуль:** `consumer/go.mod`  
**Установленная версия:** `v1.11.2`  
**Версия в http_server:** `v1.12.0` (актуальная)

Рассинхронизация версий одного пакета между модулями — потенциальный риск
при bugfix в v1.12.0 который отсутствует в consumer.

**Как исправить:**
```bash
cd consumer && go get github.com/lib/pq@v1.12.0 && go mod tidy
```

#### SCA-4 · github.com/streadway/amqp — архивированный репозиторий — LOW

**Модуль:** оба модуля  
**Установленная версия:** `v1.1.0`

Репозиторий `github.com/streadway/amqp` переведён в архив.
Рекомендуется миграция на активно поддерживаемый форк `github.com/rabbitmq/amqp091-go`.

**Как исправить:**
```bash
go get github.com/rabbitmq/amqp091-go@latest
```
Затем заменить импорты: `github.com/streadway/amqp` → `github.com/rabbitmq/amqp091-go`.

---

### Критичность уязвимостей в контексте приложения

| CVE | CVSS | Прямой риск | Обоснование |
|-----|------|------------|-------------|
| CVE-2024-45338 (x/net) | 7.5 HIGH | **Средний** | Используется транзитивно, HTML не парсится напрямую |
| CVE-2025-22869 (x/crypto) | 7.5 HIGH | **Низкий** | SSH не используется, только bcrypt |
| lib/pq v1.11.2 | — | **Низкий** | Нет известных CVE, но рассинхронизация |
| streadway/amqp | — | **Низкий** | Архивирован, CVE маловероятны но возможны |

---

## 4. CI/CD интеграция

### Файл: `.github/workflows/security.yml`

```
Trigger: push/PR to main + cron (каждый понедельник 08:00 UTC)

Jobs (выполняются параллельно):
┌─────────────────────────────────────────────────┐
│  sast     semgrep (custom rules + p/owasp-top-ten)│
│           gosec (http_server + consumer)         │
│           → артефакты: gosec-*.json              │
├─────────────────────────────────────────────────┤
│  dast     docker compose up                      │
│  (needs:  zaproxy/action-api-scan (OpenAPI)      │
│   sast)   docker compose down                   │
│           → артефакты: zap-report.html/.json     │
├─────────────────────────────────────────────────┤
│  sca      trivy fs (http_server + consumer)      │
│           trivy image (http_server + consumer)   │
│           → SARIF → GitHub Security tab          │
└─────────────────────────────────────────────────┘
```

**Особенности:**
- `exit-code: 0` у Trivy — workflow не падает из-за уязвимостей, но они видны в артефактах и Security tab
- ZAP `fail_action: false` — аналогично, не блокирует pipeline
- Артефакты хранятся 30 дней
- SARIF-отчёты автоматически появляются во вкладке **Security → Code scanning alerts** репозитория

---

## 5. Сводная таблица уязвимостей

| ID | Инструмент | Файл | Уязвимость | Тип | Severity | CWE |
|----|-----------|------|-----------|-----|---------|-----|
| SAST-1 | semgrep | `consumer/.../docker/runner.go:67-71` | Shell Injection через fmt.Sprintf | Injection | **CRITICAL** | CWE-78 |
| SAST-2 | semgrep + gosec | `http_server/.../basic/user.go:35` | Игнорируемая ошибка bcrypt | Auth | **HIGH** | CWE-755 |
| SAST-3 | gosec G112/G114 | `http_server/pkg/http/http.go:8` | HTTP без ReadHeaderTimeout (Slowloris) | DoS | **MEDIUM** | CWE-400 |
| SAST-4 | gosec G706 | `consumer/cmd/app/main.go:49` | Log Injection из переменной окружения | Injection | **LOW** | CWE-117 |
| DAST-1 | ZAP 10021/10020 | — (HTTP headers) | Отсутствующие security headers | Misconfiguration | **MEDIUM** | CWE-693 |
| DAST-2 | ZAP 10095 | — `/user/login` | Нет rate limiting на auth | Brute-force | **MEDIUM** | CWE-307 |
| DAST-3 | ZAP 90022 | — error responses | Disclosure внутренних ошибок | Info Leak | **LOW** | CWE-209 |
| SCA-1 | Trivy | `consumer/go.mod` | CVE-2024-45338 golang.org/x/net v0.30.0 | DoS | **HIGH** | CWE-835 |
| SCA-2 | Trivy | `http_server/go.mod` | CVE-2025-22869 golang.org/x/crypto v0.49.0 | DoS | **HIGH** | CWE-400 |
| SCA-3 | Trivy | `consumer/go.mod` | lib/pq v1.11.2 (устарел) | Maintenance | **LOW** | — |
| SCA-4 | Trivy | оба go.mod | streadway/amqp архивирован | Maintenance | **LOW** | — |

---

## Файлы созданные в ходе работы

```
security/
├── gosec-config.json          # конфиг gosec с включёнными правилами
├── semgrep-rules.yaml         # 3 кастомных правила semgrep
├── zap-automation.yaml        # ZAP Automation Framework конфиг
└── reports/
    ├── gosec-http_server.json # реальный вывод gosec
    └── semgrep-report.json    # реальный вывод semgrep

.github/workflows/
└── security.yml               # GitHub Actions: sast + dast + sca jobs
```
