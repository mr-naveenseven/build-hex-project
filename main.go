package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var platformDirs = []string{
	"config",
	"database",
	"auth",
	"middleware",
	"router",
	"logger",
}

var commonDirs = []string{
	"cmd/api",
	"configs",
	"docs",
	"pkg",
	"scripts",
	"migrations",
	"internal/platform",
	"internal/shared/constants",
	"internal/shared/errors",
	"internal/shared/utils",
	"internal/shared/pagination",
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Go Hexagonal Project Generator ===")

	project := prompt(reader, "Project Name")
	if project == "" {
		fmt.Println("project name required")
		return
	}

	module := prompt(reader, "Go Module (leave empty to use project name)")
	if module == "" {
		module = project
	}

	modulesInput := prompt(reader, "Modules (comma separated, leave empty for template)")
	createGit := strings.EqualFold(prompt(reader, "Initialize Git? (y/n)"), "y")

	var modules []string
	if strings.TrimSpace(modulesInput) == "" {
		modules = []string{"module"}
	} else {
		for _, m := range strings.Split(modulesInput, ",") {
			m = strings.TrimSpace(m)
			if m != "" {
				modules = append(modules, m)
			}
		}
	}

	if err := createProject(project, module, modules, createGit); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("\nProject created successfully!")
}

func prompt(r *bufio.Reader, text string) string {
	fmt.Printf("%s: ", text)
	s, _ := r.ReadString('\n')
	return strings.TrimSpace(s)
}

func createProject(root, module string, modules []string, git bool) error {
	mkdir(root)

	for _, d := range commonDirs {
		mkdir(filepath.Join(root, d))
	}

	for _, d := range platformDirs {
		mkdir(filepath.Join(root, "internal", "platform", d))
	}

	for _, m := range modules {
		createModule(root, m)
	}

	write(filepath.Join(root, ".env"), "")
	write(filepath.Join(root, "README.md"), "# "+root+"\n")
	write(filepath.Join(root, ".gitignore"), gitignore())
	write(filepath.Join(root, "Makefile"), makefile())

	initCmdMain(root)
	initGoMod(root, module)

	if git {
		run(root, "git", "init")
	}

	return nil
}

func createModule(root, name string) {
	base := filepath.Join(root, "internal", name)

	dirs := []string{
		"domain",
		"service",
		"repository",
		"http/request",
		"http/response",
		"ports",
	}

	for _, d := range dirs {
		mkdir(filepath.Join(base, d))
	}

	write(filepath.Join(base, "domain", "entity.go"),
		"package domain\n\n")
	write(filepath.Join(base, "domain", "repository.go"),
		"package domain\n\n")
	write(filepath.Join(base, "service", "service.go"),
		"package service\n\n")
	write(filepath.Join(base, "repository", "repository.go"),
		"package repository\n\n")
	write(filepath.Join(base, "http", "handler.go"),
		"package http\n\n")
	write(filepath.Join(base, "http", "routes.go"),
		"package http\n\n")
	write(filepath.Join(base, "ports", "inbound.go"),
		"package ports\n\n")
	write(filepath.Join(base, "ports", "outbound.go"),
		"package ports\n\n")
}

func initCmdMain(root string) {
	content := `package main

import "fmt"

func main() {
	fmt.Println("Server Started")
}
`
	write(filepath.Join(root, "cmd", "api", "main.go"), content)
}

func initGoMod(root, module string) {
	run(root, "go", "mod", "init", module)
}

func mkdir(path string) {
	_ = os.MkdirAll(path, 0755)
}

func write(path, content string) {
	_ = os.WriteFile(path, []byte(content), 0644)
}

func run(dir string, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = nil
	}
	_ = cmd.Run()
}

func gitignore() string {
	return `.env
vendor/
bin/
dist/
*.exe
*.out
`
}

func makefile() string {
	return `run:
	go run ./cmd/api

build:
	go build -o bin/app ./cmd/api

tidy:
	go mod tidy
`
}
