package main

func main() {
	mode := "OutOfCluster" // mode should be set from the CLI or env variable
	client := &KubernetesClient{}
	if mode == "OutOfCluster" {
		outClusterConnect(client)
	} else if mode == "InCluster" {
		inClusterConnect(client)
	}
	fetch(client)
}
