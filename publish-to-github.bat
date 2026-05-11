@echo off
REM Quick publish script for GitHub (Windows)

echo 🚀 OLC VPN Client - GitHub Publish Script
echo.

REM Check if git is initialized
if not exist .git (
    echo ❌ Not a git repository!
    exit /b 1
)

REM Check if remote exists
git remote | findstr /C:"origin" >nul
if %errorlevel% equ 0 (
    echo ✅ Remote 'origin' already exists
    git remote -v
) else (
    echo 📝 Adding remote 'origin'...
    set /p REPO_URL="Enter GitHub repository URL: "
    git remote add origin %REPO_URL%
)

echo.
echo 📤 Pushing branches...

REM Push main
echo   → Pushing main...
git push -u origin main

REM Push develop
echo   → Pushing develop...
git push -u origin develop

REM Push all branches
echo   → Pushing all branches...
git push --all origin

REM Push tags
echo   → Pushing tags...
git push --tags origin

echo.
echo ✅ Done! Repository published to GitHub
echo.
echo Next steps:
echo 1. Go to https://github.com/yanisplugg/olcvpn-client
echo 2. Configure branch protection (Settings → Branches)
echo 3. Add repository description and topics
echo 4. Check GitHub Actions (Actions tab)
echo 5. Create first release by pushing a tag: git tag v1.0.0 ^&^& git push origin v1.0.0
echo.
echo 📚 See GITHUB_PUBLISH.md for detailed instructions
