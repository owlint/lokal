package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/owlint/lokal/pkg/domain"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodDescriber struct {
	kubernetes.Interface
}

func (d *PodDescriber) ReadEnvs(ctx context.Context, namespace, pod, container string) ([]domain.EnvironmentVariable, error) {
	// TODO: retrieve from deployement
	pods, err := d.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	podDescriptor := d.latestPod(*pods, pod)
	if podDescriptor == nil {
		return nil, fmt.Errorf("invalid pod %s", pod)
	}

	for _, containerDescriptor := range podDescriptor.Spec.Containers {
		if containerDescriptor.Name == container {

			return d.extractContainerEnv(ctx, namespace, containerDescriptor)
		}
	}

	return nil, errors.New("container not found")
}

func (d *PodDescriber) latestPod(pods v1.PodList, prefix string) *v1.Pod {
	var latest *v1.Pod
	for _, pod := range pods.Items {
		pod := pod
		if !strings.HasPrefix(pod.Name, prefix) {
			continue
		}

		if latest == nil || pod.CreationTimestamp.Time.After(latest.CreationTimestamp.Time) {
			latest = &pod
		}
	}

	return latest
}

func (d *PodDescriber) extractContainerEnv(ctx context.Context, namespace string, container v1.Container) ([]domain.EnvironmentVariable, error) {
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

func (d *PodDescriber) extractContainerSecrets(ctx context.Context, namespace string, container v1.Container) ([]domain.EnvironmentVariable, error) {
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

func (d *PodDescriber) extractLocalSecrets(ctx context.Context, namespace, name string) ([]domain.EnvironmentVariable, error) {
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
