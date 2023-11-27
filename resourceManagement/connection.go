package main

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubernetesClient struct {
	client *kubernetes.Clientset
}

func inClusterConnect(k *KubernetesClient) {
	// Configures the client to use your local kubeconfig file.
	var kubeconfig string
	defaultKubeConfigPath := filepath.Join("/root", ".kube", "config") //you'll have to deal with file permissions!
	if home := homedir.HomeDir(); home != "" {
		homeKubeConfigPath := filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(homeKubeConfigPath); err == nil {
			kubeconfig = homeKubeConfigPath
		} else {
			kubeconfig = defaultKubeConfigPath
		}
	} else {
		kubeconfig = defaultKubeConfigPath
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	k.client = clientset
}

func outClusterConnect(k *KubernetesClient) {
	// Configures the client to use your local kubeconfig file.
	var kubeconfig string
	defaultKubeConfigPath := filepath.Join("/root", ".kube", "config") //you'll have to deal with file permissions!
	if home := homedir.HomeDir(); home != "" {
		homeKubeConfigPath := filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(homeKubeConfigPath); err == nil {
			kubeconfig = homeKubeConfigPath
		} else {
			kubeconfig = defaultKubeConfigPath
		}
	} else {
		kubeconfig = defaultKubeConfigPath
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	k.client = clientset // stores the client connection for later use
}
