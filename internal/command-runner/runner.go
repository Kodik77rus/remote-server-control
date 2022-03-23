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

	runner      *exec.Cmd
	OsType      string
	RunnerError error
}

type Command struct {
	Name  string
	Args  []string
	StdIn string
}

type linuxCommnd struct {
	name  string
	args  []string
	stdIn string
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
		runner:  new(exec.Cmd),
	}
}

func (c *CommandRunner) Run() *CommandRunner {
	switch c.OsType {
	case "windows":
		return c.run(newWindowsCommand(c.Command))
	case "linux":
		return c.run(newLinuxCommand(c.Command))
	default:
		return &CommandRunner{RunnerError: errUnsupportedOsType}
	}
}

func (c *CommandRunner) setRunner(runner *exec.Cmd) {
	c.runner = runner
}

func (c *CommandRunner) run(i interface{}) *CommandRunner {
	ctx, cancel := context.WithTimeout(context.Background(), _timeLimit)
	defer cancel()

	switch i := i.(type) {
	case *windowsCommnd:
		cmd := exec.CommandContext(ctx, "powershell")
		c.setRunner(cmd)
	case *linuxCommnd:
		cmd := exec.CommandContext(ctx, i.name, i.args...)
		c.setRunner(cmd)
	}

	stdin, err := c.runner.StdinPipe()
	fatalError(err)

	go func() {
		defer stdin.Close()
		_, err := io.WriteString(stdin, c.Command.StdIn)
		fatalError(err)
	}()

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	c.runner.Stdout = &stdOut
	c.runner.Stderr = &stdErr

	log.Printf("Start execute command: %+v", c.Command)

	if err := c.runner.Run(); err != nil {
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
	return &linuxCommnd{
		name:  c.Name,
		args:  c.Args,
		stdIn: c.StdIn,
	}
}

func NewCommandOutPut(stdOut, stderr *bytes.Buffer) *CommandOutPut {
	var stdt string
	var stdrr string

	stdt = stdOut.String()
	stdrr = stderr.String()

	return &CommandOutPut{
		StdOut: &stdt,
		StdErr: &stdrr,
	}
}

func fatalError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
