English | [简体中文](./README_zh.md)

# AutoSSH (Timo Edition)

A powerful interactive SSH remote client designed for macOS/Linux terminals. It helps you save complex SSH passwords and key paths, while integrating highly efficient fuzzy search and file transfer capabilities.

This project is a fork and enhancement of [Lenbo/autossh](https://github.com/islenbo/autossh).

![Demo](https://raw.githubusercontent.com/yaogh99123/autossh/main/docs/images/ezgif-5-42b5117192fc.gif)

## 🚀 Key Enhancements (vs. Original)

- **Source-level FZF Integration**: Built-in official `fzf` fuzzy search engine, **no external fzf binary required**. Supports non-fullscreen popups (40% height) to maintain context during search.
- **YAML Configuration Support**: Full support for `config.yml` format for clearer configuration and robust I/O.
- **Smart SSH Key Association**: When adding a server, it automatically scans the `~/.ssh` directory and provides a fuzzy selection interface to bind private keys instantly—no more manual path typing.
- **High-Performance Parallel Transfer**:
    - Interactive mode supports `up`/`down` commands.
    - Supports `-j` parameter for multi-threaded concurrency, significantly boosting transfer speeds for large batches of small files.
    - Full support for `-r` recursive directory transfer.
- **UI/UX Optimization**:
    - **List Folding**: Displays only the first 8 servers by default. Toggle with `a` (All) or `h` (Hide) commands to keep the UI clean.
    - **Index Retrieval**: Search box supports direct entry of list indices (e.g., `a1`, `b2`) for precise positioning.
    - **Signal Handling**: Enhanced global `Ctrl+C` capture ensures a graceful exit in any state.

## 🛠 Feature Guide

### 1. Quick SSH Login
- Run `./autossh` to enter interactive mode.
- Enter the **Index** (e.g., `1` or `a1`) or **Alias** to login directly.
- Enter **`s`** to activate fuzzy search mode.

### 2. File Upload & Download (SFTP-based)
Support direct command entry in the main interface:
- **Upload**: `up [-r] [-j 10] <local_path> <server:remote_path>`
- **Download**: `down [-r] [-j 10] <server:remote_path> <local_path>`
- **Example**: `up -r -j 10 ./dist server1:/var/www/html`

### 3. Server Management
- **Add**: Enter `add`.
- **Edit**: Enter `edit`.
- **Remove**: Enter `remove`.

### 4. Auto Update
- Enter `upgrade` to automatically detect and upgrade to the latest version from `yaogh99123/autossh`.

## 📦 Installation & Build

### Direct Run
Mac/Linux users can download the Release package and run the `install` script.

### Manual Build
This project uses `vendor` mode for dependency management, ensuring source-level builds:
```bash
# Generate binary packages for all platforms
./build.sh

# Build for current platform only
make build
```

## 📄 Acknowledgments
Thanks to the original author [islenbo](https://github.com/islenbo) for the excellent base architecture.

## ⚖️ License
MIT
