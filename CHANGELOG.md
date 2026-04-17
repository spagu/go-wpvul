# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- `--exclude` (`-e`) flag to selectively omit specific directories from the filesystem scan (e.g., skips traversing `dev.hide`, `node_modules` caches or specific backup repositories). 

### Fixed
- Addressed an issue where deeply embedded packages inside legitimate plugins raised false positive matches by ensuring matched entities appear explicitly at the standard `/plugins/` or `/mu-plugins/` folder level.

## [1.0.1] - 2026-04-17

### Added
- Integrated **Cobra CLI** framework for much more coherent positional argument parsing and automatically generated shell help outputs `--help`.
- Beautiful **UTF-8 formatted terminal output** incorporating robust ANSI string manipulations to format paths, matched slugs, and sources efficiently.
- New flag `--bw` natively included to completely suppress colored and explicit unicode blocks for shell outputs mapped to external data pipelines (clean ASCII representations).

### Fixed
- Fixed an issue causing the scanner to recursively inspect subdirectories belonging to legitimately operating, recognized WordPress plugins, greatly preventing submodule mismatching.

## [1.0.0] - 2026-04-17

### Added
- **Complete Go Implementation**: Transition to a cleanly compiled binary system entirely independent of any external prerequisites (no PHP processors required).
- Initial static compilation embedding the `cv-banned.csv` table using the standard `//go:embed` feature for blazingly fast dictionary allocations.
- Comprehensive CI/CD via GitHub Actions triggering native executable builds across: Linux, FreeBSD, and MacOS for both standard `amd64` and `arm64` chips.
- Seamless auto-detecting i18n support based on ambient `LANG` environmental strings (en, pl).
- Fully independent recursive `.php`, `.zip` and `.tar.gz` analyzer paired with `install.sh` fetch script for fast server upgrades.
