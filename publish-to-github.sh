#!/bin/bash
# Quick publish script for GitHub

echo "🚀 OLC VPN Client - GitHub Publish Script"
echo ""

# Check if git is initialized
if [ ! -d .git ]; then
    echo "❌ Not a git repository!"
    exit 1
fi

# Check if remote exists
if git remote | grep -q origin; then
    echo "✅ Remote 'origin' already exists"
    git remote -v
else
    echo "📝 Adding remote 'origin'..."
    read -p "Enter GitHub repository URL: " REPO_URL
    git remote add origin "$REPO_URL"
fi

echo ""
echo "📤 Pushing branches..."

# Push main
echo "  → Pushing main..."
git push -u origin main

# Push develop
echo "  → Pushing develop..."
git push -u origin develop

# Push all branches
echo "  → Pushing all branches..."
git push --all origin

# Push tags
echo "  → Pushing tags..."
git push --tags origin

echo ""
echo "✅ Done! Repository published to GitHub"
echo ""
echo "Next steps:"
echo "1. Go to https://github.com/yanisplugg/olcvpn-client"
echo "2. Configure branch protection (Settings → Branches)"
echo "3. Add repository description and topics"
echo "4. Check GitHub Actions (Actions tab)"
echo "5. Create first release by pushing a tag: git tag v1.0.0 && git push origin v1.0.0"
echo ""
echo "📚 See GITHUB_PUBLISH.md for detailed instructions"
