# 🚀 БЫСТРАЯ УСТАНОВКА MinGW

## Проблема
Файл `.7z` не распаковывается через PowerShell.

## ✅ Решение (выбери самое простое):

### Вариант 1: Через Chocolatey (30 секунд)
Если у тебя установлен Chocolatey:
```bash
choco install mingw
```

Затем:
```bash
build-release.bat
```

---

### Вариант 2: Через winget (1 минута)
Windows 10/11 встроенный:
```bash
winget install MSYS2.MSYS2
```

Затем в открывшемся окне MSYS2:
```bash
pacman -S mingw-w64-x86_64-gcc
```

Добавь в PATH: `C:\msys64\mingw64\bin`

---

### Вариант 3: Прямая загрузка (2 минуты)

1. **Открой:** https://github.com/niXman/mingw-builds-binaries/releases

2. **Скачай:** `x86_64-13.2.0-release-posix-seh-ucrt-rt_v11-rev0.7z` (50 MB)

3. **Распакуй 7-Zip** в `C:\` (получится `C:\mingw64`)

4. **Добавь в PATH:**
   - Система → Дополнительные параметры → Переменные среды
   - PATH → Добавить: `C:\mingw64\bin`

5. **Проверь:**
   ```bash
   gcc --version
   ```

---

### Вариант 4: Готовый установщик (самый простой)

1. **Скачай:** https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/installer/mingw-w64-install.exe

2. **Запусти** и выбери:
   - Version: 13.2.0
   - Architecture: x86_64
   - Threads: posix
   - Exception: seh

3. **Установи** в `C:\mingw64`

4. **Добавь в PATH:** `C:\mingw64\bin`

---

## После установки:

```bash
gcc --version  # Проверка
build-release.bat  # Сборка всех релизов
```

**Результат в папке `release/`!** 🎉
