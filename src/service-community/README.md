#  Community Service (Pulse & Social)

[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-316192?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=flat&logo=redis)](https://redis.io/)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-Event%20Driven-FF6600?style=flat&logo=rabbitmq)](https://www.rabbitmq.com/)
[![WebSockets](https://img.shields.io/badge/WebSockets-Real--time-blue?style=flat)](#)

**Community Service** - это социальное ядро экосистемы [Banking System](https://github.com/Adopten123/banking-system). Сервис объединяет в себе функционал социальной сети для инвесторов (аналог Пульса), новостную ленту банка и защищенный мессенджер с поддержкой P2P-переводов.

Архитектура построена по принципам **Clean Architecture**. Сервис не управляет реальными балансами, а выступает инициатором транзакций, делегируя финансовые операции в `Account Service`.

---

## Ключевые возможности (Features)

* **Профили и Социальный граф:**
    * Асинхронная синхронизация базовых данных из SSO-сервиса через RabbitMQ.
    * Настройки приватности (закрытые/открытые профили).
    * Система подписок (Followers/Following) для формирования персональной ленты.
* **Лента контента (Feed):**
    * Публикация официальных новостей банка (с пинами) и пользовательских постов.
    * Привязка постов к конкретным инвестиционным активам (тикерам).
    * Лайки, счетчики и иерархические комментарии (с поддержкой ответов).
* **Мессенджер (Чаты):**
    * Личные (P2P), групповые чаты и каналы технической поддержки.
    * *Real-time* взаимодействие через WebSockets с кэшированием сессий в Redis.
    * Передача медиа-вложений (чеки, фото, документы).
* **Переводы в сообщениях (P2P Transfers):**
    * Возможность отправить деньги напрямую в чате.
    * Сервис генерирует `idempotency_key` и отправляет REST-запрос в банковское ядро, после чего асинхронно обновляет статус сообщения (pending -> completed) через прослушивание шины RabbitMQ.
* **Модерация и Безопасность:**
    * Встроенная система жалоб (Reports) на спам, мошенничество или оскорбления.
    * Статусная модель для контента (`published`, `blocked`, `under_review`).

---

## Стек технологий (Tech Stack)

* **Язык:** Go 1.26
* **Роутинг и API:** REST API (`go-chi/chi`) + WebSockets (`gorilla/websocket`)
* **База данных:** PostgreSQL 17 (драйвер `pgx`, кодогенерация `sqlc`)
* **Кэширование:** Redis 7+ (кэш ленты новостей, сессии WS, счетчики непрочитанных)
* **Брокер сообщений:** RabbitMQ (публикация и подписка на доменные события)
* **Инфраструктура:** Docker, Taskfile, Prometheus + Grafana

---

## Интеграции с другими сервисами

1. **SSO + Auth Service:** Получение профилей юзеров при регистрации (через RabbitMQ) и валидация JWT-токенов при каждом API-вызове.
2. **Account Service (Core):** Инициация безопасных денежных переводов из чатов с защитой от двойного списания.
3. **Investment and Exchange Service:** Валидация тикеров (например, `AAPL`, `GAZP`) при их упоминании в инвестиционных идеях.
4. **Notification Service:** Отправка пуш-уведомлений о новых сообщениях, подписчиках и комментариях через RabbitMQ.

---

## Архитектура Базы Данных

База данных спроектирована для работы с графовыми и иерархическими связями, типичными для социальных сетей.

### ER-Диаграмма

```mermaid
erDiagram
    USERS_PROFILE ||--o{ FOLLOWS : "follows"
    USERS_PROFILE ||--o{ POSTS : "writes"
    USERS_PROFILE ||--o{ COMMENTS : "writes"
    USERS_PROFILE ||--o{ CHAT_MEMBERS : "joins"
    USERS_PROFILE ||--o{ REPORTS : "files"
    
    POSTS ||--o{ POST_LIKES : "receives"
    POSTS ||--o{ COMMENTS : "has"
    POST_TYPES ||--o{ POSTS : "categorizes"
    
    CHATS ||--o{ CHAT_MEMBERS : "contains"
    CHATS ||--o{ MESSAGES : "contains"
    CHAT_TYPES ||--o{ CHATS : "categorizes"

    USERS_PROFILE {
        uuid user_id PK "SSO ID"
        varchar username
        boolean is_private "Privacy setting"
    }
    POSTS {
        bigint id PK
        uuid author_id FK
        text content
        varchar status "published/blocked"
        varchar related_asset_ticker "Link to Investment API"
    }
    COMMENTS {
        bigint id PK
        bigint parent_comment_id FK "For replies"
        text content
    }
    MESSAGES {
        bigint id PK
        uuid chat_id FK
        uuid sender_id FK
        jsonb media_attachments
        boolean is_transfer
        varchar idempotency_key "Unique for safe retries"
        varchar transfer_status "pending/completed"
    }
    REPORTS {
        bigint id PK
        uuid reporter_id FK
        varchar target_type "post/comment/message"
        bigint target_id
        varchar status "open/resolved"
    }
```

### DBML (Полная структура)

```dbml
// --- ПРОФИЛИ (Проекция из SSO) ---
Table users_profile {
  user_id uuid [pk, note: 'Внешний ключ на SSO ID']
  username varchar(50) [unique]
  display_name varchar(100)
  avatar_url varchar(255)
  bio text
  
  is_verified boolean [default: false, note: 'Официальный аккаунт']
  is_staff boolean [default: false, note: 'Сотрудник банка (Support)']
  is_private boolean [default: false, note: 'Скрытый профиль']
  
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
}

// --- СОЦИАЛЬНЫЙ ГРАФ ---
Table follows {
  follower_id uuid [ref: > users_profile.user_id]
  following_id uuid [ref: > users_profile.user_id]
  created_at timestamp [default: `now()`]

  Indexes {
    (follower_id, following_id) [unique]
  }
}

// --- ЛЕНТА И КОНТЕНТ ---
Table post_types {
  id int [pk, increment]
  name varchar [note: 'news, user_post, investment_idea']
}

Table posts {
  id bigserial [pk]
  author_id uuid [ref: > users_profile.user_id]
  type_id int [ref: > post_types.id]
  
  content text
  media_attachments jsonb [null, note: 'Array of URLs']
  related_asset_ticker varchar(10) [null, note: 'Например: AAPL']
  
  status varchar(20) [default: 'published', note: 'published, blocked, under_review']
  likes_count int [default: 0]
  comments_count int [default: 0]
  is_pinned boolean [default: false, note: 'Для закрепления новостей банка']
  is_edited boolean [default: false]
  
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  deleted_at timestamp [null]

  Indexes {
    author_id
    related_asset_ticker
  }
}

Table post_likes {
  post_id bigint [ref: > posts.id]
  user_id uuid [ref: > users_profile.user_id]
  created_at timestamp [default: `now()`]
  
  Indexes {
    (post_id, user_id) [unique]
  }
}

Table comments {
  id bigserial [pk]
  post_id bigint [ref: > posts.id]
  user_id uuid [ref: > users_profile.user_id]
  parent_comment_id bigint [ref: > comments.id, null, note: 'Дерево ответов']
  
  content text
  status varchar(20) [default: 'published']
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  deleted_at timestamp [null]
}

// --- МОДЕРАЦИЯ ---
Table reports {
  id bigserial [pk]
  reporter_id uuid [ref: > users_profile.user_id]
  target_type varchar(20) [note: 'post, comment, user, message']
  target_id bigint [note: 'ID подозрительного контента']
  reason text
  status varchar(20) [default: 'open', note: 'open, in_progress, resolved']
  created_at timestamp [default: `now()`]
}

// --- МЕССЕНДЖЕР И ПЕРЕВОДЫ ---
Table chat_types {
  id int [pk, increment]
  name varchar [note: 'private, group, support_channel']
}

Table chats {
  id uuid [pk, default: `gen_random_uuid()`]
  type_id int [ref: > chat_types.id]
  
  title varchar(100) [null]
  avatar_url varchar [null]
  
  last_message_at timestamp [default: `now()`]
  created_at timestamp [default: `now()`]
}

Table chat_members {
  chat_id uuid [ref: > chats.id]
  user_id uuid [ref: > users_profile.user_id]
  
  role varchar(20) [default: 'member']
  joined_at timestamp [default: `now()`]
  last_read_message_id bigint [null]
  
  Indexes {
    (chat_id, user_id) [unique]
  }
}

Table messages {
  id bigserial [pk]
  chat_id uuid [ref: > chats.id]
  sender_id uuid [ref: > users_profile.user_id]
  reply_to_message_id bigint [ref: > messages.id, null]
  
  content text [null]
  media_attachments jsonb [null, note: 'Медиа (чеки, фото)']

  // БЛОК ФИНАНСОВ
  is_transfer boolean [default: false]
  transfer_amount numeric(20, 0) [null]
  transfer_currency varchar(3) [null]
  idempotency_key varchar(255) [null, note: 'Защита от двойного списания']
  transfer_transaction_id uuid [null, note: 'ID из bank_core.transactions']
  transfer_status varchar(20) [null, note: 'pending, completed, failed']

  is_edited boolean [default: false]
  created_at timestamp [default: `now()`]
  deleted_at timestamp [null]
  
  Indexes {
    (chat_id, id)
  }
}
```

---

## Локальный запуск (Local Setup)

### Требования
* Go 1.26+
* Docker & Docker Compose
* Утилита `task`

### 1. Подготовка инфраструктуры
Сервису требуются PostgreSQL, Redis (для WebSockets/Кэша) и RabbitMQ:
```bash
docker-compose up -d postgres redis rabbitmq
```

### 2. Конфигурация
Создайте/проверьте файл `config/local.yaml`:
```yaml
env: "local"
http:
  port: 8083
db:
  url: "postgres://user:pass@localhost:5432/community_db?sslmode=disable"
redis:
  addr: "localhost:6379"
rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
```

### 3. Запуск
```bash
# Накатить миграции и запустить сервис
task up:community
```

---

## Возможные улучшения
* [ ] **Алгоритмическая лента:** Переход от хронологического Feed к умной ленте на основе графов интересов и Machine Learning.
* [ ] **Full-Text Search:** Интеграция с **Elasticsearch** для быстрого поиска постов, тикеров и пользователей.
* [ ] **Anti-Spam:** Внедрение Rate Limiting и автоматической проверки сообщений нейросетью перед публикацией.