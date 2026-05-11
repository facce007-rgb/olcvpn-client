# 🚀 Инструкция по публикации OLC VPN Client

## Шаг 1: Инициализация Git (выполните в командной строке)

```bash
cd C:\freeclaude\olcvpn-clean

# Инициализация репозитория
git init

# Добавление всех файлов
git add .

# Первый коммит
git commit -m "Initial commit: OLC VPN Client v1.0.0"
```

## Шаг 2: Создание репозитория на GitHub

1. Откройте https://github.com/new
2. Название: `olcvpn-client` (или любое другое)
3. Описание: `Cross-platform VPN client with sing-box and olcrtc support`
4. Выберите **Public** (чтобы пользователи могли скачивать)
5. НЕ добавляйте README, .gitignore, license (уже есть в проекте)
6. Нажмите **Create repository**

## Шаг 3: Загрузка на GitHub

GitHub покажет команды, выполните их:

```bash
git remote add origin https://github.com/ВАШ_USERNAME/olcvpn-client.git
git branch -M main
git push -u origin main
```

## Шаг 4: Создание первого релиза

```bash
# Создание тега версии
git tag v1.0.0

# Отправка тега на GitHub
git push origin v1.0.0
```

## Шаг 5: Автоматическая сборка

После `git push origin v1.0.0`:

1. GitHub Actions автоматически запустится
2. Соберёт Windows, Linux, macOS, Android (займёт ~10-15 минут)
3. Создаст Release с готовыми файлами

## Шаг 6: Пользователи скачивают

Откройте: `https://github.com/ВАШ_USERNAME/olcvpn-client/releases`

Пользователи увидят:
- ✅ **olcvpn-windows.zip** - для Windows (распаковать и запустить)
- ✅ **olcvpn-linux.tar.gz** - для Linux
- ✅ **olcvpn-macos.zip** - для macOS
- ✅ **olcvpn.apk** - для Android (скачать и установить)

## Обновление версии (в будущем)

```bash
# Внесите изменения в код
git add .
git commit -m "Описание изменений"
git push

# Создайте новый релиз
git tag v1.0.1
git push origin v1.0.1
```

GitHub Actions автоматически соберёт новую версию!

---

## Альтернатива: Локальная сборка (если нет GitHub)

Если не хотите использовать GitHub:

1. Установите Docker Desktop
2. Запустите `build-all-docker.bat`
3. Файлы появятся в папке `release\`
4. Загрузите их на любой файлообменник (Google Drive, Dropbox, etc.)

---

**Рекомендация:** Используйте GitHub - это бесплатно, автоматически и профессионально.
