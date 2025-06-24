# Goputer

A lightweight system information Terminal User Interface (TUI) built with Go that provides real-time monitoring of your system's performance and resources.

## Features

- **CPU Monitoring**: View CPU usage, cores, threads, and detailed processor information
- **Memory Statistics**: Monitor RAM usage, cache, and memory distribution
- **Disk Information**: Display all mounted disks with available free space
- **Process Management**: View top processes and their current resource consumption
- **Real-time Updates**: Live monitoring with automatic refresh
- **Cross-platform**: Currently supports Windows (macOS and Linux support under testing)

## Installation

### Prerequisites
- Go 1.24.4 or higher

### From Source
```bash
git clone https://github.com/mkael1/goputer.git
cd goputer
go build -o goputer
```

### Using Go Install
```bash
go install github.com/mkael1/goputer@latest
```

## Usage

Simply run the executable to start the TUI:

```bash
./goputer
```

### Navigation
- Use arrow keys to navigate between different sections
- Press `q` or `Ctrl+C` to quit

## Screenshots

```
TBD
```

## Platform Support

| Platform | Status |
|----------|---------|
| Windows  | ✅ Fully Supported |
| macOS    | 🧪 Testing Required |
| Linux    | 🧪 Testing Required |

## Core Dependencies
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Modern TUI framework for building terminal applications
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - Common TUI components (spinners, text inputs, etc.)
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Styling and layout for terminal interfaces
- **[gopsutil](https://github.com/shirou/gopsutil)** - Cross-platform system and process utilities for Go
- **golang.org/x/term** - Terminal control and utilities

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. Areas where help is especially needed:

- Testing and validation on macOS and Linux platforms
- Performance optimizations
- Additional system metrics
- UI/UX improvements

### Development Setup

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] Complete macOS and Linux testing
- [ ] Add details to each main panel (CPU,MEM,Disk)
- [ ] Add a network panel
- [ ] Add some tests! :D

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/mkael1/goputer/issues) on GitHub.

---

Made with ❤️ in Go
