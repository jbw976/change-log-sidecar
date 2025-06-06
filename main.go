package main

import (
	"errors"
	"io/fs"
	"log"
	"net"
	"os"

	changelogs "github.com/crossplane/crossplane-runtime/apis/changelogs/proto/v1alpha1"

	"github.com/jbw976/change-log-sidecar/server"

	"google.golang.org/grpc"
)

func main() {
	socketPath := "/var/run/change-logs/change-logs.sock"

	if err := os.Remove(socketPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("failed to remove existing unix domain socket at %s: %+v", socketPath, err)
	}

	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("failed to listen to unix domain socket at %s: %+v", socketPath, err)
	}

	s := server.Server{}
	grpcServer := grpc.NewServer()
	changelogs.RegisterChangeLogServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
}
