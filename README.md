# 🛡️ wpvul

![Go](https://img.shields.io/github/languages/top/spagu/go-wpvul)
![Build](https://img.shields.io/github/actions/workflow/status/spagu/go-wpvul/Build%20and%20Release?branch=main)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

**wpvul** is a lightning-fast, zero-dependency WordPress plugin scanner written in Go. It recursively scans your site directory (e.g., `wp-content/plugins`) and flags known malicious, duplicate, and redundant plugins based on a comprehensive, built-in blacklist.

---

## ✨ Features

- **🚀 Extremely Fast & Lightweight:** Compiled into a single embedded binary. No external dependencies, PHP engines or large setups needed.
- **🛡️ Embedded Blacklist:** Scans using a continuously updated internal index matching malicious entries (checks directories and standalone PHP/ZIP/TAR.GZ archives).
- **🌍 i18n Support:** Fully translated into **English** (default) and **Polish**! Handled automatically via `LANG` env variable.
- **💻 Cross-Platform:** Available and natively built for **Linux**, **macOS** (Darwin), and **FreeBSD** across both `amd64` and `arm64` architectures.

---

## 📦 Installation

To quickly get up and running, you can install the binary pointing to GitHub Releases using our handy installer.

### Quick cURL Installer (Linux / macOS / FreeBSD)
```bash
curl -sSL https://raw.githubusercontent.com/spagu/go-wpvul/main/install.sh | bash
```

### Homebrew / macOS Local Make
You can use `make brew` to copy the natively generated binary (if compiled from source via Makefile) to `/usr/local/bin` seamlessly on Mac.
```bash
make brew
```

### Build from Source
Ensure you have **Go 1.21+** installed on your system.
```bash
git clone https://github.com/spagu/go-wpvul.git
cd go-wpvul

# Compile binaries for ALL OS architectures (Linux, Mac, FreeBSD; arm + amd)
make compile-all
```
Your compiled binaries for different platforms will be waiting for you inside the `./build/` directory!

---

## 🚀 Usage

Execute `wpvul` by passing the absolute or relative directory you want to scan:

```bash
# General Usage
wpvul /var/www/html/wp-content/plugins

# Example Output
# [DETECTED] /var/www/html/wp-content/plugins/bad-behavior
#  ├─ matched slug: bad-behavior
#  └─ blacklist source: User list
```

The script returns an explicit exit code `!= 0` if at least 1 vulnerability/banned plugin was identified, making it perfect for CI/CD integrations.

---

## 🌐 Language Settings (i18n)

The app automatically adjusts its interface to your operating system's environment variables:
- `LANG=pl_PL.UTF-8 wpvul ...` outputs alerts in **Polish**.
- `LANG=en_US.UTF-8 wpvul ...` outputs in **English**.

All messages are handled smoothly in-memory without extra JSON payloads.

---

## 🤝 Contributing & Blacklist Upgrades

The tool checks against `cv-banned.csv`. Because Go embeds the CSV at compile time using `//go:embed`, you can update the script by modifying the CSV file and hitting:
```bash
make build
```

## License
MIT License
