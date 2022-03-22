package runner

import (
	"bytes"
	"context"
	"io"
	"log"
	"os/exec"
	"time"
)

type CommandRunner struct {
	Command       *Command
	CommandOutPut *CommandOutPut

	osType string
	// err    error
}

type Command struct {
	Name  string
	Flags []string
	StdIn string
}

type CommandOutPut struct {
	StdOut *bytes.Buffer
	StdErr *bytes.Buffer
}

func New(c *Command, os string) *CommandRunner {

	return &CommandRunner{
		Command: c,
		osType:  os,
	}
}

func (c *CommandRunner) Run() *CommandRunner {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch c.osType {
	case "windows":
		return c.runWindowsCommand(&ctx)
	case "linux":
		return c.runLinuxCommand(&ctx)
	default:
		return nil
	}
}

func (c *CommandRunner) runWindowsCommand(ctx *context.Context) *CommandRunner {
	cmd := exec.CommandContext(*ctx, "powershell")

	stdin, err := cmd.StdinPipe()
	fatalError(err)

	go func() {
		defer stdin.Close()
		_, err := io.WriteString(stdin, c.Command.Name)
		fatalError(err)
	}()

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		log.Print(err)
	}

	return c.setOutput(&stdOut, &stdErr)
}

func (c *CommandRunner) runLinuxCommand(ctx *context.Context) *CommandRunner {

	return nil
}

func (c *CommandRunner) setOutput(stdOut, Stderr *bytes.Buffer) *CommandRunner {
	out := &CommandOutPut{
		StdOut: stdOut,
		StdErr: Stderr,
	}

	c.CommandOutPut = out
	return c
}

func fatalError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
