package main

import (
	"flag"
	"fmt"
)

func main() {
	fp := flag.String("file", "default", "path of file to examine")
	inCluster := flag.Bool("in-cluster", false, "Connection mode. Set to false if it is executed from a pod inside the cluster")
	flag.Parse()
	fmt.Println("path:", *fp)
	option := flag.Arg(0)
	fmt.Println("Option:", option)

	client := &KubernetesClient{}
	if !*inCluster {
		outClusterConnect(client)
	} else {
		inClusterConnect(client)
	}

	switch option {
	case "resetdb":
		fmt.Sprintf("Resetting DB using %s file", *fp)
		resetDB(client, *fp)
	case "fetch":
		fmt.Println("Doing nothing")
	case "":
		panic(fmt.Sprint("Empty argument"))
	default:
		panic(fmt.Sprintf("Invalid argument: %s", option))
	}
	fetch(client)
}
