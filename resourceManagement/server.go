package kube

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/Acarnesecchi/GoKubeTool/proto"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
)

type devToolsServer struct {
	pb.UnimplementedDevToolsServiceServer
	k8sClient *kubernetes.Clientset
}

func NewDevToolsServer(inCluster, useMicrok8s bool) (*devToolsServer, error) {
	server := &devToolsServer{}
	var err error
	if !inCluster {
		err = outClusterConnect(server, useMicrok8s)
	} else {
		err = inClusterConnect(server.k8sClient, useMicrok8s)
	}
	return server, err
}

func (s *devToolsServer) Fetch(ctx context.Context, fr *pb.FetchRequest) (*pb.FetchResponse, error) {
	if (fr.Resource == pb.ResourceType_JOB || fr.Resource == pb.ResourceType_POD) && fr.Namespace == "" {
		return nil, fmt.Errorf("namespace is required for fetching this resource")
	}
	resources, err := fetchResources(s.k8sClient, fr.Resource, fr.Namespace)
	if err != nil {
		return nil, err
	}
	return &pb.FetchResponse{Resource: resources}, nil
}

func StartServer() {
	port := 5051
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to lister: %v", err)
	}
	grpcServer := grpc.NewServer()
	newServer, _ := NewDevToolsServer(false, false)
	pb.RegisterDevToolsServiceServer(grpcServer, newServer)
	grpcServer.Serve(lis)
}
