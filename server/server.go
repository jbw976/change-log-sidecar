package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	changelogs "github.com/crossplane/crossplane-runtime/apis/changelogs/proto/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
)

type Server struct {
	changelogs.UnimplementedChangeLogServiceServer
}

func (s *Server) SendChangeLog(ctx context.Context, req *changelogs.SendChangeLogRequest) (*changelogs.SendChangeLogResponse, error) {
	if req == nil || req.Entry == nil {
		st := status.New(codes.Internal, "Request and change logs entry must not be nil")
		return &changelogs.SendChangeLogResponse{}, st.Err()
	}

	if req.Entry.Timestamp != nil {
		// We only care about resolution of the timestamps to seconds, so discard
		// the nanoseconds.
		req.Entry.Timestamp.Nanos = 0
	}

	// Marshal the change log entry coming over the wire to JSON using the
	// protojson helper
	b, err := protojson.Marshal(req.Entry)
	if err != nil {
		st := status.New(codes.Internal, errors.Wrap(err, "Failed to marshall input entry").Error())
		return &changelogs.SendChangeLogResponse{}, st.Err()
	}

	// write the final change log entry to stdout
	fmt.Printf("%s\n", string(b))

	return &changelogs.SendChangeLogResponse{}, nil
}
