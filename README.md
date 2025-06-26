# Woman App Backend - Шаблон Go API с Keycloak и PostgreSQL

Этот проект представляет собой готовый шаблон для создания Go API с авторизацией через Keycloak и базой данных PostgreSQL.

## 🛠 Стек технологий

- **Backend**: Go 1.23+ с Chi router
- **База данных**: PostgreSQL 16
- **Авторизация**: Keycloak 24
- **Документация**: Swagger UI
- **Контейнеризация**: Docker & Docker Compose
- **Тестирование**: Go testing + Testify
- **Логирование**: Zap logger

## 📋 Что входит в шаблон

- ✅ Настроенный Keycloak с realm и client
- ✅ PostgreSQL с миграциями и архивированием
- ✅ API сервер с middleware авторизации
- ✅ Debug сервер для разработки
- ✅ Swagger UI с документацией
- ✅ Тестовая база данных
- ✅ Готовые эндпоинты регистрации и авторизации
- ✅ Валидация пользовательских данных
- ✅ Настроенный CI/CD pipeline

## 🚀 Быстрый старт

### 1. Клонирование и настройка

```bash
# Клонируем репозиторий
git clone <your-repo-url>
cd woman-app

# Копируем конфигурационные файлы
cp configs/config.example.yaml configs/config.yaml

# Создаем переменные окружения (если нужно)
cp .env.example .env
```

### 2. Настройка конфигурации

Отредактируйте `configs/config.yaml`:

```yaml
# Измените на ваши настройки
keycloak:
  url: "http://localhost:38081"
  realm: "YourRealm"          # Измените
  client_id: "your-client"    # Измените

database:
  host: "localhost"
  port: 37546
  user: "your_user"           # Измените
  password: "your_password"   # Измените
  dbname: "your_db"           # Измените

server:
  port: ":38080"
  debug_port: ":38081"
```

### 3. Запуск инфраструктуры

```bash
# Поднимаем всю инфраструктуру
docker-compose -f deploy/dev/docker-compose.yml up -d

# Ждем запуска служб (30-60 секунд)
docker-compose -f deploy/dev/docker-compose.yml logs -f
```

### 4. Настройка Keycloak

1. Откройте Keycloak Admin Console: http://localhost:38081
2. Войдите как admin (логин/пароль в docker-compose.yml)
3. Создайте новый realm или используйте существующий
4. Настройте client для вашего API
5. Обновите конфигурацию в `configs/config.yaml`

### 5. Инициализация базы данных

```bash
# Применяем миграции
task db:migrate:up

# Или через Docker
docker exec -it woman-app-db psql -U youruser -d yourdb -f /migrations/001_initial_schema.sql
```

### 6. Сборка и запуск API

```bash
# Сборка приложения
go build -o main-service ./cmd/main-service

# Запуск API сервера
./main-service

# Или через Docker
docker-compose -f deploy/dev/docker-compose.yml up api-server
```

### 7. Проверка работоспособности

```bash
# Проверка health endpoint
curl http://localhost:38080/health

# Проверка Swagger UI
open http://localhost:38080/docs/

# Проверка регистрации
curl -X POST http://localhost:38080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## 📁 Структура проекта

```
woman-app/
├── api/                    # HTTP handlers
│   ├── auth/              # Авторизация
│   └── user/              # Пользователи
├── cmd/                   # Точки входа
│   └── main-service/      # Основной сервис
├── configs/               # Конфигурации
├── deploy/                # Docker и развертывание
│   ├── dev/              # Development окружение
│   └── local/            # Локальная разработка
├── internal/              # Внутренняя логика
│   ├── clients/          # Внешние клиенты (Keycloak)
│   ├── config/           # Парсинг конфигурации
│   ├── middlewares/      # HTTP middleware
│   ├── models/           # Модели данных
│   ├── service/          # Бизнес-логика
│   └── store/            # Слой данных
├── docs/                  # Документация
└── tests/                 # Тесты
```

## 🛠 Разработка

### Добавление новых эндпоинтов

1. Создайте handler в `api/`
2. Добавьте бизнес-логику в `internal/service/`
3. Добавьте работу с БД в `internal/store/`
4. Обновите документацию в `deploy/dev/docs/api.yaml`
5. Добавьте тесты

### Работа с базой данных

```bash
# Создание новой миграции
task db:migration:create NAME=add_new_table

# Применение миграций
task db:migrate:up

# Откат миграций
task db:migrate:down

# Архивирование данных
task db:backup
```

### Тестирование

```bash
# Запуск всех тестов
go test ./...

# Тесты с покрытием
go test -cover ./...

# Запуск тестовой БД
task db:test:up
```

## 🐳 Docker команды

```bash
# Полная пересборка
docker-compose -f deploy/dev/docker-compose.yml build --no-cache

# Логи конкретного сервиса
docker-compose -f deploy/dev/docker-compose.yml logs -f api-server

# Остановка всех сервисов
docker-compose -f deploy/dev/docker-compose.yml down

# Очистка данных
docker-compose -f deploy/dev/docker-compose.yml down -v
```

## 🔧 Адаптация под новый проект

### 1. Переименование модуля

```bash
# Замените все упоминания модуля
find . -name "*.go" -exec sed -i 's/github.com\/Fisher-Development\/woman-app-backend/your-new-module/g' {} \;

# Обновите go.mod
go mod edit -module your-new-module
go mod tidy
```

### 2. Настройка портов

Измените порты в:
- `deploy/dev/docker-compose.yml`
- `configs/config.yaml`
- `deploy/dev/docs/api.yaml`

### 3. Настройка Keycloak

- Создайте новый realm
- Настройте client_id и секреты
- Обновите конфигурацию

### 4. Кастомизация API

- Замените модели в `internal/models/`
- Обновите валидацию в `internal/service/validation.go`
- Добавьте новые эндпоинты

## 📊 Мониторинг

Доступные эндпоинты для мониторинга:

- `GET /health` - Проверка здоровья сервиса
- `GET /info` - Информация о сервисе
- `GET /api/v1/version` - Версия API

## 🔐 Безопасность

- Все защищенные эндпоинты требуют JWT токен
- Валидация через Keycloak
- CORS настроен для development
- Готовые middleware для авторизации

## 📝 Примеры API

### Регистрация пользователя
```bash
curl -X POST http://localhost:38080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Авторизация
```bash
curl -X POST http://localhost:38080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Обновление профиля
```bash
curl -X PUT http://localhost:38080/api/v1/user/update \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"firstName":"John","lastName":"Doe"}'
```

## 🤝 Участие в разработке

1. Fork проекта
2. Создайте feature branch
3. Внесите изменения и добавьте тесты
4. Запустите тесты и линтеры
5. Создайте Pull Request

## 📄 Лицензия

MIT License - подробности в файле LICENSE

## 📞 Поддержка

- Issues: создавайте issue в GitHub
- Email: support@your-domain.com
- Документация: доступна в Swagger UI

---

**🎯 Готово к продакшену:** Этот шаблон включает все необходимые компоненты для создания production-ready API сервиса.
