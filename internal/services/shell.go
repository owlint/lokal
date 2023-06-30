package services

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/owlint/lokal/pkg/domain"
)

func StreamingCommand(command string) *exec.Cmd {
	fmt.Printf("> %s\n", command)
	cmd := exec.Command("bash", "-c", command)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd
}

func AddEnvs(command *exec.Cmd, envs []domain.EnvironmentVariable) {
	command.Env = []string{}

	for _, env := range envs {
		command.Env = append(command.Env, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}
	// append os environnement after to enable overwriting
	command.Env = append(command.Env, os.Environ()...)
}
