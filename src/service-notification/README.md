# 🔔 Notification Service (Event Consumer)

[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)](https://go.dev/)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-Event%20Driven-FF6600?style=flat&logo=rabbitmq)](https://www.rabbitmq.com/)
[![Architecture](https://img.shields.io/badge/Architecture-Pub%2FSub-4CAF50?style=flat)]()

**Notification Service** - это легковесный асинхронный микросервис в экосистеме [Banking System](https://github.com/Adopten123/banking-system), отвечающий за обработку доменных событий и отправку уведомлений пользователям (Push, SMS, Email).

Сервис построен по паттерну **Event-Driven Architecture**. Он не имеет собственной базы данных и не блокирует транзакционное ядро (Account Service). Вместо этого он "слушает" шину сообщений RabbitMQ и реагирует на бизнес-события в реальном времени.

---

## Ключевые возможности (Features)

* **Асинхронная обработка (Async Processing):** События вычитываются из очереди RabbitMQ в фоновых горутинах, обеспечивая высокую пропускную способность.
* **Динамический роутинг событий:** Реализован паттерн `Event Dispatcher` — входящие сообщения автоматически маршрутизируются в нужный хендлер на основе поля типа события (например, `AccountCreatedEvent` -> `handleAccountCreated`).
* **Мультиканальность (Симуляция):** Сервис определяет оптимальный канал доставки в зависимости от типа и важности события:
    * **[EMAIL]** — Создание аккаунтов, изменение кредитных лимитов.
    *  **[SMS]** — Изменение статусов счетов, выпуск новых карт.
    *  **[PUSH]** — Финансовые операции (переводы, пополнения, списания), блокировка карт.
* **Изоляция сбоев:** Падение или перезагрузка сервиса нотификаций никак не влияет на работоспособность ядра системы. При перезапуске сервис просто вычитает накопившиеся в очереди сообщения.
* **Graceful Shutdown:** Корректное завершение работы при получении сигналов `SIGINT`/`SIGTERM` с безопасным закрытием соединений и каналов RabbitMQ.

---

## Стек технологий

* **Язык:** Go 1.26
* **Брокер сообщений:** RabbitMQ
* **Клиент AMQP:** `github.com/rabbitmq/amqp091-go`
* **Формат данных:** JSON

---

## Domain Events

Сервис обрабатывает следующий пул событий, генерируемых транзакционным ядром:

| Событие | Триггер в ядре | Канал уведомления |
| :--- | :--- | :--- |
| `AccountCreatedEvent` | Открытие нового банковского счета | 📧 Email |
| `AccountStatusChangedEvent` | Заморозка, блокировка или закрытие счета | 💬 SMS |
| `CreditLimitChangedEvent` | Выделение или изменение овердрафта | 📧 Email |
| `CardIssuedEvent` | Выпуск новой физической или виртуальной карты | 💬 SMS |
| `CardStatusChangedEvent` | Блокировка или разблокировка карты | 📱 Push |
| `DepositCompletedEvent` | Успешное пополнение баланса | 📱 Push |
| `WithdrawalCompletedEvent` | Успешное снятие наличных / оплата | 📱 Push |
| `TransferCreatedEvent` | Перевод средств (содержит данные отправителя и получателя) | 📱 Push |

---

##  Архитектура взаимодействия (Pub/Sub)

```text
[ Account Service ] --(JSON Event)--> [ RabbitMQ Exchange (fanout) ]
                                                |
                                                +--> [ Queue: notifications ]
                                                            |
                                                            v
                                            [ Notification Service ]
                                            /          |           \
                                       [SMS]        [EMAIL]      [PUSH]
```

При запуске консьюмер (`NotificationConsumer`) автоматически декларирует Fanout-обменник, создает очередь и выполняет Binding (связывание). Это гарантирует, что сервис готов к приему сообщений сразу после старта, даже если ядро еще не было запущено.

---

## Локальный запуск (Local Setup)

### 1. Подготовка RabbitMQ
Сервису требуется запущенный экземпляр RabbitMQ. Вы можете поднять его через Docker:
```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

### 2. Конфигурация
Настройки задаются через переменные окружения или конфигурационный файл:
```yaml
# config/local.yaml
env: "local"
rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
  exchange: "banking_events_exchange"
  queue: "notification_queue"
```

### 3. Запуск
```bash
go mod tidy
go run cmd/main.go
```

При успешном старте вы увидите в консоли:
```text
Successfully connected to RabbitMQ
[*] Waiting for messages in queue: notification_queue
Notification service is successfully running!
```

---

## Возможные улучшения (Roadmap)
* Интеграция с реальными провайдерами (Twilio для SMS, SendGrid для Email, Firebase Cloud Messaging для Push).
* Реализация механизма **Retry/Dead Letter Queue (DLQ)** для повторной отправки уведомлений при временной недоступности внешних провайдеров доставки.
* Шаблонизатор (например, `html/template`) для генерации красивых HTML-писем.