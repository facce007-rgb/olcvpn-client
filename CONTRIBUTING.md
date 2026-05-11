# Contributing to OLC VPN Client

Спасибо за интерес к проекту! 🎉

## Как внести вклад

### Сообщить о баге

1. Проверьте, что баг еще не был сообщен в [Issues](https://github.com/yanisplugg/olcvpn-client/issues)
2. Создайте новый issue с подробным описанием:
   - Шаги для воспроизведения
   - Ожидаемое поведение
   - Фактическое поведение
   - Версия ОС и приложения
   - Логи (если есть)

### Предложить улучшение

1. Откройте issue с тегом `enhancement`
2. Опишите, что вы хотите улучшить и почему
3. Дождитесь обсуждения перед началом работы

### Pull Request процесс

1. **Fork** репозитория
2. Создайте **feature branch** от `develop`:
   ```bash
   git checkout develop
   git checkout -b feature/amazing-feature
   ```
3. Внесите изменения
4. Добавьте тесты (если применимо)
5. Убедитесь, что все тесты проходят:
   ```bash
   go test ./...
   ```
6. Commit с понятным сообщением:
   ```bash
   git commit -m "feat: add amazing feature"
   ```
7. Push в ваш fork:
   ```bash
   git push origin feature/amazing-feature
   ```
8. Откройте Pull Request в `develop` ветку

## Gitflow структура

```
main          ← стабильные релизы (только merge из release/hotfix)
  ↑
develop       ← активная разработка
  ↑
feature/*     ← новые фичи (merge в develop)
hotfix/*      ← срочные исправления (merge в main и develop)
release/*     ← подготовка релиза (merge в main и develop)
```

## Commit сообщения

Используем [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: добавить новую фичу
fix: исправить баг
docs: обновить документацию
style: форматирование кода
refactor: рефакторинг без изменения функциональности
perf: улучшение производительности
test: добавить тесты
build: изменения в сборке
ci: изменения в CI/CD
chore: прочие изменения
```

Примеры:
```
feat: add VLESS Reality support
fix: resolve SOCKS5 authentication bypass
docs: update installation guide for Windows
refactor: simplify profile storage logic
```

## Стиль кода

- Следуйте стандартам Go: `gofmt`, `golint`
- Комментарии на английском для кода, на русском для документации
- Тесты обязательны для новой функциональности
- Покрытие тестами должно быть > 70%

## Тестирование

```bash
# Запустить все тесты
go test ./...

# С покрытием
go test -cover ./...

# С race detector
go test -race ./...

# Конкретный пакет
go test ./internal/core/...
```

## Линтер

```bash
golangci-lint run ./...
```

## Вопросы?

Открывайте issue с тегом `question` или пишите в Discussions.

---

**Спасибо за вклад в проект!** 🚀
