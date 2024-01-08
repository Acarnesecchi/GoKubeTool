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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	JobName        string            `yaml:"jobname"`
	Namespace      string            `yaml:"namespace"`
	ContainerImage string            `yaml:"containerImage"`
	Command        string            `yaml:"command"`
	RestartPolicy  string            `yaml:"restartPolicy"`
	BackoffLimit   int32             `yaml:"backoffLimit"`
	FailCondition  string            `yaml:"failCondition"`
	GracePeriod    int               `yaml:"gracePeriod"`
	EnvVariables   map[string]string `yaml:"envVariables"`
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
	var podName string
	config, err := parseConfig(f)
	check(err)
	podName, err = createJob(k, config)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			fmt.Printf("Job %s already exists. Do you want to delete it before proceeding? [y/N]\n", config.JobName)
			var r string
			fmt.Scanln(&r)
			r = strings.ToLower(strings.TrimSpace(r))
			if r == "y" || r == "yes" {
				deleteJob(k, config)
				time.Sleep(5 * time.Second)
				podName, err = createJob(k, config)
				check(err)
			} else {
				fmt.Println("Exiting session...")
				return true
			}
		} else {
			panic(fmt.Sprintf("Failed to create job: %v\n", err))
		}
	}

	podLogOpts := v1.PodLogOptions{Follow: true}
	req := k.client.CoreV1().Pods(config.Namespace).GetLogs(podName, &podLogOpts)

	success := make(chan bool)
	go func() {
		logStream, err := req.Stream(context.TODO())
		check(err)
		fmt.Println("Connected to pod. Retrieving logs...")
		defer logStream.Close()

		reader := bufio.NewReader(logStream)
		ticker := time.NewTicker(15 * time.Second) // Adjust the duration as needed
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C: // refresh the pod status every tick
				updatedPod, err := k.client.CoreV1().Pods(config.Namespace).Get(context.TODO(), podName, metav1.GetOptions{})
				check(err)
				if updatedPod.Status.Phase == v1.PodSucceeded {
					success <- true
					fmt.Println("Job ended succesfully.")
					return
				} else if updatedPod.Status.Phase == v1.PodFailed {
					success <- false
					return
				}
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						continue
					} else {
						check(err)
					}
				}
				fmt.Print(line)
				if strings.Contains(line, config.FailCondition) {
					success <- false
					return
				}
			}
		}
	}()

	val := <-success
	if config.GracePeriod > 0 {
		fmt.Printf("Waiting %d seconds before attempting pod deletion\n", config.GracePeriod)
	}
	deleteJob(k, config)
	return val
}

func Deployment(k *KubernetesClient) {
	deploymentsClient := k.client.AppsV1().Deployments(v1.NamespaceDefault)
	deployment := &appsv1.Deployment{}
	fmt.Printf("%s %s", deploymentsClient, deployment)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseConfig(f string) (*Config, error) {
	file, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func createJob(k *KubernetesClient, c *Config) (string, error) {
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

	var envVars []v1.EnvVar
	for key, value := range c.EnvVariables {
		envVars = append(envVars, v1.EnvVar{
			Name:  key,
			Value: value,
		})
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
							Env:     envVars,
						},
					},
					RestartPolicy: restartPolicy,
				},
			},
			BackoffLimit: &c.BackoffLimit,
		},
	}
	job, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	labelSelector := fmt.Sprintf("job-name=%s", job.ObjectMeta.Name)
	watcher, err := k.client.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	check(err)
	var podName string
	for event := range watcher.ResultChan() {
		pod, ok := event.Object.(*v1.Pod)
		if !ok {
			return "", fmt.Errorf("unexpected type")
		}
		if pod.Status.Phase == v1.PodRunning {
			podName = pod.Name
			fmt.Println("Job started")
			break
		}
	}

	// timeout := time.After(30 * time.Minute) // Example timeout
	// for {
	// 	select {
	// 	case event := <-watcher.ResultChan():
	// 		e, ok := event.Object.(*v1.Event)
	// 		if !ok {
	// 			return fmt.Errorf("unexpected type")
	// 		}

	// 		// Check for specific events like image pull failures or PVC issues
	// 		if strings.Contains(e.Message, "Failed to pull image") {
	// 			// Handle image pull failure
	// 			return fmt.Errorf("image pull failure: %s", e.Message)
	// 		} else if strings.Contains(e.Message, "persistentvolumeclaim \"<PVC_NAME>\" not found") {
	// 			// Handle PVC binding issue
	// 			return fmt.Errorf("PVC binding issue: %s", e.Message)
	// 		}

	// 	case <-timeout:
	// 		// Handle timeout scenario
	// 		return fmt.Errorf("operation timed out")
	// 	}
	// }

	return podName, err
}

func deleteJob(k *KubernetesClient, c *Config) {
	fmt.Println("Attemping pod deletion...")
	propagationPolicy := metav1.DeletePropagationBackground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	}

	err := k.client.BatchV1().Jobs(c.Namespace).Delete(context.TODO(), c.JobName, deleteOptions)
	check(err)

	fmt.Println("Job and its associated pods are being deleted.")
}
