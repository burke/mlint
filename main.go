package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(cmd, path string) error {
	command := exec.Command("/bin/sh", "-c", cmd+" "+path)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func banner(color, emoji, cmd, path string) {
	cols := termWidth()
	pw := cols - len(cmd) - len(path) - 7
	padding := ""
	for i := 0; i < pw; i++ {
		padding += " "
	}
	fmt.Printf("\x1b[4%sm %s  \x1b[0;3%sm %s \x1b[0;4%s;39m %s%s\x1b[0m\n",
		color, emoji, color, cmd, color, path, padding)
}

func termWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	parts := strings.Split(strings.TrimSpace(string(out)), " ")
	if len(parts) < 2 {
		return 0
	}

	i, _ := strconv.Atoi(parts[1])
	return i
}

func check(cmd, path string) (ok bool) {
	banner("4", "🔍", cmd+" running", path)
	if err := run(cmd, path); err != nil {
		banner("1", "❌", cmd+" failed", "")
		return false
	} else {
		banner("2", "✅", cmd+" passed", "")
		return true
	}
}

func setup() {
	os.Mkdir("/tmp/mlint", 0755)
	gp := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", gp)
	os.Setenv("GOPATH", "/tmp/mlint")
	os.Setenv("PATH", "/tmp/mlint/bin:"+os.Getenv("PATH"))
	if _, err := os.Stat("/tmp/mlint/bin/vet"); err != nil {
		fmt.Println(exec.Command("go", "get", "golang.org/x/tools/cmd/vet").Output())
	}
	if _, err := os.Stat("/tmp/mlint/bin/errcheck"); err != nil {
		fmt.Println(exec.Command("go", "get", "github.com/burke/errcheck").Output())
	}
	if _, err := os.Stat("/tmp/mlint/bin/golint"); err != nil {
		fmt.Println(exec.Command("go", "get", "github.com/golang/lint/golint").Output())
	}
}

func main() {
	setup()

	ok := true
	checks := []string{"go fmt", "go vet", "errcheck", "golint", "go test"}
	for _, chk := range checks {
		if !check(chk, "./...") {
			ok = false
		}
	}

	if ok {
		os.Exit(0)
	}
	os.Exit(1)
}
