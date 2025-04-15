# React + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react/README.md)
  uses [Babel](https://babeljs.io/) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react-swc) uses [SWC](https://swc.rs/) for Fast
  Refresh

## Expanding the ESLint configuration

If you are developing a production application, we recommend using TypeScript and enable type-aware lint rules. Check
out the [TS template](https://github.com/vitejs/vite/tree/main/packages/create-vite/template-react-ts) to integrate
TypeScript and [`typescript-eslint`](https://typescript-eslint.io) in your project.

# Описание функционала фронтенда

## 1. Авторизация/регистрация

- **Режимы:**
    - **Обычная регистрация:** Пользователь регистрируется с помощью логина и пароля (данные хранятся в основном
      `README`).
    - **Guest Mode:** Пользователь получает быстрый доступ к функциям сайта через автоматическую генерацию аккаунта.

## 2. Страница популярных фильмов

- **Функционал:** Вывод самых популярных и рейтинговых фильмов.

## 3. Страница просмотра лайков

- **Функционал:** Отображение фильмов, которые пользователь оценил лайком.

## 4. Страница создания и просмотра дивана

- **Функционал:** Пользователь может:
    - Входить в общие комнаты.
    - Совместно выбирать фильмы.
    - Просматривать персонализированные рекомендации, адаптированные для всех пользователей в комнате.

## 5. Тиндер-режим

- **Функционал:** Пользователь может свайпами (лайк/дизлайк) настраивать персонализированную ленту в зависимости от
  своих интересов.

# Админ-панель

## Основные особенности

- **Доступ:** Для входа в админ-панель используйте логин `admin` и пароль `adminadmin`.
- **Адрес:** Панель доступна по пути `/adminPanel`.

## Основные функции

1. **Изменение фильмов:**
    - Поиск нужного фильма.
    - Редактирование данных фильма.

2. **Создание фильмов:**
    - Добавление новых фильмов в базу данных.

3. **Удаление фильмов:**
    - Поиск нужного фильма.
    - Удаление фильма из базы данных.

## Преимущества

- **Контроль данных:** Легко управляйте контентом и поддерживайте сайт в актуальном состоянии.
- **Развитие проекта:** Сосредоточьтесь на улучшении и развитии вашего проекта, не отвлекаясь на технические сложности.

ссылка: https://prod-team-32-n26k57br.REDACTED