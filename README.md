# Development Setup

This project uses Go and can be developed using any IDE or editor of your choice. Below are some recommended setups, but feel free to use whatever tools you prefer.

## Option 1: Using Nix (Recommended)

This project uses Nix to manage the Go development environment. This ensures consistent development environments across all developers.

> **Note:** Nix is only supported on Linux and macOS environments. Windows users should use Option 2 below.

### Prerequisites
- Install [Nix](https://nixos.org/download.html) (Linux/macOS only)
- Install your preferred IDE/editor

### Launching with Nix

If you're using VSCode or Cursor, you can launch them directly from nix-shell:

#### For VSCode:
```bash
nix-shell --command "code ."
```

#### For Cursor:
```bash
nix-shell --command "cursor ."
```

This will launch your IDE with the correct Go environment and all necessary dependencies.

## Option 2: Using Global Go Installation

If you prefer not to use Nix, or if you're on Windows, you can install Go globally on your system:

### Prerequisites
- Install Go from [golang.org](https://golang.org/dl/)
- Install your preferred IDE/editor

### Setup
1. Install Go following the official installation guide for your operating system
2. Set up your GOPATH and add it to your PATH
3. Install any Go-related extensions for your IDE (if needed)
4. Open the project directly in your IDE

## IDE Extensions (Optional)

### VSCode
- Install the official Go extension: `golang.go`

### Cursor
- Go support is built-in, no additional extensions required

## Verifying Setup

To verify your Go installation, run:
```bash
go version
```

You should see the Go version information displayed.

## Note
These instructions are meant to be helpful starting points. Feel free to use any development setup that works best for you - the only requirement is having Go installed and properly configured on your system. 