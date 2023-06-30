package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/owlint/lokal/internal/services"
	localConfig "github.com/owlint/lokal/pkg/config"
	"github.com/owlint/lokal/pkg/services/k8s"
	"github.com/urfave/cli/v2"
)

var RunCommand = &cli.Command{
	Name:  "run",
	Usage: "Run a local application",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "namespace",
			Usage:    "Namespace within cluster. Overrides config file if provided.",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "deployment",
			Usage:    "Deployment within namespace. Overrides config file if provided.",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "container",
			Usage:    "Container within deployment. Will use pod name if ommited.",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "command",
			Usage:    "Command to run.",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "force-namespace",
			Usage:    "Append namespace to url environnement variables referencing local deployment. Ex: http://myapp/ will be converted to http://myapp.mynamespace/. This feature is useful when running lokal alongside Telepresence.",
			Value:    true,
			Required: false,
		},
		&cli.StringFlag{
			Name:     "config",
			Usage:    "Path to the local config file.",
			Required: false,
			Value:    "./lokal.yaml",
		},
		&cli.StringFlag{
			Name:     "kube-config",
			Usage:    "Path to the kube config file.",
			Required: false,
			Value: filepath.Join(
				os.Getenv("HOME"), ".kube", "config",
			),
		},
	},
	Action: func(c *cli.Context) error {
		configPath := c.String("config")
		config, err := localConfig.ReadLocalConfig(configPath)
		if err != nil {
			fmt.Printf("Couldn't read lokal config %s, ignoring...\n", configPath)
			config = &localConfig.LocalConfig{}
		}

		forceNamespace := c.Bool("force-namespace")
		if forceNamespace {
			config.ForceNamespace = forceNamespace
		}
		namespace := c.String("namespace")
		if namespace != "" {
			config.Namespace = namespace
		}
		deployment := c.String("deployment")
		if deployment != "" {
			config.Deployment = deployment
		}
		container := c.String("container")
		if container != "" {
			config.Container = container
		}
		if config.Container == "" {
			config.Container = config.Deployment
		}

		command := c.String("command")
		if command != "" {
			config.Command = command
		}

		err = config.EnsureValid()
		if err != nil {
			return err
		}

		clientset, err := k8s.NewClientSet(c.String("kube-config"))
		if err != nil {
			return err
		}

		describer := k8s.NewDeploymentDescriber(clientset, config.ForceNamespace)

		ctx, cancel := context.WithTimeout(c.Context, 30*time.Second)
		defer cancel()
		envs, err := describer.ReadEnvs(ctx, config.Namespace, config.Deployment, config.Container)
		if err != nil {
			return err
		}

		process := services.StreamingCommand(config.Command)
		services.AddEnvs(process, append(envs, config.Env...))

		return process.Run()
	},
}
