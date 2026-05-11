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
	return sh.Run("go", "build", "-o", "olcvpn", "./cmd/olcvpn")
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

	// macOS - только на macOS хосте
	if runtime.GOOS == "darwin" {
		if err := buildMacOS(); err != nil {
			fmt.Println("⚠️  macOS build skipped (requires macOS host)")
		}
	} else {
		fmt.Println("⚠️  macOS build skipped (requires macOS host)")
	}

	// Linux - только на Linux хосте
	if runtime.GOOS == "linux" {
		if err := buildLinux(); err != nil {
			fmt.Println("⚠️  Linux build skipped (requires Linux host)")
		}
	} else {
		fmt.Println("⚠️  Linux build skipped (requires Linux host)")
	}

	// Android - требует gomobile
	fmt.Println("⚠️  Android build skipped (run: mage android)")

	// iOS - требует gomobile и macOS
	fmt.Println("⚠️  iOS build skipped (run: mage ios on macOS)")

	fmt.Println("\n✅ Windows release built successfully!")
	fmt.Println("📦 Check release/ directory")
	return nil
}

// Windows собирает только Windows релиз (для CI)
func Windows() error {
	fmt.Println("Building Windows release...")
	if err := os.MkdirAll("release", 0755); err != nil {
		return err
	}
	return buildWindows()
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
	if err := sh.RunWith(env, "go", "build", "-o", output, "./cmd/olcvpn"); err != nil {
		return err
	}

	// Копируем README и LICENSE
	sh.Copy(filepath.Join(tmpDir, "README.md"), "README.md")
	sh.Copy(filepath.Join(tmpDir, "LICENSE"), "LICENSE")

	// Создаём zip
	zipFile := "release/olcvpn-windows-amd64.zip"
	if runtime.GOOS == "windows" {
		return sh.Run("powershell", "Compress-Archive", "-Force", "-Path", tmpDir+"/*", "-DestinationPath", zipFile)
	}
	return sh.Run("zip", "-r", zipFile, tmpDir)
}

// MacOS собирает только macOS релиз (для CI)
func MacOS() error {
	fmt.Println("Building macOS release...")
	if err := os.MkdirAll("release", 0755); err != nil {
		return err
	}
	return buildMacOS()
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
	if err := sh.RunWith(env, "go", "build", "-o", outputIntel, "./cmd/olcvpn"); err != nil {
		return err
	}

	// Собираем для Apple Silicon
	env["GOARCH"] = "arm64"
	outputArm := filepath.Join(tmpDir, "olcvpn-arm64")
	if err := sh.RunWith(env, "go", "build", "-o", outputArm, "./cmd/olcvpn"); err != nil {
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

	// Копируем README и LICENSE
	sh.Copy(filepath.Join(tmpDir, "README.md"), "README.md")
	sh.Copy(filepath.Join(tmpDir, "LICENSE"), "LICENSE")

	// Создаём zip
	return sh.Run("zip", "-r", "release/olcvpn-macos-universal.zip", tmpDir)
}

// Linux собирает только Linux релиз (для CI)
func Linux() error {
	fmt.Println("Building Linux release...")
	if err := os.MkdirAll("release", 0755); err != nil {
		return err
	}
	return buildLinux()
}

// buildLinux собирает Linux релиз
func buildLinux() error {
	fmt.Println("\n📦 Building Linux...")

	tmpDir := "release/tmp-linux"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	env := map[string]string{
		"GOOS":        "linux",
		"GOARCH":      "amd64",
		"CGO_ENABLED": "1",
	}

	output := filepath.Join(tmpDir, "olcvpn")
	if err := sh.RunWith(env, "go", "build", "-o", output, "./cmd/olcvpn"); err != nil {
		return err
	}

	// Делаем исполняемым
	os.Chmod(output, 0755)

	// Копируем README и LICENSE
	sh.Copy(filepath.Join(tmpDir, "README.md"), "README.md")
	sh.Copy(filepath.Join(tmpDir, "LICENSE"), "LICENSE")

	// Создаём tar.gz
	return sh.Run("tar", "-czf", "release/olcvpn-linux-amd64.tar.gz", "-C", tmpDir, ".")
}

// buildAndroid собирает Android AAR
func buildAndroid() error {
	fmt.Println("\n📦 Building Android AAR...")

	// Проверяем gomobile
	if err := sh.Run("gomobile", "version"); err != nil {
		return fmt.Errorf("gomobile not found. Install: go install golang.org/x/mobile/cmd/gomobile@latest && gomobile init")
	}

	// Создаём директорию
	if err := os.MkdirAll("android/app/libs", 0755); err != nil {
		return err
	}

	// Собираем AAR
	return sh.Run("gomobile", "bind",
		"-target", "android",
		"-androidapi", "21",
		"-o", "android/app/libs/vpncore.aar",
		"./mobile/")
}

// buildIOS собирает iOS xcframework
func buildIOS() error {
	fmt.Println("\n📦 Building iOS xcframework...")

	// Проверяем gomobile
	if err := sh.Run("gomobile", "version"); err != nil {
		return fmt.Errorf("gomobile not found. Install: go install golang.org/x/mobile/cmd/gomobile@latest && gomobile init")
	}

	// Проверяем что мы на macOS
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("iOS build requires macOS")
	}

	// Создаём директорию
	if err := os.MkdirAll("ios/Frameworks", 0755); err != nil {
		return err
	}

	// Собираем xcframework
	return sh.Run("gomobile", "bind",
		"-target", "ios",
		"-o", "ios/Frameworks/VPNCore.xcframework",
		"./mobile/")
}

// Android собирает только Android AAR
func Android() error {
	return buildAndroid()
}

// IOS собирает только iOS xcframework
func IOS() error {
	return buildIOS()
}

// Run запускает приложение
func Run() error {
	fmt.Println("Running OLC VPN Client...")
	return sh.Run("go", "run", "./cmd/olcvpn")
}
