# Название проекта
### Автор
Автор: [Nipup] (https://github.com/Nipup1)
### Обзор проекта
Этот проект представляет собой распределенное приложение, состоящее из нескольких микросервисов и фронтенда, предназначенное для управления событиями, обмена сообщениями и аутентификации.

### Структура проекта
Проект организован в следующие директории:

- calendar-service: Микросервис, отвечающий за создание событий и отправку напоминаний через WebSocket. Управляет расписанием событий и уведомлениями.
- event-messenger-frontend: Клиентская часть приложения, предоставляющая пользовательский интерфейс для взаимодействия с функциями событий и сообщений.
- messenger-service: Микросервис, обрабатывающий чаты и отправку запросов. Реализует REST API и WebSocket для реального времени, а также интегрируется с сервисом SSO для аутентификации.
- sso: Сервис для аутентификации и авторизации, обеспечивающий безопасный доступ к приложению и его микросервисам.