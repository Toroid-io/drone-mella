package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
)

type (
	// Remote defines the server parameters
	Remote struct {
		Server string // URL of ownCloud server. No /remote[...]
		Folder string // Folder where to store files
	}

	// Local defines the local files parameters
	Local struct {
		Folder string // Local folder to upload
		Files  string // Local files to upload
	}

	// Auth handles authentification
	Auth struct {
		User string // ownCloud username
		Pass string // ownCloud password
	}

	// Plugin defines the KiCad plugin parameters
	Plugin struct {
		Remote  Remote // Remote configuration
		Local   Local  // Local configuration
		Auth    Auth   // Authentification
		Verbose bool   // Add -v to mella script
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

	if p.Local.Folder == "" && p.Local.Files == "" {
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

	genConfig(p.Auth)
	cmds = append(cmds, commandTAR(p.Local))
	cmds = append(cmds, commandUPLOAD(p.Remote, p.Local, p.Verbose))

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

func commandTAR(l Local) *exec.Cmd {

	var buffer bytes.Buffer
	buffer.WriteString("tar -czf ")
	buffer.WriteString(l.Folder)
	buffer.WriteString(".tgz ")
	buffer.WriteString(path.Join(l.Folder, l.Files))

	// Calling bash allows wildcard expansion in files
	return exec.Command(
		"/bin/bash",
		"-c",
		buffer.String(),
	)
}

func commandUPLOAD(r Remote, l Local, v bool) *exec.Cmd {

	u, _ := url.Parse(r.Server)
	u.Path = path.Join(u.Path, r.Folder)

	var buffer bytes.Buffer
	buffer.WriteString(l.Folder)
	buffer.WriteString(".tgz")

	verbose := ""
	if v {
		verbose = "-v"
	}

	return exec.Command(
		"mella",
		"-c",
		"auth.conf",
		verbose,
		buffer.String(),
		u.String(),
	)
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
