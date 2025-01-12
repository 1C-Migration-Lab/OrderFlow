# Структура проекта учета заказов

## 1. Архитектура системы

```
[Frontend (HTML + JavaScript)] 
           ↓ HTTP/REST
[Backend (Go + PostgreSQL)]
           ↓ SQL
[Database (PostgreSQL)]
```

## 2. Структура базы данных

База данных уже определена в migrations/001_initial_schema.sql и включает:

### Основные таблицы:
1. clients (клиенты)
   - id: SERIAL PRIMARY KEY
   - name: VARCHAR(255) NOT NULL
   - inn: VARCHAR(50)

2. products (товары)
   - id: SERIAL PRIMARY KEY
   - name: VARCHAR(255) NOT NULL
   - unit: VARCHAR(50) NOT NULL

3. orders (заказы)
   - id: SERIAL PRIMARY KEY
   - client_id: INTEGER (FK -> clients)
   - date: TIMESTAMP
   - number: VARCHAR(50) UNIQUE
   - total_amount: DECIMAL(15,2)
   - is_confirmed: BOOLEAN
   - created_at: TIMESTAMP

4. order_items (позиции заказа)
   - id: SERIAL PRIMARY KEY
   - order_id: INTEGER (FK -> orders)
   - product_id: INTEGER (FK -> products)
   - quantity: DECIMAL(15,3)
   - price: DECIMAL(15,2)
   - line_amount: DECIMAL(15,2)

5. orders_by_client (агрегация сумм)
   - client_id: INTEGER (FK -> clients)
   - orders_sum: DECIMAL(15,2)

## 3. Структура Backend (Go)

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Точка входа приложения
├── internal/
│   ├── api/                  # HTTP handlers и роутинг
│   │   ├── handlers/
│   │   │   ├── clients.go    # Обработчики для клиентов
│   │   │   ├── products.go   # Обработчики для товаров
│   │   │   └── orders.go     # Обработчики для заказов
│   │   └── router.go         # Настройка маршрутизации
│   ├── models/               # Структуры данных
│   │   └── models.go         # Определения моделей
│   ├── repository/           # Слой доступа к данным
│   │   └── repository.go     # Интерфейсы и реализация
│   └── service/              # Бизнес-логика
│       └── service.go        # Сервисные функции
├── pkg/                      # Публичные пакеты
│   └── utils/               # Утилиты
└── config/                  # Конфигурация
    └── config.go           # Настройки приложения
```

## 4. Структура Frontend

```
frontend/
├── index.html              # Главная страница
├── css/                    # Стили
│   └── styles.css
├── js/                     # JavaScript модули
│   ├── api.js             # Работа с API
│   ├── clients.js         # Логика работы с клиентами
│   ├── products.js        # Логика работы с товарами
│   └── orders.js          # Логика работы с заказами
└── components/            # HTML компоненты
    ├── client-form.html   # Форма клиента
    ├── product-form.html  # Форма товара
    └── order-form.html    # Форма заказа
```

## 5. API Endpoints

### Клиенты
- GET /api/clients - список клиентов
- POST /api/clients - создание клиента
- GET /api/clients/{id} - получение клиента
- PUT /api/clients/{id} - обновление клиента
- DELETE /api/clients/{id} - удаление клиента

### Товары
- GET /api/products - список товаров
- POST /api/products - создание товара
- GET /api/products/{id} - получение товара
- PUT /api/products/{id} - обновление товара
- DELETE /api/products/{id} - удаление товара

### Заказы
- GET /api/orders - список заказов
- POST /api/orders - создание заказа
- GET /api/orders/{id} - получение заказа
- PUT /api/orders/{id} - обновление заказа
- DELETE /api/orders/{id} - удаление заказа
- POST /api/orders/{id}/confirm - подтверждение заказа

### Агрегация
- GET /api/orders-by-client - суммы заказов по клиентам

## 6. Основные компоненты фронтенда

### Страницы
1. Список клиентов
   - Таблица с колонками: ID, Название, ИНН
   - Кнопки: Создать, Редактировать, Удалить

2. Список товаров
   - Таблица с колонками: ID, Название, Единица измерения
   - Кнопки: Создать, Редактировать, Удалить

3. Список заказов
   - Таблица с колонками: ID, Номер, Клиент, Дата, Сумма, Статус
   - Кнопки: Создать, Редактировать, Подтвердить, Удалить

4. Форма заказа
   - Шапка: выбор клиента, дата, номер
   - Табличная часть: выбор товара, количество, цена
   - Автоматический расчет сумм
   - Кнопки: Сохранить, Отмена

5. Виджет сумм по клиентам
   - Таблица: Клиент, Общая сумма заказов

## 7. Технологический стек

### Backend
- Go 1.18+
- PostgreSQL
- Gin Web Framework
- GORM (ORM для Go)

### Frontend
- Vanilla JavaScript (без фреймворков)
- HTML5
- CSS3

## 8. Безопасность и валидация

### Backend
- Валидация входных данных
- Проверка бизнес-правил
- Защита от SQL-инъекций (через ORM)
- Логирование операций

### Frontend
- Валидация форм
- Обработка ошибок API
- Подтверждение удаления
- Блокировка подтвержденных заказов

## 9. Дальнейшие шаги разработки

1. Создание базовой структуры проекта
2. Реализация backend API
3. Разработка frontend интерфейса
4. Тестирование
5. Документирование
