package main

import kube "github.com/Acarnesecchi/GoKubeTool/resourceManagement"

func main() {
	kube.StartServer()
}

//option string, inCluster bool, useMicrok8s bool, fp string /opt/devtools/config.yaml
