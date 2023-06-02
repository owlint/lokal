package commands

import (
	"log"
	"os"
	"path/filepath"

	"github.com/owlint/lokal/pkg/services"
	"github.com/urfave/cli/v2"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var RunCommand = &cli.Command{
	Name:  "run",
	Usage: "Run command alongside services",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "namespace",
			Usage:    "Namespace",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "pod",
			Usage:    "Pod",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "container",
			Usage:    "Container",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "command",
			Usage:    "Command to run",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		// TODO: get from env
		kubeconfig := filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatal(err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatal(err)
		}

		describer := services.NewPodDescriber(clientset)

		// TODO: merge with config file
		envs, err := describer.ReadEnvs(c.Context, c.String("namespace"), c.String("pod"), c.String("container"))
		if err != nil {
			return err
		}

		process := services.StreamingCommand("bash", "-c", c.String("command"))
		services.AddEnvs(process, envs)

		return process.Run()
	},
}
