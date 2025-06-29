# DTree 🌳

A fast, interactive directory tree viewer for the terminal with smooth scrolling viewport. Navigate your filesystem with ease and open files instantly.

## ✨ Features

- **Smooth Scrolling Viewport** - Handles large directories with responsive scrolling
- **Interactive Navigation** - Use arrow keys or vim-style (j/k) to explore
- **Smart Lazy Loading** - Performance optimized for deep directory structures
- **File Integration** - Press Enter to open files with default applications
- **Configurable Depth** - Control initial expansion with `--depth` flag
- **Cross-Platform** - Works on macOS, Linux, and Windows
- **Clean Interface** - Minimal design focused on productivity

## 🚀 Quick Start

```bash
# View current directory
dtree

# View specific directory  
dtree /path/to/project

# Expand 3 levels deep
dtree --depth 3 .
```

## 📦 Installation

### 🚀 Quick Install (Recommended)
```bash
curl -sSL https://raw.githubusercontent.com/C0ldSmi1e/dtree/main/install.sh | bash
```

### 📥 Download Pre-built Binary
- [macOS (Intel)](https://github.com/C0ldSmi1e/dtree/releases/latest/download/dtree-darwin-amd64)
- [macOS (Apple Silicon)](https://github.com/C0ldSmi1e/dtree/releases/latest/download/dtree-darwin-arm64) 
- [Linux (x64)](https://github.com/C0ldSmi1e/dtree/releases/latest/download/dtree-linux-amd64)
- [Linux (ARM64)](https://github.com/C0ldSmi1e/dtree/releases/latest/download/dtree-linux-arm64)
- [Windows (x64)](https://github.com/C0ldSmi1e/dtree/releases/latest/download/dtree-windows-amd64.exe)

### 🛠️ For Go Developers
```bash
go install github.com/C0ldSmi1e/dtree@latest
```

### 🔧 Build from Source
```bash
git clone https://github.com/C0ldSmi1e/dtree.git
cd dtree
go build -o dtree
```

## 🎮 Controls

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate up/down |
| `Ctrl+U/D` | Jump half-screen up/down |
| `Ctrl+B/F` | Jump full-screen up/down |
| `gg/G` | Go to top/bottom |
| `Enter/Space` | Expand/collapse directories |
| `Enter` | Open files with default app |
| `q/Ctrl+C/Esc` | Quit |

## 🔧 Usage

```bash
dtree [options] [directory]

Options:
  -d, --depth <num>   Initial depth to expand (default: 1)
  -h, --help          Show help message

Examples:
  dtree               # Current directory, depth 1
  dtree ~/Projects    # Specific directory
  dtree -d 2 .        # Expand 2 levels deep
```

## 📋 Examples

### Basic Navigation
```
myproject/
▶ src/
▶ docs/
  README.md
  go.mod
```

### Expanded View
```
myproject/
▼ src/
  ├── main.go
  ├── handlers/
  └── utils/
▶ docs/
  README.md
  go.mod
```

## 🏗️ Development

### Requirements
- Go 1.21 or later

### Running Tests
```bash
go test ./tests/...           # All tests
go test ./tests/... -v        # Verbose output
go test ./tests/... -cover    # With coverage
```

### Project Structure
```
dtree/
├── main.go           # Entry point
├── internal/         # Private packages
│   ├── tree/        # Tree data structures
│   ├── ui/          # Terminal interface  
│   └── fileops/     # File operations
└── tests/           # Test suite
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) 🧋
- Styled with [Lip Gloss](https://github.com/charmbracelet/lipgloss) 💄
- Inspired by the classic `tree` command 🌲

---

**Made with ❤️ for developers who live in the terminal**