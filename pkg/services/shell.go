package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/owlint/lokal/pkg/domain"
)

type ShellResult struct {
	Error error
}

func StartStreamingCommand(name string, args ...string) *exec.Cmd {
	cmd := StreamingCommand(name, args...)
	cmd.Start()
	return cmd
}

func StreamingCommand(name string, args ...string) *exec.Cmd {
	fmt.Printf("> %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd
}

func WaitChan(command *exec.Cmd) <-chan ShellResult {
	done := make(chan ShellResult)
	go func() {
		err := command.Wait()
		done <- ShellResult{
			Error: err,
		}
		close(done)
	}()

	return done
}

func KillAll(commands []*exec.Cmd) {
	for _, cmd := range commands {
		cmd.Process.Kill()
	}
}

func AddEnvs(command *exec.Cmd, envs []domain.EnvironmentVariable) {
	command.Env = []string{}

	for _, env := range envs {
		command.Env = append(command.Env, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}
	// append os environnement after to enable overwriting
	command.Env = append(command.Env, os.Environ()...)
}
