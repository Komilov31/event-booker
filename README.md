# Event Booker

Система бронирования мероприятий - это полнофункциональное веб-приложение для создания и управления событиями, бронирования мест и обработки платежей. Приложение построено с использованием архитектуры микросервисов и предоставляет REST API с веб-интерфейсом.

## 🛠 Технологии

### Backend
- **Go** - основной язык программирования
- **Gin Framework** - HTTP веб-фреймворк для Go
- **PostgreSQL** - основная база данных
- **RabbitMQ** - очередь сообщений для асинхронной обработки
- **Goose** - инструмент для миграций базы данных
- **Zerolog** - структурированное логирование

### Frontend
- **HTML5/CSS3/JavaScript** - веб-интерфейс

### DevOps & Infrastructure
- **Docker** - контейнеризация приложения
- **Docker Compose** - оркестрация многоконтейнерных приложений
- **Swagger** - документация API

### Дополнительные библиотеки
- **Telegram Bot API** - интеграция с Telegram
- **WB Framework** - фреймворк от WB для легкой работы с разными инструментами
- **pq** - драйвер PostgreSQL для Go
- **AMQP** - клиент для RabbitMQ

## 🏗 Архитектура

Приложение построено с использованием **чистой архитектуры** (Clean Architecture) с четким разделением слоев:

### Слои архитектуры

1. **Handler Layer** (`internal/handler/`)
   - HTTP обработчики запросов
   - Валидация входящих данных
   - Формирование ответов
   - Swagger документация

2. **Service Layer** (`internal/service/`)
   - Бизнес-логика приложения
   - Координация между слоями
   - Обработка очередей сообщений
   - Валидация бизнес-правил

3. **Repository Layer** (`internal/repository/`)
   - Доступ к данным
   - Абстракция от конкретной БД

4. **Model Layer** (`internal/model/`)
   - Структуры данных
   - Сущности предметной области

5. **DTO Layer** (`internal/dto/`)
   - Data Transfer Objects
   - Структуры для API запросов/ответов

### Дополнительные компоненты

- **Config** (`internal/config/`) - управление конфигурацией
- **RabbitMQ** (`internal/rabbitmq/`) - работа с очередями сообщений
- **Sender** (`internal/sender/`) - отправка уведомлений


## 🚀 API Endpoints

### POST Endpoints

#### Создание нового мероприятия
```bash
curl -X POST "http://localhost:8080/events" \
  -H "Content-Type: application/json" \
  -d '{
    "event_name": "Концерт группы Rock",
    "event_at": "2024-12-31T20:00:00Z",
    "all_seats": 100
  }'
```

#### Бронирование места на мероприятие
```bash
curl -X POST "http://localhost:8080/events/{event_id}/book" \
  -H "Content-Type: application/json" \
  -d '{
    "telegram_id": 123456789,
    "places_count": 2
  }'
```

#### Подтверждение оплаты бронирования
```bash
curl -X POST "http://localhost:8080/events/{booking_id}/confirm" \
  -H "Content-Type: application/json"
```

### GET Endpoints

#### Получение всех мероприятий
```bash
curl -X GET "http://localhost:8080/events" \
  -H "Content-Type: application/json"
```

#### Получение конкретного мероприятия
```bash
curl -X GET "http://localhost:8080/events/{id}" \
  -H "Content-Type: application/json"
```

#### Получение информации о бронировании
```bash
curl -X GET "http://localhost:8080/events/{id}" \
  -H "Content-Type: application/json"
```

#### Главная страница
```bash
curl -X GET "http://localhost:8080/"
```

#### Админ-панель
```bash
curl -X GET "http://localhost:8080/admin"
```

#### Пользовательская страница
```bash
curl -X GET "http://localhost:8080/user"
```

#### Swagger документация
```bash
curl -X GET "http://localhost:8080/swagger/index.html"
```


## 🐳 Запуск с Docker

### Предварительные требования
- Docker
- Docker Compose

### Запуск приложения
```bash
# Клонирование репозитория
git clone https://github.com/Komilov31/event-booker
cd event-booker

# Создание .env файла
cp .env.example .env

# Запуск всех сервисов
docker-compose up -d

# Просмотр логов
docker-compose logs -f app
```

### Остановка приложения
```bash
docker-compose down
```

## 📁 Структура проекта

```
event-booker/
├── cmd/                    # Точки входа в приложение
│   ├── main.go            # Главная функция
│   └── app/               # Инициализация приложения
├── internal/              # Внутренние пакеты
│   ├── config/           # Конфигурация
│   ├── dto/              # Data Transfer Objects
│   ├── handler/          # HTTP обработчики
│   ├── model/            # Модели данных
│   ├── repository/       # Репозитории
│   ├── service/          # Сервисы
│   ├── rabbitmq/         # Работа с очередями
│   └── sender/           # Отправка уведомлений
├── migrations/           # Миграции базы данных
├── static/               # Статические файлы
├── config/               # Конфигурационные файлы
├── docs/                 # Swagger документация
├── docker-compose.yml    # Docker Compose
├── Dockerfile           # Dockerfile
├── go.mod               # Модули Go
└── README.md            # Документация
```

## 🔄 Работа с очередями

Приложение использует RabbitMQ для асинхронной обработки:

1. **Создание бронирования** → отправка в очередь данные о брони через 15 минту после заказа
2. **Обработка очереди** → проверка статуса оплаты бронирования
3. **Подтверждение/отмена** → обновление статуса бронирования(отмена в случае не оплаты)

## 🌐 Веб-интерфейс

Приложение предоставляет три основных интерфейса:

1. **Главная страница** (`/`) - выбор режима работы
2. **Админ-панель** (`/admin`) - управление мероприятиями
3. **Пользовательский интерфейс** (`/user`) - бронирование мест

## 📚 Документация API

Swagger документация доступна по адресу:
```
http://localhost:8080/swagger/index.html
```
