<p align="center">
  <a href="https://github.com/mrtoadie/go-velero-manager/releases"><img src="https://img.shields.io/github/v/release/mrtoadie/go-velero-manager?style=flat&label=Release" alt="Latest Release"></a>
  <a href="https://github.com/mrtoadie/go-velero-manager/releases"><img src="https://img.shields.io/github/downloads/mrtoadie/go-velero-manager/total?style=flat&label=Downloads" alt="Total Downloads"></a>
  <img src="https://img.shields.io/github/stars/mrtoadie/go-velero-manager?style=flat&label=Stars" alt="Stars">
  <img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white" alt="Made with Go">
  <img src="https://img.shields.io/github/license/mrtoadie/go-velero-manager?style=flat&label=License" alt="License">
</p>

<h3 align="center">Interactive TUI for managing Velero backups — list, inspect, and clean up without wrangling CLI flags.</h3>

---

**go-velero-manager** is a terminal UI tool for managing [Velero](https://velero.io/) backups in [Kubernetes](https://kubernetes.io/). Built with Go and [huh](https://github.com/charmbracelet/huh) from Charm Bracelet.

Instead of typing `velero backup` commands over and over, you get a guided, interactive workflow: browse your backups, inspect their status, and delete obsolete ones — with a safety net at every step.

## ✨ Features

- 📋 **Overview** — list all Velero backups with status, errors, and warnings
- 🗑 **Multi-select deletion** — mark multiple backups with the spacebar for batch cleanup
- ⚠️ **Confirmation prompt** — explicit confirmation before any irreversible deletion
- 📅 **Sorted by creation date** — newest first, so current backups are always on top

<p align="center">
  <img src="assets/screenshot_01.png" alt="go-velero-manager screenshot" width="700"/>
</p>

## 📋 Requirements

- Go 1.19+
- [Velero](https://velero.io/) installed and configured
- Access to a [Kubernetes](https://kubernetes.io/) cluster

## 🚀 Installation

### Build from source

```bash
git clone https://github.com/mrtoadie/go-velero-manager.git
cd go-velero-manager
go build -o velero-manager
./velero-manager
```

### Run directly (without build)

```bash
go run main.go
```

## ✅ Tested on

Developed and tested on [Arch Linux](https://archlinux.org/) against a [k3s](https://k3s.io/) cluster.

| Environment | Status |
|---|---|
| [Arch Linux](https://archlinux.org/) | ✅ |
| [k3s](https://k3s.io/) cluster | ✅ |

## 📄 License

go-velero-manager is licensed under the [MIT License](LICENSE).
