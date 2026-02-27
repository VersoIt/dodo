# 🍕 Pizza Delivery System — Microservices Architecture

[![CI Pipeline](https://github.com/versoit/diploma/actions/workflows/ci.yml/badge.svg)](https://github.com/versoit/diploma/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.x-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Профессиональная микросервисная экосистема для автоматизации службы доставки еды. Проект спроектирован с учетом требований к высокой нагрузке, отказоустойчивости и консистентности данных в распределенных системах.

## 🏗 Архитектурные принципы

Проект построен на базе **современных** паттернов и технологий:

*   **Clean Architecture & DDD**: Четкое разделение на слои (Domain, UseCase, Infrastructure), использование Value Objects и Aggregate Roots.
*   **Event-Driven Architecture**: Асинхронное взаимодействие через **NATS** вместо хрупких gRPC-цепочек.
*   **Transactional Outbox**: Гарантия "at-least-once delivery" событий. Заказ и событие сохраняются в БД атомарно.
*   **Go Workspaces**: Современное управление зависимостями в монорепозитории.
*   **Observability Stack**: Полная прозрачность системы через OpenTelemetry, Prometheus, Tempo (Tracing) и Loki (Logging).
*   **Performance**: Отсутствие проблем N+1 в репозиториях (оптимизированные SQL-запросы).

## 🛰 Технологический стек

| Слой | Технологии |
| :--- | :--- |
| **Backend** | Go 1.24, gRPC, NATS, PostgreSQL, Squirrel (Query Builder), Uber FX (DI) |
| **Frontend** | Vue.js 3, TypeScript, TailwindCSS, DaisyUI, Axios, Pinia |
| **DevOps** | Docker, Docker Compose, GitHub Actions (CI), Goose (Migrations) |
| **Monitoring** | Grafana, Prometheus, OpenTelemetry (OTEL), Tempo, Loki |

## 📂 Структура проекта

```text
├── be/                 # Backend микросервисы
│   ├── auth/           # Авторизация и JWT
│   ├── catalog/        # Меню и ингредиенты
│   ├── orders/         # Ядро системы (Outbox, Relay)
│   ├── kitchen/        # Модуль кухни
│   ├── logistics/      # Доставка и курьеры
│   ├── notification/   # Push/Email уведомления
│   ├── treasury/       # Финансы и платежи
│   ├── chat/           # WebSocket чат между ролями
│   └── gateway/        # API Gateway (REST -> gRPC)
├── fe/                 # Frontend (Vue.js 3 SPA)
├── pkg/                # Общие библиотеки (Logger, DB, Telemetry)
└── deploy/             # Конфиги для Kubernetes и мониторинга
```

## 🚀 Быстрый старт

### Требования
*   Docker & Docker Compose
*   Go 1.24+ (для локальной разработки)

### Запуск всей системы
```bash
# Клонировать репозиторий
git clone https://github.com/versoit/diploma.git
cd diploma

# Запустить все сервисы в Docker
docker-compose up --build -d
```

Система будет доступна по адресам:
*   **Frontend**: `http://localhost:3000`
*   **API Gateway**: `http://localhost:8080`
*   **Grafana**: `http://localhost:3001` (панели мониторинга)
*   **Prometheus**: `http://localhost:9090`

## 📊 Основные фичи

1.  **Умный экспорт**: Генерация Excel-отчетов с использованием стриминга (`StreamWriter`), что позволяет выгружать миллионы строк без утечек памяти.
2.  **Надежные транзакции**: Благодаря Outbox паттерну, даже если NATS упадет, события будут доставлены сразу после его поднятия.
3.  **Real-time чат**: Прямая связь клиента с курьером и поддержкой через WebSockets.
4.  **Аналитика**: Автоматический расчет среднего чека, топ-товаров и выручки.

## 🛠 Разработка

### Генерация Proto-файлов
```bash
task proto
```

### Запуск линтера
```bash
task lint
```

### Запуск тестов
```bash
task test
```

---
⭐ *Разработано в рамках дипломного проекта с упором на современные архитектурные паттерны.*
