# 🚀 Публикация на GitHub

## Шаг 1: Создать репозиторий на GitHub

1. Перейти на https://github.com/new
2. Название: `olcvpn-client`
3. Описание: `Cross-platform VPN client with sing-box and olcrtc support`
4. Visibility: **Public**
5. **НЕ** инициализировать с README, .gitignore или LICENSE (уже есть локально)
6. Нажать **Create repository**

## Шаг 2: Подключить локальный репозиторий

```bash
# Добавить remote
git remote add origin https://github.com/yanisplugg/olcvpn-client.git

# Проверить remote
git remote -v

# Push main ветки
git push -u origin main

# Push develop ветки
git push -u origin develop

# Push всех веток
git push --all origin

# Push тегов
git push --tags origin
```

## Шаг 3: Настроить GitHub репозиторий

### Описание репозитория
```
Cross-platform VPN client with sing-box and olcrtc engines. Windows, Linux, macOS, Android, iOS support.
```

### Topics (теги)
```
vpn, vpn-client, sing-box, olcrtc, vless, reality, shadowsocks, trojan, 
webrtc, golang, fyne, android, ios, cross-platform, wireguard-alternative
```

### Settings → General
- ✅ Issues
- ✅ Discussions (опционально)
- ✅ Projects (опционально)
- ✅ Wiki (опционально)

### Settings → Branches
Настроить защиту веток:

**main branch:**
- ✅ Require a pull request before merging
- ✅ Require status checks to pass before merging
  - CI / test
  - CI / lint
  - CI / build-test
- ✅ Require conversation resolution before merging
- ✅ Do not allow bypassing the above settings

**develop branch:**
- ✅ Require a pull request before merging
- ✅ Require status checks to pass before merging

## Шаг 4: Создать первый релиз

### Автоматический (через GitHub Actions)
```bash
# Создать и push тег
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions автоматически:
# 1. Соберёт все платформы
# 2. Создаст Release
# 3. Загрузит артефакты
```

### Ручной (если нужно)
1. Перейти в **Releases** → **Draft a new release**
2. Tag: `v1.0.0`
3. Title: `v1.0.0 - Initial Release`
4. Description: скопировать из `.github/workflows/release.yml`
5. Загрузить артефакты из `release/` папки
6. Нажать **Publish release**

## Шаг 5: Проверить CI/CD

После push проверить:
- ✅ GitHub Actions запустились
- ✅ CI тесты прошли
- ✅ Release workflow готов к запуску

## Шаг 6: Обновить README

После создания релиза обновить ссылки в README:
```bash
git checkout develop
# Отредактировать README.md (ссылки уже правильные)
git add README.md
git commit -m "docs: update release links"
git push origin develop
```

## Структура веток

```
main                    ← стабильные релизы (защищена)
  ↓
  v1.0.0 (tag)
  
develop                 ← активная разработка (защищена)
  ↓
  feature/*             ← новые фичи
  └─ feature/subscription-support
  
hotfix/*                ← срочные исправления
  └─ hotfix/critical-security-fix
```

## Workflow для разработки

### Новая фича
```bash
git checkout develop
git pull origin develop
git checkout -b feature/my-feature
# ... работа ...
git add .
git commit -m "feat: add my feature"
git push origin feature/my-feature
# Создать PR в develop на GitHub
```

### Hotfix
```bash
git checkout main
git pull origin main
git checkout -b hotfix/fix-name
# ... исправление ...
git add .
git commit -m "fix: critical bug"
git push origin hotfix/fix-name
# Создать PR в main И develop на GitHub
```

### Релиз
```bash
git checkout main
git merge develop
git tag v1.1.0 -m "Release v1.1.0"
git push origin main --tags
# GitHub Actions автоматически создаст релиз
```

## Секреты для Android подписи (опционально)

Если хотите подписывать APK:

1. Settings → Secrets and variables → Actions
2. Добавить секреты:
   - `KEYSTORE_FILE` - base64 keystore.jks
   - `KEYSTORE_PASSWORD` - пароль keystore
   - `KEY_ALIAS` - алиас ключа
   - `KEY_PASSWORD` - пароль ключа

Создать keystore:
```bash
keytool -genkey -v -keystore keystore.jks -keyalg RSA -keysize 2048 -validity 10000 -alias olcvpn
base64 keystore.jks > keystore.b64
# Содержимое keystore.b64 → KEYSTORE_FILE secret
```

## Проверка перед публикацией

- [ ] Все тесты проходят локально: `go test ./...`
- [ ] Линтер проходит: `golangci-lint run ./...`
- [ ] README актуален
- [ ] CLAUDE.md актуален
- [ ] LICENSE корректен
- [ ] .gitignore настроен
- [ ] GitHub Actions workflows корректны
- [ ] Версия в go.mod актуальна

## После публикации

1. Добавить бейдж в README:
```markdown
[![Release](https://img.shields.io/github/v/release/yanisplugg/olcvpn-client)](https://github.com/yanisplugg/olcvpn-client/releases)
[![CI](https://github.com/yanisplugg/olcvpn-client/workflows/CI/badge.svg)](https://github.com/yanisplugg/olcvpn-client/actions)
[![License](https://img.shields.io/github/license/yanisplugg/olcvpn-client)](LICENSE)
```

2. Создать Discussions для сообщества
3. Настроить GitHub Pages для документации (опционально)
4. Добавить в awesome-vpn списки

## Полезные команды

```bash
# Проверить статус
git status

# Посмотреть историю
git log --oneline --graph --all

# Посмотреть теги
git tag -l

# Удалить локальную ветку
git branch -d feature/old-feature

# Удалить remote ветку
git push origin --delete feature/old-feature

# Синхронизировать с remote
git fetch --all --prune
```

---

**Готово к публикации!** 🎉
