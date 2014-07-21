package gordon

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	VerboseOutput = false
)

type remote struct {
	Name string
	Url  string
}

func PrintCommand(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stderr, "executing %q ...\n", strings.Join(cmd.Args, " "))
}

func Fatalf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func GetOriginUrl() (string, string, error) {
	return GetRemoteUrl("origin")
}

func GetRemoteUrl(remote string) (string, string, error) {
	remotes, err := getRemotes()
	if err != nil {
		return "", "", err
	}
	for _, r := range remotes {
		if r.Name == remote {
			parts := strings.Split(r.Url, "/")

			org := parts[len(parts)-2]
			if i := strings.LastIndex(org, ":"); i > 0 {
				org = org[i+1:]
			}

			name := parts[len(parts)-1]
			name = strings.TrimRight(name, ".git")

			return org, name, nil
		}
	}
	return "", "", nil
}

func GetMaintainerManagerEmail() (string, error) {
	cmd := exec.Command("git", "config", "user.email")
	if VerboseOutput {
		PrintCommand(cmd)
	}
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git config user.email: %v", err)
	}
	return string(bytes.Split(output, []byte("\n"))[0]), err
}

// Return the remotes for the current dir
func getRemotes() ([]remote, error) {
	cmd := exec.Command("git", "remote", "-v")
	if VerboseOutput {
		PrintCommand(cmd)
	}
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	out := []remote{}
	s := bufio.NewScanner(bytes.NewBuffer(output))
	for s.Scan() {
		o := remote{}
		if _, err := fmt.Sscan(s.Text(), &o.Name, &o.Url); err != nil {
			return nil, err
		}
		out = append(out, o)
	}

	return out, nil
}

// Execute git commands and output to
// Stdout and Stderr
func Git(args ...string) error {
	cmd := exec.Command("git", args...)
	if VerboseOutput {
		PrintCommand(cmd)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func GetTopLevelGitRepo() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	if VerboseOutput {
		PrintCommand(cmd)
	}
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(output), "\n"), nil
}
