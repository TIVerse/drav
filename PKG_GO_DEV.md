# Publishing DRAV to pkg.go.dev

## ✅ Completed Steps

1. ✅ Created version tag `v0.1.0`
2. ✅ Pushed tag to GitHub
3. ✅ Module indexed in Go proxy (proxy.golang.org)

## 🔄 Indexing on pkg.go.dev

Your module is now available in the Go module proxy! To complete indexing on pkg.go.dev:

### Method 1: Direct Browser Access (Recommended)
Visit this URL in your browser:
```
https://pkg.go.dev/github.com/TIVerse/drav@v0.1.0
```

This will trigger pkg.go.dev to fetch and index your module. It may take 2-3 minutes for the first indexing.

### Method 2: Go Get
Run this command from any directory:
```bash
go get github.com/TIVerse/drav@v0.1.0
```

### Method 3: Wait for Discovery
pkg.go.dev will automatically discover and index your module within a few hours.

## 📦 Module URLs

Once indexed, your module will be available at:

- **Main package**: https://pkg.go.dev/github.com/TIVerse/drav
- **Version v0.1.0**: https://pkg.go.dev/github.com/TIVerse/drav@v0.1.0
- **Dravya package**: https://pkg.go.dev/github.com/TIVerse/drav/pkg/dravya
- **Maya package**: https://pkg.go.dev/github.com/TIVerse/drav/pkg/maya
- **Prana package**: https://pkg.go.dev/github.com/TIVerse/drav/pkg/prana
- **Agni package**: https://pkg.go.dev/github.com/TIVerse/drav/pkg/agni
- **Vak package**: https://pkg.go.dev/github.com/TIVerse/drav/pkg/vak
- **Vayu package**: https://pkg.go.dev/github.com/TIVerse/drav/pkg/vayu
- **Sri package**: https://pkg.go.dev/github.com/TIVerse/drav/pkg/sri

## 🚀 Installing Your Module

Users can now install DRAV with:
```bash
go get github.com/TIVerse/drav@v0.1.0
```

Or using the latest version:
```bash
go get github.com/TIVerse/drav@latest
```

## 📝 Next Release

When ready for the next version:
```bash
# Create new tag
git tag v0.2.0 -m "Release message"

# Push tag
git push origin v0.2.0

# Index on pkg.go.dev (visit in browser)
https://pkg.go.dev/github.com/TIVerse/drav@v0.2.0
```

## ✨ Badge for README

Add this badge to your README.md:
```markdown
[![Go Reference](https://pkg.go.dev/badge/github.com/TIVerse/drav.svg)](https://pkg.go.dev/github.com/TIVerse/drav)
```

## 🔍 Verification

Check if your module is indexed:
```bash
# Check Go proxy
curl https://proxy.golang.org/github.com/\!t\!i\!verse/drav/@v/v0.1.0.info

# List module
go list -m github.com/TIVerse/drav@v0.1.0
```

Your module is **LIVE** on the Go proxy! 🎉
