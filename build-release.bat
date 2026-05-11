@echo off
REM Скрипт для быстрой сборки релизов на Windows

echo ================================
echo OLC VPN Client - Release Builder
echo ================================

REM Добавляем MinGW в PATH
set MINGW_BIN=C:\ProgramData\mingw64\mingw64\bin
if exist "%MINGW_BIN%\gcc.exe" (
    set PATH=%MINGW_BIN%;%PATH%
)

REM Проверяем Go
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Go не установлен
    echo Установите Go с https://go.dev/dl/
    pause
    exit /b 1
)

REM Проверяем gcc (MinGW)
where gcc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] gcc не найден - MinGW не установлен
    echo Запустите install-mingw-fixed.bat для установки MinGW
    pause
    exit /b 1
)

REM Создаём директорию release
if not exist release mkdir release
del /Q release\* 2>nul

echo.
echo [BUILD] Сборка Windows...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w -H windowsgui" -o release\olcvpn.exe .\cmd\olcvpn
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Ошибка сборки Windows
    pause
    exit /b 1
)

REM Создаём ZIP
powershell Compress-Archive -Path release\olcvpn.exe -DestinationPath release\olcvpn-windows.zip -Force
del release\olcvpn.exe

echo.
echo [INFO] Linux и macOS сборка пропущена
echo Для кросс-компиляции с GUI нужен Docker или нативная система
echo Используйте GitHub Actions или соберите на Linux/macOS напрямую

echo.
echo [BUILD] Сборка Android AAR...
where gomobile >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [WARN] gomobile не установлен, устанавливаю...
    go install golang.org/x/mobile/cmd/gomobile@latest
    gomobile init
)

if not exist android\app\libs mkdir android\app\libs
gomobile bind -target android -androidapi 21 -o android\app\libs\vpncore.aar .\mobile\
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Ошибка сборки Android AAR
    pause
    exit /b 1
)

echo.
echo [BUILD] Сборка Android APK...
if exist android\gradlew.bat (
    cd android
    call gradlew.bat assembleRelease
    if %ERRORLEVEL% EQU 0 (
        copy app\build\outputs\apk\release\app-release-unsigned.apk ..\release\olcvpn.apk
        echo [WARN] APK не подписан. Для подписи используйте jarsigner.
    ) else (
        echo [ERROR] Ошибка сборки APK
    )
    cd ..
) else (
    echo [WARN] Android Gradle не найден, пропускаю сборку APK
)

echo.
echo ================================
echo [SUCCESS] Сборка завершена!
echo ================================
echo.
echo Файлы в директории release\:
dir /B release
echo.
echo Готово к распространению!
pause
