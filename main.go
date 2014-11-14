package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(cmd, path string) error {
	return exec.Command("/bin/sh", "-c", cmd+" "+path).Run()
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

func main() {
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
