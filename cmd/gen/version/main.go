package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	version, _ := os.LookupEnv("BRING_VERSION")
	if !strings.HasPrefix(version, "v") {
		version = "v0.0.0-edge"
	}

	build_time := time.Now().Format(time.RFC3339)

	git_hash, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching git hash: %v\n", err)
		os.Exit(1)
		return
	}

	git_status := exec.Command("git", "status", "--porcelain")
	var git_status_output bytes.Buffer
	git_status.Stdout = &git_status_output
	if err := git_status.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error checking git status: %v\n", err)
		os.Exit(1)
		return
	}

	git_dirty := strings.TrimSpace(git_status_output.String()) != ""

	// Write the generated code to a file
	content := fmt.Sprintf(`// Code generated by go generate; DO NOT EDIT.
package cmd

func init(){
	b := &_buildInfo
	b.Version   = %q
	b.TimeBuild = %q
	b.GitRev    = %q
	b.GitDirty  = %v 
}
`,
		version,
		build_time,
		strings.TrimSpace(string(git_hash)),
		git_dirty,
	)

	err = os.WriteFile("version.g.go", []byte(content), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing build_info.go: %v\n", err)
		os.Exit(1)
	}
}
