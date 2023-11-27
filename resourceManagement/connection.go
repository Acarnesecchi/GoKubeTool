package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubernetesClient struct {
	client *kubernetes.Clientset
}

func inClusterConnect() {
	// Configures the client to use your local kubeconfig file.
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = filepath.Join("root", ".kube", "config")
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

	// Fetches the list of namespaces.
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Namespaces in the cluster:")
	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
	}
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

	// Fetches the list of namespaces.
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Namespaces in the cluster:")
	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
	}
}
