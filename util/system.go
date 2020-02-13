package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// Directory is the directory address of the cli and all darknodes data.
var Directory = filepath.Join(os.Getenv("HOME"), ".darknode")

// NodePath return the absolute directory of the node with given name.
func NodePath(name string) string {
	return filepath.Join(Directory, "darknodes", name)
}

// BackUpConfig copies the config file of the node to the backup folder under
// .darknode directory in case something unexpected happens.
func BackUpConfig(name string) error{
	path := NodePath(name)
	backupFolder := filepath.Join(Directory, "backup", name)
	if err := Run("mkdir", "-p", backupFolder); err != nil {
		return err
	}
	backup := fmt.Sprintf("cp %v %v", filepath.Join(path, "config.json"), backupFolder)
	return Run("bash", "-c", backup)
}

// run the command and pipe the output to the stdout
func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SilentRun runs the commands with no output
func SilentRun(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	return cmd.Run()
}

// CommandOutput runs a series of commands with bash
func CommandOutput(commands string) (string, error) {
	cmd := exec.Command("bash", "-c", commands)
	output, err := cmd.Output()
	return string(output), err
}

// RemoteRun runs the script on the instance which host the darknode of given name.
func RemoteRun(name, script string) error {
	return RemoteRunWithUser(name, script, "darknode")
}

// RemoteRun runs the script on the instance as specific system user.
func RemoteRunWithUser(name, script, user string) error {
	key, err := ParseSshPrivateKey(name)
	if err != nil {
		return err
	}
	config := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the instance using ssh
	ip, err := IP(name)
	if err != nil {
		return err
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:22", ip), &config)
	if err != nil {
		return err
	}

	// Create a new session to run the script
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Redirect the remote stdin, stdout and stderr to local.
	sessStdIn, err := session.StdinPipe()
	if err != nil {
		return err
	}
	go io.Copy(sessStdIn, os.Stdin)
	sessStdOut, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStdErr, err := session.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stderr, sessStdErr)

	return session.Run(script)
}

// OpenInBrowser tries to open the url with system default browser. It ignores the error if failing.
func OpenInBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		SilentRun("open", url)
	case "linux":
		if CheckWSL() {
			SilentRun("cmd.exe", "/C", "start", url)
		} else {
			SilentRun("xdg-open", url)
		}
	}
	return nil
}

// CheckWSL if the linux system is a Subsystem of window.
func CheckWSL() bool {
	file, err := ioutil.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(string(file), "Microsoft")
}
