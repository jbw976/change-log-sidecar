package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	changelogs "github.com/crossplane/crossplane-runtime/apis/changelogs/proto/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
)

type Server struct {
	changelogs.UnimplementedChangeLogServiceServer
}

func (s *Server) SendChangeLog(ctx context.Context, entry *changelogs.SendChangeLogRequest) (*changelogs.SendChangeLogResponse, error) {
	// Marshal the change log entry coming over the wire to JSON using the
	// protojson helper
	b, err := protojson.Marshal(entry)
	if err != nil {
		st := status.New(codes.Internal, errors.Wrap(err, "Failed to marshall input entry").Error())
		return &changelogs.SendChangeLogResponse{}, st.Err()
	}

	// Now unmarshal those bytes into a map so we can inject fields/data that
	// are server side only
	var entryData map[string]interface{}
	if err := json.Unmarshal(b, &entryData); err != nil {
		st := status.New(codes.Internal, errors.Wrap(err, "Failed to unmarshal input entry").Error())
		return &changelogs.SendChangeLogResponse{}, st.Err()
	}

	// Include a server side current timestamp into the data
	entryData["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	// Marshal the data back to JSON bytes once again so we can write the final
	// product to the output
	b, err = json.Marshal(entryData)
	if err != nil {
		st := status.New(codes.Internal, errors.Wrap(err, "Failed to marshal final entry").Error())
		return &changelogs.SendChangeLogResponse{}, st.Err()
	}

	// write the final change log entry to stdout
	fmt.Printf("%s\n", string(b))

	return &changelogs.SendChangeLogResponse{}, nil
}
