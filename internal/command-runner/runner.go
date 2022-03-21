package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

type CommandRunner struct {
	Command       *Command
	CommandOutPut *CommandOutPut
	err           error
}
type d struct{}
type Command struct {
	Name  string
	Flags []string
	StdIn string
}

type CommandOutPut struct {
	StdOut string
	StdErr string
}

func New() *CommandRunner {
	return &CommandRunner{}
}

func (c *CommandRunner) Run() {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, "cmd")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Print("gjhgg")
		log.Fatalln(err)
	}

	go func() {
		defer stdin.Close()
		_, err := io.WriteString(stdin, "ls")
		if err != nil {
			log.Fatalln(err)
		}
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(err)
	}

	fmt.Print(bytes.NewBuffer(out).String())
}
