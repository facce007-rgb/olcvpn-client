@echo off
setlocal

:: ============================================================
::  Claude Code через OmniRoute
::  Настрой переменные ниже, потом запускай этот батник
:: ============================================================

:: API ключ из Dashboard -> Endpoints в OmniRoute
set OMNIROUTE_KEY=sk-d3a5199541fb0342-e39c42-bd5ca237

:: Адрес OmniRoute (по умолчанию локальный)
set OMNIROUTE_URL=http://localhost:20128/v1

:: Модель (примеры: cc/claude-sonnet-4-5, cc/claude-opus-4-7, gg/gemini-2.5-pro)
:: Оставь пустым — OmniRoute выберет по combo/fallback цепочке
set OMNIROUTE_MODEL=kr/claude-sonnet-4.5

:: ============================================================

:: Проверка: заполнен ли ключ
if "%OMNIROUTE_KEY%"==sk-d3a5199541fb0342-e39c42-bd5ca237 (
    echo [ОШИБКА] Укажи OMNIROUTE_KEY в этом батнике!
    echo Ключ берётся из OmniRoute Dashboard ^> Endpoints
    pause
    exit /b 1
)

:: Проверка: запущен ли OmniRoute
curl -s --max-time 3 "%OMNIROUTE_URL%/models" -H "Authorization: Bearer %OMNIROUTE_KEY%" >nul 2>&1
if errorlevel 1 (
    echo [ПРЕДУПРЕЖДЕНИЕ] OmniRoute недоступен по адресу %OMNIROUTE_URL%
    echo Убедись что omniroute запущен, затем нажми любую клавишу...
    pause >nul
)

:: Устанавливаем переменные для Claude Code
set ANTHROPIC_BASE_URL=%OMNIROUTE_URL%
set ANTHROPIC_AUTH_TOKEN=%OMNIROUTE_KEY%

set CLAUDE_EXE=C:\Users\ian\AppData\Roaming\npm\claude.cmd

if not "%OMNIROUTE_MODEL%"=="" (
    echo [INFO] Запуск claude --model %OMNIROUTE_MODEL%
    echo [INFO] Base URL: %OMNIROUTE_URL%
    echo.
    "%CLAUDE_EXE%" --model %OMNIROUTE_MODEL% %*
) else (
    echo [INFO] Запуск claude (модель по combo/fallback)
    echo [INFO] Base URL: %OMNIROUTE_URL%
    echo.
    "%CLAUDE_EXE%" %*
)

endlocal
