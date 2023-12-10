package main

import (
	"flag"
	"fmt"
	"os"
)

const max int = 3

func main() {
	fp := flag.String("f", "/opt/devtools/config.yaml", "path of file to examine")
	inCluster := flag.Bool("ioc", false, "Connection mode. Set to false if it is executed from a pod inside the cluster")
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
		val := true
		fmt.Printf("Resetting DB using %s file \n", *fp)
		for i := 0; i < max; i++ {
			if !val {
				fmt.Printf("Job ended in a failure. Retrying again up to %d more times\n", max-i)
			}
			val = resetDB(client, *fp)
			if val {
				fmt.Println("Job ended succesfully.")
				os.Exit(0)
			}
		}
		fmt.Printf("Job failed %d times. Check everything is working and try again\n", max)
	case "fetch":
		fmt.Println("Doing nothing")
		fetch(client)
	case "":
		panic(("Empty argument"))
	default:
		panic(fmt.Sprintf("Invalid argument: %s", option))
	}
}
