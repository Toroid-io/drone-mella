package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type (
	// Remote defines the server parameters
	Remote struct {
		Server string // URL of server. No /remote[...]
		Folder string // Folder where to store files
	}

	// Auth handles authentification
	Auth struct {
		User string // username
		Pass string // password
	}

	// Commit handles commit information
	Commit struct {
		Tag  string // tag if tag event
		Sha  string // commit sha
		Repo string // repo name
	}

	// Plugin defines the KiCad plugin parameters
	Plugin struct {
		Remote    Remote   // Remote configuration
		Files     []string // Local files
		TgzName   string   // Suffix
		Parentdir bool     // Include dir structure
		Auth      Auth     // Authentification
		Commit    Commit   // Commit information
		Verbose   bool     // Add -v to mella script
	}
)

func (p Plugin) Exec() error {

	var cmds []*exec.Cmd

	// Sanity checks
	if p.Auth.User == "" {
		return fmt.Errorf("No username provided!")
	}

	if p.Auth.Pass == "" {
		return fmt.Errorf("No password provided!")
	}

	if len(p.Files) == 0 {
		return fmt.Errorf("No local files provided!")
	}

	if p.Remote.Server == "" {
		return fmt.Errorf("No server provided!")
	}

	// Add webdav url
	u, err := url.Parse(p.Remote.Server)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, "remote.php/webdav")
	p.Remote.Server = u.String()

	var tarFile []string
	tarFile = append(tarFile, p.Commit.Repo, "_")
	if p.TgzName != "" {
		tarFile = append(tarFile, p.TgzName, "_")
	}
	if p.Commit.Tag != "" {
		tarFile = append(tarFile, p.Commit.Tag)
	} else {
		tarFile = append(tarFile, p.Commit.Sha[:7])
	}
	tarFile = append(tarFile, ".tar")

	var tgzFile = append(tarFile, ".gz")

	genConfig(p.Auth)
	for i, file := range p.Files {
		cmds = append(cmds, commandTAR(i, file, p.Parentdir, strings.Join(tarFile, "")))
	}
	cmds = append(cmds, commandGZIP(strings.Join(tarFile, "")))
	cmds = append(cmds, commandUPLOAD(p.Remote, strings.Join(tgzFile, ""), p.Verbose))

	// execute all commands in batch mode.
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func genConfig(a Auth) error {

	var buffer bytes.Buffer
	buffer.WriteString(a.User)
	buffer.WriteString(":")
	buffer.WriteString(a.Pass)
	buffer.WriteString("\n")

	return ioutil.WriteFile("auth.conf", buffer.Bytes(), 0777)
}

func commandTAR(index int, f string, parentdir bool, tarFile string) *exec.Cmd {

	var tarCmd []string
	var abs string

	if !parentdir {
		abs, _ = filepath.Abs(path.Dir(f))
		tarCmd = append(tarCmd, "cd", abs, "&&")
	}

	if index == 0 {
		tarCmd = append(tarCmd, "tar -cf")
	} else {
		tarCmd = append(tarCmd, "tar -uf")
	}
	abs, _ = filepath.Abs(tarFile)
	tarCmd = append(tarCmd, abs)

	if !parentdir {
		tarCmd = append(tarCmd, path.Base(f))
	} else {
		tarCmd = append(tarCmd, f)
	}

	// Calling bash allows wildcard expansion in files
	return exec.Command(
		"/bin/bash",
		"-c",
		strings.Join(tarCmd, " "),
	)
}

func commandGZIP(tarFile string) *exec.Cmd {

	return exec.Command(
		"gzip",
		tarFile,
	)
}

func commandUPLOAD(r Remote, tgzFile string, v bool) *exec.Cmd {

	u, _ := url.Parse(r.Server)
	u.Path = path.Join(u.Path, r.Folder)

	args := []string{"-c", "auth.conf"}
	if v {
		args = append(args, "-v")
	}
	args = append(args, tgzFile)
	args = append(args, u.String())

	return exec.Command(
		"mella",
		args...,
	)
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
