//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/sh"
)

// Build собирает desktop приложение для текущей платформы
func Build() error {
	fmt.Println("Building OLC VPN Client...")
	return sh.Run("go", "build", "-o", "olcvpn", "./cmd/olcvpn/main_v2.go")
}

// Test запускает все тесты
func Test() error {
	fmt.Println("Running tests...")
	return sh.Run("go", "test", "./...")
}

// Lint запускает golangci-lint
func Lint() error {
	fmt.Println("Running linter...")
	return sh.Run("golangci-lint", "run", "./...")
}

// Clean очищает артефакты сборки
func Clean() error {
	fmt.Println("Cleaning...")
	os.RemoveAll("build/")
	os.RemoveAll("release/")
	return nil
}

// Release собирает релизы для всех платформ в папку release/
func Release() error {
	fmt.Println("Building releases for all platforms...")

	if err := Clean(); err != nil {
		return err
	}

	// Создаём директории
	if err := os.MkdirAll("release", 0755); err != nil {
		return err
	}

	// Windows
	if err := buildWindows(); err != nil {
		return err
	}

	// macOS
	if err := buildMacOS(); err != nil {
		return err
	}

	// Linux
	if err := buildLinux(); err != nil {
		return err
	}

	// Android
	if err := buildAndroid(); err != nil {
		return err
	}

	// iOS
	if err := buildIOS(); err != nil {
		return err
	}

	fmt.Println("\n✅ All releases built successfully!")
	fmt.Println("📦 Check release/ directory")
	return nil
}

// buildWindows собирает Windows релиз
func buildWindows() error {
	fmt.Println("\n📦 Building Windows...")

	// Создаём временную директорию
	tmpDir := "release/tmp-windows"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Собираем exe
	env := map[string]string{
		"GOOS":        "windows",
		"GOARCH":      "amd64",
		"CGO_ENABLED": "1",
	}

	output := filepath.Join(tmpDir, "olcvpn.exe")
	if err := sh.RunWith(env, "go", "build", "-o", output, "./cmd/olcvpn/main_v2.go"); err != nil {
		return err
	}

	// Создаём zip
	zipFile := "release/olcvpn-windows.zip"
	if runtime.GOOS == "windows" {
		return sh.Run("powershell", "Compress-Archive", "-Path", tmpDir+"/*", "-DestinationPath", zipFile)
	}
	return sh.Run("zip", "-r", zipFile, tmpDir)
}

// buildMacOS собирает macOS релиз
func buildMacOS() error {
	fmt.Println("\n📦 Building macOS...")

	tmpDir := "release/tmp-macos"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Собираем для Intel
	env := map[string]string{
		"GOOS":        "darwin",
		"GOARCH":      "amd64",
		"CGO_ENABLED": "1",
	}
	outputIntel := filepath.Join(tmpDir, "olcvpn-amd64")
	if err := sh.RunWith(env, "go", "build", "-o", outputIntel, "./cmd/olcvpn/main_v2.go"); err != nil {
		return err
	}

	// Собираем для Apple Silicon
	env["GOARCH"] = "arm64"
	outputArm := filepath.Join(tmpDir, "olcvpn-arm64")
	if err := sh.RunWith(env, "go", "build", "-o", outputArm, "./cmd/olcvpn/main_v2.go"); err != nil {
		return err
	}

	// Создаём universal binary
	outputUniversal := filepath.Join(tmpDir, "olcvpn")
	if err := sh.Run("lipo", "-create", "-output", outputUniversal, outputIntel, outputArm); err != nil {
		// Если lipo недоступен, просто копируем один из бинарников
		if err := sh.Run("cp", outputArm, outputUniversal); err != nil {
			return err
		}
	}

	// Делаем исполняемым
	os.Chmod(outputUniversal, 0755)

	// Создаём zip
	zipFile := "release/olcvpn-macos.zip"
	return sh.Run("zip", "-r", zipFile, tmpDir)
}

// buildLinux собирает Linux релиз
func buildLinux() error {
	fmt.Println("\n📦 Building Linux...")

	tmpDir := "release/tmp-linux"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Собираем бинарник
	env := map[string]string{
		"GOOS":        "linux",
		"GOARCH":      "amd64",
		"CGO_ENABLED": "1",
	}

	output := filepath.Join(tmpDir, "olcvpn")
	if err := sh.RunWith(env, "go", "build", "-o", output, "./cmd/olcvpn/main_v2.go"); err != nil {
		return err
	}

	// Делаем исполняемым
	os.Chmod(output, 0755)

	// Создаём tar.gz
	tarFile := "release/olcvpn-linux.tar.gz"
	return sh.Run("tar", "-czf", tarFile, "-C", tmpDir, ".")
}

// buildAndroid собирает Android APK
func buildAndroid() error {
	fmt.Println("\n📦 Building Android...")

	// Сначала собираем AAR
	if err := os.MkdirAll("android/app/libs", 0755); err != nil {
		return err
	}

	if err := sh.Run("gomobile", "bind",
		"-target", "android",
		"-androidapi", "21",
		"-o", "android/app/libs/vpncore.aar",
		"./mobile/"); err != nil {
		return err
	}

	// Собираем APK через Gradle
	if err := os.Chdir("android"); err != nil {
		return err
	}
	defer os.Chdir("..")

	var gradleCmd string
	if runtime.GOOS == "windows" {
		gradleCmd = "gradlew.bat"
	} else {
		gradleCmd = "./gradlew"
	}

	if err := sh.Run(gradleCmd, "assembleRelease"); err != nil {
		return err
	}

	// Копируем APK в release
	apkSrc := "app/build/outputs/apk/release/app-release.apk"
	apkDst := "../release/olcvpn.apk"
	return sh.Run("cp", apkSrc, apkDst)
}

// buildIOS собирает iOS xcframework
func buildIOS() error {
	fmt.Println("\n📦 Building iOS...")

	if err := os.MkdirAll("ios/Frameworks", 0755); err != nil {
		return err
	}

	if err := sh.Run("gomobile", "bind",
		"-target", "ios",
		"-o", "ios/Frameworks/VPNCore.xcframework",
		"./mobile/"); err != nil {
		return err
	}

	// Создаём zip с xcframework
	zipFile := "release/olcvpn-ios-framework.zip"
	return sh.Run("zip", "-r", zipFile, "ios/Frameworks/VPNCore.xcframework")
}

// Dev запускает приложение в режиме разработки
func Dev() error {
	fmt.Println("Running in development mode...")
	return sh.Run("go", "run", "./cmd/olcvpn/main_v2.go")
}

// Deps устанавливает зависимости
func Deps() error {
	fmt.Println("Installing dependencies...")
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "mod", "tidy")
}

// Android собирает только Android AAR
func Android() error {
	fmt.Println("Building Android AAR...")

	if err := os.MkdirAll("android/app/libs", 0755); err != nil {
		return err
	}

	return sh.Run("gomobile", "bind",
		"-target", "android",
		"-androidapi", "21",
		"-o", "android/app/libs/vpncore.aar",
		"./mobile/")
}

// IOS собирает только iOS xcframework
func IOS() error {
	fmt.Println("Building iOS xcframework...")

	if err := os.MkdirAll("ios/Frameworks", 0755); err != nil {
		return err
	}

	return sh.Run("gomobile", "bind",
		"-target", "ios",
		"-o", "ios/Frameworks/VPNCore.xcframework",
		"./mobile/")
}
