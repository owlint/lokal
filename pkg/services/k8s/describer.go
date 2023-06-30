package k8s

import (
	"context"
	"fmt"
	"net/url"

	"github.com/owlint/lokal/pkg/domain"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentDescriber struct {
	kubernetes.Interface

	overrideNamespace bool
}

func NewDeploymentDescriber(clientset *kubernetes.Clientset, overrideNamespace bool) *DeploymentDescriber {
	return &DeploymentDescriber{
		Interface:         clientset,
		overrideNamespace: overrideNamespace,
	}
}

func (d *DeploymentDescriber) ReadEnvs(ctx context.Context, namespace, deployment, container string) ([]domain.EnvironmentVariable, error) {
	deploymentDescriptor, err := d.AppsV1().Deployments(namespace).Get(ctx, deployment, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	for _, containerDescriptor := range deploymentDescriptor.Spec.Template.Spec.Containers {
		if containerDescriptor.Name == container {
			envs, err := d.extractContainerEnv(ctx, namespace, containerDescriptor)
			if err != nil {
				return nil, err
			}

			if d.overrideNamespace {
				envs, err = d.OverrideNamespace(ctx, namespace, envs)
				if err != nil {
					return nil, err
				}
			}

			return envs, nil
		}
	}

	return nil, fmt.Errorf("container %s not found", container)
}

func (d *DeploymentDescriber) extractContainerEnv(ctx context.Context, namespace string, container v1.Container) ([]domain.EnvironmentVariable, error) {
	envs := []domain.EnvironmentVariable{}
	for _, env := range container.Env {
		envs = append(envs, domain.EnvironmentVariable{
			Name:  env.Name,
			Value: env.Value,
		})
	}

	secrets, err := d.extractContainerSecrets(ctx, namespace, container)
	if err != nil {
		return nil, err
	}

	envs = append(envs, secrets...)

	return envs, nil
}

func (d *DeploymentDescriber) extractContainerSecrets(ctx context.Context, namespace string, container v1.Container) ([]domain.EnvironmentVariable, error) {
	envs := []domain.EnvironmentVariable{}
	for _, env := range container.EnvFrom {
		if env.SecretRef == nil {
			continue
		}
		secrets, err := d.extractLocalSecrets(ctx, namespace, env.SecretRef.LocalObjectReference.Name)
		if err != nil {
			return nil, err
		}

		envs = append(envs, secrets...)
	}

	return envs, nil
}

func (d *DeploymentDescriber) extractLocalSecrets(ctx context.Context, namespace, name string) ([]domain.EnvironmentVariable, error) {
	secrets, err := d.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	envs := []domain.EnvironmentVariable{}
	for name, secret := range secrets.Data {
		envs = append(envs, domain.EnvironmentVariable{
			Name:  name,
			Value: string(secret),
		})
	}

	return envs, nil
}

func (d *DeploymentDescriber) OverrideNamespace(ctx context.Context, namespace string, envs []domain.EnvironmentVariable) ([]domain.EnvironmentVariable, error) {
	deployments, err := d.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	deploymentNames := make(map[string]struct{}, len(deployments.Items))
	for _, deployment := range deployments.Items {
		deploymentNames[deployment.GetName()] = struct{}{}
	}

	overridedEnvs := make([]domain.EnvironmentVariable, 0, len(envs))
	for _, env := range envs {
		env := env
		if u, err := url.Parse(env.Value); err == nil {
			if _, exists := deploymentNames[u.Hostname()]; exists {
				port := u.Port()
				if len(port) > 0 {
					port = fmt.Sprintf(":%s", port)
				}
				u.Host = fmt.Sprintf("%s.%s%s", u.Hostname(), namespace, port)

				env.Value = u.String()
			}

			overridedEnvs = append(overridedEnvs, env)
		}
	}

	return overridedEnvs, nil
}
