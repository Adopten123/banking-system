# 🏦 Account Service (Banking Core)

[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-316192?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-Event%20Driven-FF6600?style=flat&logo=rabbitmq)](https://www.rabbitmq.com/)
[![gRPC](https://img.shields.io/badge/gRPC-Vault%20Integration-244c5a?style=flat&logo=grpc)](https://grpc.io/)

**Account Service** - это ядро банковской системы [github.com/Adopten123/banking-system](https://github.com/Adopten123/banking-system), отвечающее за управление счетами пользователей, картами, балансами и безопасным проведением финансовых транзакций.

Сервис спроектирован с учетом строгих требований к финансовым системам (ACID, консистентность данных). Под капотом реализована защита от двойных списаний (идемпотентность), оптимистичные блокировки, а также асинхронное событийно-ориентированное взаимодействие с другими узлами системы.

---

##  Фичи

* **Управление счетами (Account Management):**
    * Создание, закрытие, заморозка (Антифрод) и блокировка счетов.
    * Управление кредитными лимитами (поддержка овердрафта).
* **Транзакционное ядро (Ledger):**
    * Реализация принципа двойной записи (Double-entry bookkeeping).
    * Безопасные операции пополнения (Deposits), снятия (Withdrawals) и переводов (Transfers) между счетами и картами.
* **Идемпотентность (Idempotency):**
    * Все финансовые операции требуют заголовок `Idempotency-Key` для абсолютной защиты от дублирования платежей при сетевых сбоях (Retry/Timeout).
* **Card Management (Интеграция с Vault):**
    * Выпуск физических и виртуальных карт, привязанных к счетам.
    * Чувствительные данные (PAN, CVV, PIN) **не хранятся** в БД сервиса — используется gRPC-интеграция с защищенным Card Vault.
* **Регулярные платежи (Автоплатежи):**
    * Встроенный background-воркер для автоматического списания средств по CRON-расписанию.
    * Безопасное параллельное выполнение выборки платежей благодаря механизму `FOR UPDATE SKIP LOCKED`.
* **Мультивалютность и кросс-курсы:**
    * Интеграция с внешним Exchanger Service (HTTP) для конвертации валют «на лету» при переводах.
* **Event-Driven Architecture:**
    * Асинхронная публикация доменных событий (TransferCreated, AccountCreated и др.) в RabbitMQ для последующей обработки (нотификации, аналитика).

## 🛠 Стек технологий (Tech Stack)

Сервис построен на базе современных и производительных инструментов экосистемы Go:

* **Язык:** Go 1.26 (Golang)
* **База данных:** PostgreSQL 17
* **Работа с БД:** `pgx/v5` (driver & pool) + `sqlc` (Type-safe SQL генератор)
* **Роутинг HTTP:** `go-chi/chi` (легковесный и быстрый роутер)
* **Брокер сообщений:** RabbitMQ (публикация событий)
* **gRPC Client:** `google.golang.org/grpc` (интеграция с защищенным хранилищем карт)
* **Миграции БД:** `golang-migrate`
* **Математика денег:** `shopspring/decimal` (для точных вычислений без потери плавающей точки)
* **Тестирование:** `testcontainers-go` (PostgreSQL в Docker), `testify`, `mockgen`
* **Планировщик задач:** `robfig/cron/v3`

---

## База данных

База данных спроектирована по принципам строгого финансового учета с разделением сущностей для минимизации блокировок при конкурентных запросах.

### Основные архитектурные решения:
1. **Разделение Accounts и Balances (1:1):** Таблица `accounts` хранит метаданные счета и поле `version` для реализации механизма **Optimistic Locking**. Таблица `account_balances` выступает как агрегированный кэш балансов и хранит кредитные лимиты. Это позволяет обновлять баланс, не блокируя чтение метаданных счета.
2. **Двойная запись (Double-entry bookkeeping):**
   Таблица `transactions` отражает бизнес-суть перевода (кто, кому, сколько, ключ идемпотентности). Таблица `postings` содержит физические финансовые проводки (списание со счета А, зачисление на счет Б) с учетом кросс-курсов конвертации (`exchange_rate`).
3. **Безопасность карт (Cards):**
   Таблица `cards` содержит только маскированный номер (`pan_mask`) и токены. Сами карточные данные живут во внешнем контуре безопасности.

### ER-Диаграмма (Entity-Relationship Diagram)

```mermaid
erDiagram
    USERS ||--o{ ACCOUNTS : "has"
    ACCOUNTS ||--|| ACCOUNT_BALANCES : "1:1 relation"
    ACCOUNTS ||--o{ CARDS : "issues"
    ACCOUNTS ||--o{ POSTINGS : "receives/sends"
    TRANSACTIONS ||--|{ POSTINGS : "contains (double-entry)"
    RECURRING_PAYMENTS }o--|| ACCOUNTS : "targets (via UUID)"

    USERS {
        uuid id PK
        varchar username
    }
    ACCOUNTS {
        bigint id PK
        uuid public_id "For External API"
        uuid user_id FK
        varchar currency_code
        int version "Optimistic Lock"
    }
    ACCOUNT_BALANCES {
        bigint account_id PK, FK
        numeric balance "Cached sum of postings"
        numeric credit_limit
    }
    CARDS {
        uuid id PK "Token to Vault"
        bigint account_id FK
        varchar pan_mask
        varchar status
    }
    TRANSACTIONS {
        uuid id PK
        uuid source_id
        uuid destination_id
        varchar idempotency_key "Unique"
        int status_id
    }
    POSTINGS {
        bigint id PK
        uuid transaction_id FK
        bigint account_id FK
        numeric amount "Negative = debit, Positive = credit"
        numeric exchange_rate
    }
    RECURRING_PAYMENTS {
        uuid id PK
        uuid source_id
        numeric amount
        varchar cron_expression
        timestamptz next_execution_time
        boolean is_active
    }
```

## API & Интеграции

Сервис предоставляет REST API для клиентских приложений и общается с другими сервисами через gRPC и брокер сообщений.

### REST API

API спроектировано в стиле RESTful. Все транзакционные запросы (POST) требуют заголовок `Idempotency-Key`.

| Группа | Метод | Эндпоинт | Описание |
| :--- | :--- | :--- | :--- |
| **Accounts** | `POST` | `/api/accounts` | Создание нового счета (с нулевым балансом) |
| | `GET` | `/api/accounts/{id}` | Получение информации о счете |
| | `GET` | `/api/accounts/{id}/balance` | Получение текущего баланса |
| | `PUT` | `/api/accounts/{id}/credit-limit` | Установка кредитного лимита (овердрафт) |
| | `POST` | `/api/accounts/{id}/[status]` | Изменение статуса: `activate`, `block`, `freeze`, `close` |
| | `GET` | `/api/accounts/{id}/transactions` | История транзакций по счету (с пагинацией) |
| **Cards** | `POST` | `/api/accounts/{id}/card` | Выпуск новой карты к счету |
| | `GET` | `/api/accounts/{id}/cards` | Список привязанных карт (маскированный PAN) |
| | `GET` | `/api/cards/{card_id}/details` | Получение полных реквизитов (через Vault) |
| | `POST` | `/api/cards/{card_id}/pin` | Установка / проверка PIN-кода |
| | `PATCH` | `/api/cards/{card_id}/status` | Блокировка / разблокировка карты |
| **Transactions**| `POST` | `/api/deposits` | Пополнение счета или карты |
| | `POST` | `/api/withdrawals` | Снятие наличных / списание средств |
| | `POST` | `/api/transfers` | Внутренний перевод (между счетами/картами) |
| **Recurring** | `POST` | `/api/recurring-payments` | Создание автоплатежа по CRON-расписанию |
| | `DELETE`| `/api/recurring-payments/{id}` | Отмена подписки / автоплатежа |
| **Internal** | `POST` | `/api/internal/payments/verify-card` | Верификация сырых реквизитов для платежных шлюзов |

> 📚 **Swagger UI:** При запущенном сервисе полная документация со схемами запросов и ответов доступна по адресу `http://localhost:8081/swagger/index.html`.

---

###  Внутренние интеграции (Microservices Communication)

Чтобы не превращать сервис в монолит и соблюсти стандарты безопасности (PCI-DSS), часть логики вынесена во внешние зависимости:

#### 1. Card Vault (gRPC)
Хранилище чувствительных данных. Account Service **никогда не сохраняет** полный номер карты (PAN), CVV или PIN-код в своей базе PostgreSQL. 
* При выпуске карты ядро отправляет gRPC-запрос в Vault, который генерирует данные, шифрует их и возвращает `card_id` (токен).
* При просмотре реквизитов в приложении ядро проксирует токен в Vault и отдает расшифрованные данные клиенту.

#### 2. Currency Exchanger (HTTP REST)
Сервис актуальных валютных курсов. Вызывается синхронно при попытке совершить кросс-валютный перевод (например, со счета в `RUB` на счет в `USD`). Ядро запрашивает курс конвертации и фиксирует его в таблице `postings`.

#### 3. Event Publisher (RabbitMQ)
Ядро реализует паттерн **Event-Driven**. При каждом важном бизнес-действии (создан перевод, пополнен баланс, выпущена карта) сервис асинхронно публикует `DomainEvent` в обменник RabbitMQ.
* Это позволяет другим сервисам (например, `Notification Service` для отправки Push/SMS или `Analytics Service` для построения отчетов) реагировать на события без жесткой связности с транзакционным ядром.

## Локальный запуск (Local Setup)

### Требования (Prerequisites)
* **Go** 1.26
* **Docker** и **Docker Compose** (для инфраструктуры и тестов)
* Утилита task
* Утилита [golang-migrate](https://github.com/golang-migrate/migrate) (для ручного наката миграций)

### 1. Подготовка инфраструктуры
Для работы ядра требуются PostgreSQL и RabbitMQ. В корне проекта (или в репозитории инфраструктуры) поднимите контейнеры:

```bash
# Поднимаем БД и брокер сообщений
docker-compose up -d postgres rabbitmq
```

### 2. Конфигурация
Сервис читает настройки из YAML-файла. Убедитесь, что у вас есть файл `config/local.yaml` (или передайте путь через переменную окружения `CONFIG_PATH`):

```yaml
db:
  url: "postgres://user:pass@localhost:5432/accounts_db?sslmode=disable"
rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
vault:
  address: "localhost:50051"
exchanger:
  url: "http://localhost:8082"
  timeout: 5s
```

### 3. Миграции
Накатите структуру базы данных:

```bash
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/accounts_db?sslmode=disable" up
```

### 4. Запуск сервиса
```bash
task up

task up:account (персональная команда сервиса)
```
*Сервис запустится на порту `8081`. API будет доступно по адресу `http://localhost:8081/api/`.*

---

## Тестирование (Testing)

В проекте реализован строгий подход к тестированию бизнес-логики. Мы не используем in-memory базы (sqlite) для транзакционного ядра, так как они не поддерживают специфичные фичи PostgreSQL (например, `FOR UPDATE SKIP LOCKED`).

Вместо этого мы используем **Testcontainers**.

### Как работают тесты:
1. При запуске `go test`, библиотека `testcontainers-go` автоматически скачивает образ `postgres:17-alpine` и поднимает чистый Docker-контейнер.
2. Автоматически применяются SQL-миграции из папки `migrations`.
3. Прогоняются Table-Driven тесты (включая проверку идемпотентности, транзакций, параллельного списания автоплатежей).
4. После завершения тестов контейнер с БД автоматически удаляется.

Для слоя инфраструктуры (Vault, RabbitMQ, Exchanger) используются сгенерированные интерфейсы и моки (Mockery / ручные моки).

### Запуск тестов:
Убедитесь, что демон Docker запущен, и выполните:

```bash
# Запуск всех тестов с подробным выводом
go test -v ./tests/...
```

---

## Возможные улучшения (Roadmap)
* [ ] **Transactional Outbox Pattern:** Сохранение событий RabbitMQ в отдельную таблицу `outbox_events` в рамках единой транзакции БД (для 100% гарантии доставки при падении брокера).
* [ ] **Observability:** Интеграция с Prometheus (сбор метрик бизнес-логики: количество переводов, отказов) и Jaeger (OpenTelemetry) для распределенной трассировки запросов между Account, Vault и Exchanger сервисами.