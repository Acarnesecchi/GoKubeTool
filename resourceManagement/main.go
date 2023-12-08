package main

import (
	"flag"
	"fmt"
)

func main() {
	fp := flag.String("file", "default", "path of file to examine")
	inCluster := flag.Bool("in-cluster", false, "Connection mode. Set to false if it is executed from a pod inside the cluster")
	flag.Parse()
	option := flag.Arg(0)

	client := &KubernetesClient{}
	if !*inCluster {
		outClusterConnect(client)
	} else {
		inClusterConnect(client)
	}

	switch option {
	case "resetdb":
		for i := 0; i < 3; i++ {
			fmt.Printf("Resetting DB using %s file \n", *fp)
			val := resetDB(client, *fp)
			if val {
				fmt.Println("Job ended succesfully.")
				break
			}
			fmt.Printf("Job ended in a failure. Retrying again up to %d more times\n", 3-i)
		}
	case "fetch":
		fmt.Println("Doing nothing")
	case "":
		panic(("Empty argument"))
	default:
		panic(fmt.Sprintf("Invalid argument: %s", option))
	}
	//fetch(client)
}
