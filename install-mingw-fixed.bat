@echo off
REM Упрощённая установка MinGW через winget или прямую ссылку

echo ========================================
echo Установка MinGW-w64 для OLC VPN Client
echo ========================================
echo.

REM Проверяем, не установлен ли уже
if exist "C:\mingw64\bin\gcc.exe" (
    echo [OK] MinGW уже установлен!
    goto check_gcc
)

REM Пробуем через winget (Windows 10/11)
where winget >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    echo [1/2] Устанавливаю через winget...
    winget install -e --id=MSYS2.MSYS2 --silent

    if exist "C:\msys64\mingw64\bin\gcc.exe" (
        set MINGW_PATH=C:\msys64\mingw64\bin
        goto add_to_path
    )
)

REM Если winget не сработал, показываем инструкцию
echo.
echo ========================================
echo Ручная установка (2 минуты):
echo ========================================
echo.
echo 1. Открой в браузере:
echo    https://github.com/niXman/mingw-builds-binaries/releases
echo.
echo 2. Скачай файл:
echo    x86_64-13.2.0-release-posix-seh-ucrt-rt_v11-rev0.7z
echo.
echo 3. Распакуй 7-Zip в C:\ (получится C:\mingw64)
echo.
echo 4. Запусти этот скрипт снова
echo.
echo ========================================
echo Или установи через Chocolatey:
echo    choco install mingw
echo ========================================
echo.
pause
exit /b 1

:add_to_path
echo [2/2] Добавляю в PATH...
set PATH=%MINGW_PATH%;%PATH%
setx PATH "%MINGW_PATH%;%PATH%" >nul 2>&1

:check_gcc
echo.
echo ========================================
echo Проверка:
echo ========================================

REM Ищем gcc
set GCC_PATH=
if exist "C:\mingw64\bin\gcc.exe" set GCC_PATH=C:\mingw64\bin
if exist "C:\msys64\mingw64\bin\gcc.exe" set GCC_PATH=C:\msys64\mingw64\bin

if "%GCC_PATH%"=="" (
    echo [ERROR] gcc не найден
    echo.
    echo Установи MinGW вручную (см. инструкцию выше)
    pause
    exit /b 1
)

REM Добавляем в PATH если нужно
set PATH=%GCC_PATH%;%PATH%
setx PATH "%GCC_PATH%;%PATH%" >nul 2>&1

where gcc
gcc --version

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo [SUCCESS] MinGW готов!
    echo ========================================
    echo.
    echo Теперь собери релизы:
    echo   build-release.bat
    echo.
) else (
    echo.
    echo [WARN] Перезапусти командную строку
    echo.
)

pause
