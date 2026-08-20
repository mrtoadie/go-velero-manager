# go-velero-manager

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white) ![GitHub License](https://img.shields.io/github/license/mrtoadie/go-velero-manager) ![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/mrtoadie/go-velero-manager/total) ![GitHub Release](https://img.shields.io/github/v/release/mrtoadie/go-velero-manager) ![GitHub Repo stars](https://img.shields.io/github/stars/mrtoadie/go-velero-manager?style=flat&label=Stars) 
Interactive TUI tool for managing Velero backups in Kubernetes. Built with Go and Charm Huh.

## Features

- List all Velero backups with status, errors, warnings
- Multi-select backups for deletion
- Color-coded status: ✓ Completed | ✗ Failed | ⚠️ Partial
- Confirmation prompt before irreversible deletion
- Sorted by creation date (newest first)

## Requirements

- Go 1.19+
- `velero` CLI installed and configured
- Kubernetes cluster access

## Install
```bash 
git clone https://github.com/toadie/velero-manager.git 
cd velero-manager go build -o velero-manager 
./velero-manager
```

Or run directly:

```bash
go run main.go
```

---
:white_check_mark: Tested and developed on [Arch Linux](https://archlinux.org/) and a [K3S](https://k3s.io/) cluster

## License
go-check-cert is licensed under the [MIT License](LICENSE).