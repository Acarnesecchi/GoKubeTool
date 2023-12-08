package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
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
	namespaces, err := k.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	check(err)

	fmt.Println("Namespaces in the cluster:")
	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
	}
}

func resetDB(k *KubernetesClient, f string) bool {
	config, err := parseConfig(f)
	check(err)
	err = createJob(k, config)
	check(err)
	time.Sleep(5 * time.Second)

	var pod corev1.Pod
	found := false
	for !found {
		pods, err := k.client.CoreV1().Pods(config.Namespace).List(context.TODO(),
			metav1.ListOptions{
				LabelSelector: "job-name=" + config.JobName,
			})
		check(err)

		for _, p := range pods.Items {
			if strings.HasPrefix(p.Name, config.JobName) {
				pod = p
				found = true
				break
			}
		}
	}
	podLogOpts := corev1.PodLogOptions{Follow: true}
	req := k.client.CoreV1().Pods(config.Namespace).GetLogs(pod.Name, &podLogOpts)

	success := make(chan bool)
	go func() {
		logStream, err := req.Stream(context.TODO())
		check(err)
		fmt.Println("Connected to pod. Retrieving logs...")
		defer logStream.Close()

		reader := bufio.NewReader(logStream)
		ticker := time.NewTicker(5 * time.Second) // Adjust the duration as needed
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Refresh the pod object to get the latest status
				updatedPod, err := k.client.CoreV1().Pods(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
				check(err)
				if updatedPod.Status.Phase == corev1.PodSucceeded {
					success <- true
					return
				} else if updatedPod.Status.Phase == corev1.PodFailed {
					success <- false
					return
				}
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						continue // wait until pod status gets updated
					} else {
						check(err)
					}
				}
				fmt.Print(line)
			}
		}
	}()

	val := <-success
	fmt.Println(val)

	fmt.Println("Job ended")
	fmt.Println("Attemping pod deletion...")
	deleteJob(k, config)
	return val
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
	file, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	var c JobConfig
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func createJob(k *KubernetesClient, c *JobConfig) error {
	var cmd []string
	var restartPolicy v1.RestartPolicy
	namespace := "default"
	if c.Namespace != "" {
		namespace = c.Namespace
	}
	if c.Command != "" {
		cmd = strings.Fields(c.Command) //entrypoint formatted ["", ""]
	}
	if c.RestartPolicy != "" {
		restartPolicy = v1.RestartPolicy(c.RestartPolicy)
	} else {
		restartPolicy = v1.RestartPolicyNever
	}
	jobs := k.client.BatchV1().Jobs(namespace)
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.JobName,
			Namespace: c.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    c.JobName,
							Image:   c.ContainerImage,
							Command: cmd,
						},
					},
					RestartPolicy: restartPolicy,
				},
			},
			BackoffLimit: &c.BackoffLimit,
		},
	}
	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	return err
}

func deleteJob(k *KubernetesClient, c *JobConfig) {
	propagationPolicy := metav1.DeletePropagationBackground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	}

	err := k.client.BatchV1().Jobs(c.Namespace).Delete(context.TODO(), c.JobName, deleteOptions)
	check(err)

	fmt.Println("Job and its associated pods are being deleted.")
}
