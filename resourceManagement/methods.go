package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func fetch(k *KubernetesClient) {
	// Fetches the list of namespaces.
	namespaces, err := k.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Namespaces in the cluster:")
	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
	}
}
