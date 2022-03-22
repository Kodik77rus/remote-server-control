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

const _timeLimit = 5 * time.Second // max time for executing command

var (
	errUnsupportedOsType  = errors.New("unsupported os type")
	errExecutionTimeLimit = errors.New("interrupted: execution time limit exceeded")
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
	StdOut *string
	StdErr *string
}

type FailedExecutionCommand struct {
	err error
}

func New(c *Command, os string) *CommandRunner {
	return &CommandRunner{
		Command: c,
		OsType:  os,
	}
}

func (c *CommandRunner) Run() *CommandRunner {

	switch c.OsType {
	case "windows":
		return c.runWindowsCommand(newWindowsCommand(c.Command))
	case "linux":
		return c.runLinuxCommand(newLinuxCommand(c.Command))
	default:
		return &CommandRunner{RunnerError: errUnsupportedOsType}
	}
}

func (c *CommandRunner) runWindowsCommand(winCmd *windowsCommnd) *CommandRunner {
	ctx, cancel := context.WithTimeout(context.Background(), _timeLimit)
	defer cancel()

	cmd := exec.CommandContext(ctx, "powershell")

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

	log.Printf("Start execute command: %s", winCmd.stdIn)

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Print(errExecutionTimeLimit.Error())
			return &CommandRunner{RunnerError: errExecutionTimeLimit}
		}

		log.Print(FailedExecutionCommand{err: err}.Error())
		return &CommandRunner{RunnerError: FailedExecutionCommand{err: err}}
	}

	log.Printf("Done output:\n %s", stdOut.String())

	cmdOutPut := NewCommandOutPut(&stdOut, &stdErr)

	c.setOutput(cmdOutPut)

	return c
}

func (c *CommandRunner) runLinuxCommand(linCmd *linuxCommnd) *CommandRunner {
	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// defer cancel()
	return nil
}

func (c *CommandRunner) setOutput(out *CommandOutPut) {
	c.CommandOutPut = out
}

func (f FailedExecutionCommand) Error() string {
	return fmt.Sprintf("failed execute command: %v", f.err)
}

func newWindowsCommand(c *Command) *windowsCommnd {
	return &windowsCommnd{
		stdIn: fmt.Sprintf("%v %v %v", c.Name, strings.Join(c.Args, ""), c.StdIn),
	}
}

func newLinuxCommand(c *Command) *linuxCommnd {
	return nil
}

func NewCommandOutPut(stdOut, stderr *bytes.Buffer) *CommandOutPut {
	stdt := new(string)
	stdrr := new(string)

	*stdt = stdOut.String()
	*stdrr = stderr.String()

	return &CommandOutPut{
		StdOut: stdt,
		StdErr: stdrr,
	}
}

func fatalError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
