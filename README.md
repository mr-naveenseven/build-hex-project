# Build Hex Project

A simple command-line tool to quickly scaffold a **Hexagonal Architecture (Ports & Adapters)** project in Go.

## Features

* Create a new Go project in seconds
* Generate a clean Hexagonal Architecture folder structure
* Create common project directories (`cmd`, `internal`, `configs`, `docs`, `scripts`, etc.)
* Automatically generate module folders (e.g. `users`, `orders`, `products`)
* Generate starter files such as:

  * `README.md`
  * `.env`
  * `.gitignore`
  * `Makefile`
  * `cmd/api/main.go`
* Initialize a Go module
* Optionally initialize a Git repository

## Usage

```bash
go run main.go
```

Follow the interactive prompts to provide:

* Project name
* Go module name
* Feature modules (comma-separated)
* Git initialization option

## Example

```text
Project Name: ecommerce

Go Module: github.com/username/ecommerce

Modules:
users,products,orders

Initialize Git? (y/n): y
```

The tool will generate a production-friendly Go project structure based on Hexagonal Architecture.

## License

MIT
