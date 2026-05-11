@echo off
REM Сборка всех платформ через Docker

echo ========================================
echo OLC VPN Client - Docker Cross-Compile
echo ========================================

REM Проверяем Docker
where docker >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Docker не установлен
    echo Установите Docker Desktop: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

REM Проверяем, запущен ли Docker
docker info >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Docker не запущен
    echo Запустите Docker Desktop и попробуйте снова
    pause
    exit /b 1
)

echo.
echo [1/4] Создание Docker образа для сборки...
docker build -f Dockerfile.builder -t olcvpn-builder --target base .
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Ошибка создания образа
    pause
    exit /b 1
)

REM Создаём директорию release
if not exist release mkdir release
del /Q release\* 2>nul

echo.
echo [2/4] Сборка Windows через Docker...
docker build -f Dockerfile.builder --target windows-builder -t olcvpn-windows .
docker create --name olcvpn-win-temp olcvpn-windows
docker cp olcvpn-win-temp:/out/olcvpn.exe release/
docker rm olcvpn-win-temp
if exist release\olcvpn.exe (
    powershell Compress-Archive -Path release\olcvpn.exe -DestinationPath release\olcvpn-windows.zip -Force
    del release\olcvpn.exe
    echo [OK] Windows build complete
) else (
    echo [ERROR] Windows build failed
)

echo.
echo [3/4] Сборка Linux через Docker...
docker build -f Dockerfile.builder --target linux-builder -t olcvpn-linux .
docker create --name olcvpn-linux-temp olcvpn-linux
docker cp olcvpn-linux-temp:/out/olcvpn release/
docker rm olcvpn-linux-temp
if exist release\olcvpn (
    tar -czf release\olcvpn-linux.tar.gz -C release olcvpn
    del release\olcvpn
    echo [OK] Linux build complete
) else (
    echo [ERROR] Linux build failed
)

echo.
echo [4/4] Сборка macOS через Docker...
echo [WARN] macOS сборка требует osxcross SDK
echo Для полноценной macOS сборки используйте GitHub Actions или Mac

echo.
echo ========================================
echo [SUCCESS] Сборка завершена!
echo ========================================
echo.
echo Файлы в release\:
dir /B release
echo.
pause
