package kube

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubernetesClient struct {
	client *kubernetes.Clientset
}

func deleteConfigFile() error {
	cwd, err := os.Getwd() // Getting current working dir
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}

	configFilePath := filepath.Join(cwd, "mk8s-config")
	if _, err := os.Stat(configFilePath); err == nil {
		// File exists, so delete it
		if err := os.Remove(configFilePath); err != nil {
			return fmt.Errorf("error deleting 'mk8s-config': %w", err)
		}
		fmt.Println("'mk8s-config' successfully deleted")
	} else if os.IsNotExist(err) {
		fmt.Println("no config to remove")
	} else {
		return fmt.Errorf("error checking 'mk8s-config': %w", err)
	}

	return nil
}

func getKubeConfig(isMicrok8s bool) string {
	if isMicrok8s {
		cmd := exec.Command("microk8s", "config")
		cmdOutput, err := cmd.Output()
		if err != nil {
			panic("Error executing microk8s config command")
		}

		configPath := filepath.Join(".", "mk8s-config")
		err = os.WriteFile(configPath, cmdOutput, 0644)
		if err != nil {
			panic("Error writing mk8s-config file")
		}
		return configPath
	} else {
		if home := homedir.HomeDir(); home != "" {
			homeKubeConfigPath := filepath.Join(home, ".kube", "config")
			if _, err := os.Stat(homeKubeConfigPath); err == nil {
				return homeKubeConfigPath
			}
		}
		return filepath.Join(".", "config")
	}
}

func inClusterConnect(k *KubernetesClient, useMicrok8s bool) {
	// Configures the client to use your local kubeconfig file.
	kubeconfig := getKubeConfig(useMicrok8s) // Getting kubeconfig path

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic("Error looking for kubeconfig")
		}
	}

	k.client, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic("Could not start clientset. Revise your kubeconfig roles and permissiones")
	}
	// Removing config file
	if err := deleteConfigFile(); err != nil {
		fmt.Println("Error:", err)
	}
}

func outClusterConnect(k *KubernetesClient, useMicrok8s bool) {
	// Configures the client to use your local kubeconfig file.
	kubeconfig := getKubeConfig(useMicrok8s) // Getting kubeconfig path
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic("Error looking for kubeconfig")
	}

	k.client, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic("Could not start clientset. Revise your kubeconfig roles and permissiones")
	}
	// Removing config file
	if err := deleteConfigFile(); err != nil {
		fmt.Println("Error:", err)
	}
}
