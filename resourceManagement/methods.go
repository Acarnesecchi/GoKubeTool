package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobConfig struct {
	JobName        string `yaml:"jobname"`
	Namespace      string `yaml:"namespace"`
	ContainerImage string `yaml:"containerImage"`
	Command        string `yaml:"command"`
	RestartPolicy  string `yaml:"restartPolicy"`
	BackoffLimit   int32  `yaml:"backoffLimit"`
}

func fetch(k *KubernetesClient) {
	// Fetches the list of namespaces.
	namespaces, err := k.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	check(err)

	fmt.Println("Namespaces in the cluster:")
	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
	}
}

func resetDB(k *KubernetesClient, f string) {
	config, err := parseConfig(f)
	check(err)
	var cmd []string
	var restartPolicy v1.RestartPolicy
	namespace := "default"
	if config.Namespace != "" {
		namespace = config.Namespace
	}
	if config.Command != "" {
		cmd = strings.Fields(config.Command) //entrypoint formatted ["", ""]
	}
	if config.RestartPolicy != "" {
		restartPolicy = v1.RestartPolicy(config.RestartPolicy)
	} else {
		restartPolicy = v1.RestartPolicyNever
	}
	fmt.Printf("container image: %s", config.ContainerImage)
	jobs := k.client.BatchV1().Jobs(namespace)
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.JobName,
			Namespace: config.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    config.JobName,
							Image:   config.ContainerImage,
							Command: cmd,
						},
					},
					RestartPolicy: restartPolicy,
				},
			},
			BackoffLimit: &config.BackoffLimit,
		},
	}
	_, err = jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	check(err)
}

func Deployment(k *KubernetesClient) {
	deploymentsClient := k.client.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{}
	fmt.Sprintf("%s %s", deploymentsClient, deployment)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseConfig(f string) (*JobConfig, error) {
	fmt.Printf("File: %s", f)
	file, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	var config JobConfig
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	fmt.Println(config)
	return &config, nil
}
