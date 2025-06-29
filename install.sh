#!/bin/bash

# DTree Installer Script
# Usage: curl -sSL https://raw.githubusercontent.com/C0ldSmi1e/dtree/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values - prefer directories that don't need sudo
if [[ -d "/opt/homebrew/bin" && -w "/opt/homebrew/bin" ]]; then
    INSTALL_DIR="/opt/homebrew/bin"
elif [[ -d "/usr/local/bin" && -w "/usr/local/bin" ]]; then
    INSTALL_DIR="/usr/local/bin"
elif [[ -d "$HOME/.local/bin" ]]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$HOME/.local/bin"
else
    INSTALL_DIR="/usr/local/bin"
fi
REPO="C0ldSmi1e/dtree"

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case "$os" in
        darwin*)
            case "$arch" in
                x86_64) echo "darwin-amd64" ;;
                arm64) echo "darwin-arm64" ;;
                *) log_error "Unsupported architecture: $arch"; exit 1 ;;
            esac
            ;;
        linux*)
            case "$arch" in
                x86_64) echo "linux-amd64" ;;
                aarch64|arm64) echo "linux-arm64" ;;
                *) log_error "Unsupported architecture: $arch"; exit 1 ;;
            esac
            ;;
        mingw*|msys*|cygwin*)
            echo "windows-amd64.exe"
            ;;
        *)
            log_error "Unsupported OS: $os"
            exit 1
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# Download and install
install_dtree() {
    local platform=$(detect_platform)
    local version=$(get_latest_version)
    local binary_name="dtree-$platform"
    local download_url="https://github.com/$REPO/releases/download/$version/$binary_name"
    
    log_info "Detected platform: $platform"
    log_info "Latest version: $version"
    log_info "Download URL: $download_url"
    
    # Create temp directory
    local tmp_dir=$(mktemp -d)
    local tmp_file="$tmp_dir/dtree"
    
    # Download binary
    log_info "Downloading dtree..."
    if command -v curl >/dev/null 2>&1; then
        curl -sL "$download_url" -o "$tmp_file"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$download_url" -O "$tmp_file"
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    # Make executable
    chmod +x "$tmp_file"
    
    # Install to system
    if [[ -w "$INSTALL_DIR" ]]; then
        mv "$tmp_file" "$INSTALL_DIR/dtree"
        log_info "dtree installed to $INSTALL_DIR/dtree"
    else
        log_error "No write permission to $INSTALL_DIR"
        log_info "Please run one of these commands manually:"
        log_info "  sudo mv $tmp_file $INSTALL_DIR/dtree"
        log_info "  mv $tmp_file ~/bin/dtree  # (make sure ~/bin is in your PATH)"
        exit 1
    fi
    
    # Cleanup
    rm -rf "$tmp_dir"
    
    # Verify installation
    if command -v dtree >/dev/null 2>&1; then
        log_info "Installation successful! 🎉"
        log_info "Run 'dtree --help' to get started."
        dtree --help
    else
        log_error "Installation failed. Please check your PATH includes $INSTALL_DIR"
        exit 1
    fi
}

# Main
main() {
    log_info "Installing DTree - Interactive Directory Tree Viewer"
    log_info "Repository: https://github.com/$REPO"
    echo
    
    install_dtree
}

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi