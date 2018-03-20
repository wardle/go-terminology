package server

import (
	"fmt"
	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

//go:generate protoc -I. -I.. --go_out=plugins=grpc:. server.proto

type myServer struct {
	svc *terminology.Svc
}

// RunRPCServer runs the RPC server
func RunRPCServer(svc *terminology.Svc, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)
	RegisterSnomedCTServer(server, &myServer{svc: svc})
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
	return nil
}

func (ss *myServer) GetConcept(ctx context.Context, id *SctID) (*snomed.Concept, error) {
	return ss.svc.GetConcept(id.Identifier)
}

var _ SnomedCTServer = (*myServer)(nil)
