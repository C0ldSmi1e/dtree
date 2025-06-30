# DTree 🌳

The `tree` command, but interactive - a terminal filesystem navigator.

## ✨ Features

- **Interactive Navigation** - Navigate with arrow keys or vim-style controls (j/k/gg/G)
- **Instant File Opening** - Press Enter to open with default applications
- **Configurable Depth** - See as much or as little as you want
- **Cross-Platform** - Works on macOS, Linux and WSL
- **Zero Dependencies** - Single binary, no installation complexity

## 🚀 Quick Start

```bash
# View current directory
dtree

# View specific directory  
dtree /path/to/project

# Expand 3 levels deep
dtree --depth 3 .

# Usage
dtree -h
```

## 📦 Installation

### 🚀 Quick Install (Recommended)
```bash
curl -sSL https://raw.githubusercontent.com/C0ldSmi1e/dtree/main/install.sh -o dtree_install.sh && bash dtree_install.sh && rm dtree_install.sh
```

### 📥 Download Pre-built Binary
- [macOS (Intel)](https://github.com/C0ldSmi1e/dtree/releases/latest/download/dtree-darwin-amd64)
- [macOS (Apple Silicon)](https://github.com/C0ldSmi1e/dtree/releases/latest/download/dtree-darwin-arm64) 
- [Linux (x64)](https://github.com/C0ldSmi1e/dtree/releases/latest/download/dtree-linux-amd64)
- [Linux (ARM64)](https://github.com/C0ldSmi1e/dtree/releases/latest/download/dtree-linux-arm64)

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