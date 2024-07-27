package server

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"

	crv1alpha1 "github.com/crossplane/crossplane-runtime/apis/proto/v1alpha1"
)

type Server struct {
	crv1alpha1.UnimplementedChangeLogServiceServer
}

func (s *Server) SendChangeLog(ctx context.Context, entry *crv1alpha1.ChangeLogEntry) (*crv1alpha1.ChangeLogResponse, error) {
	b, err := protojson.Marshal(entry)
	if err != nil {
		return &crv1alpha1.ChangeLogResponse{
			Success: false,
			Message: "Failed to marshall entry",
		}, err
	}

	// write the change log entry to stdout
	fmt.Printf("%s\n", string(b))

	return &crv1alpha1.ChangeLogResponse{
		Success: true,
		Message: "Change log entry received",
	}, nil
}
