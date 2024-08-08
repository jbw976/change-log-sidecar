package server

import (
	"encoding/json"
	"fmt"
	"io"
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

func (s *Server) SendChangeLog(stream changelogs.ChangeLogService_SendChangeLogServer) error {
	callCount := 0
	for {
		callCount++
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&changelogs.SendChangeLogResponse{})
		}
		if err != nil {
			return err
		}

		if req != nil {
			st := status.New(codes.Internal, fmt.Sprintf("fake error %d from jared", callCount))
			return st.Err()
		}

		if req == nil || req.Entry == nil {
			st := status.New(codes.Internal, "Request and change logs entry must not be nil")
			return st.Err()
		}

		// Marshal the change log entry coming over the wire to JSON using the
		// protojson helper
		b, err := protojson.Marshal(req.Entry)
		if err != nil {
			st := status.New(codes.Internal, errors.Wrap(err, "Failed to marshall input entry").Error())
			return st.Err()
		}

		// Now unmarshal those bytes into a map so we can inject fields/data that
		// are server side only
		var entryData map[string]interface{}
		if err := json.Unmarshal(b, &entryData); err != nil {
			st := status.New(codes.Internal, errors.Wrap(err, "Failed to unmarshal input entry").Error())
			return st.Err()
		}

		// Include a server side current timestamp into the data
		entryData["timestamp"] = time.Now().UTC().Format(time.RFC3339)

		// Marshal the data back to JSON bytes once again so we can write the final
		// product to the output
		b, err = json.Marshal(entryData)
		if err != nil {
			st := status.New(codes.Internal, errors.Wrap(err, "Failed to marshal final entry").Error())
			return st.Err()
		}

		// write the final change log entry to stdout
		fmt.Printf("%s\n", string(b))
	}
}
