package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

type CommandRunner struct {
	Command       *Command
	CommandOutPut *CommandOutPut

	OsType      string
	RunnerError error
}

type Command struct {
	Name  string
	Args  []string
	StdIn string
}

type linuxCommnd struct {
	name string
	Args []string
}

type windowsCommnd struct {
	stdIn string
}

type CommandOutPut struct {
	StdOut *bytes.Buffer
	StdErr *bytes.Buffer
}

var errUnsupportedOsType = errors.New("unsupported os type")

func New(c *Command, os string) *CommandRunner {
	return &CommandRunner{
		Command: c,
		OsType:  os,
	}
}

func (c *CommandRunner) Run() *CommandRunner {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch c.OsType {
	case "windows":
		return c.runWindowsCommand(&ctx, newWindowsCommand(c.Command))
	case "linux":
		return c.runLinuxCommand(&ctx, newLinuxCommand(c.Command))
	default:
		return &CommandRunner{RunnerError: errUnsupportedOsType}
	}
}

func (c *CommandRunner) runWindowsCommand(ctx *context.Context, winCmd *windowsCommnd) *CommandRunner {
	cmd := exec.CommandContext(*ctx, "powershell")

	stdin, err := cmd.StdinPipe()
	fatalError(err)

	go func() {
		defer stdin.Close()
		_, err := io.WriteString(stdin, winCmd.stdIn)
		fatalError(err)
	}()

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return &CommandRunner{RunnerError: err}
	}

	c.setOutput(&stdOut, &stdErr)

	return c
}

func (c *CommandRunner) runLinuxCommand(ctx *context.Context, linCmd *linuxCommnd) *CommandRunner {
	return nil
}

func newWindowsCommand(c *Command) *windowsCommnd {
	return &windowsCommnd{
		stdIn: fmt.Sprintf("%v %v %v", c.Name, strings.Join(c.Args, ""), c.StdIn),
	}
}

func newLinuxCommand(c *Command) *linuxCommnd {
	return nil
}

func (c *CommandRunner) setOutput(stdOut, Stderr *bytes.Buffer) {
	out := &CommandOutPut{
		StdOut: stdOut,
		StdErr: Stderr,
	}

	c.CommandOutPut = out
}

func fatalError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
