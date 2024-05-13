package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/CarlFlo/steamWorkshopDownloader/protos/workshopParser"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	validBearerToken = "some-secret-token"
	port             = 9000
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func Launch() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port '%d'. Reason: '%v'\n", port, err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken),
	}

	grpcServer := grpc.NewServer(opts...)
	workshopParserService := &MyWorkshopParserServer{}

	workshopParser.RegisterWorkshopParserServer(grpcServer, workshopParserService)

	reflection.Register(grpcServer)

	fmt.Printf("Server running and listening on port %d\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port '%d'. Reason: '%v'", port, err)
	}
}

func ensureValidToken(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	if !isAuthorized(metadata["authorization"]) {
		return nil, errInvalidToken
	}
	// Token is valid, so we continue
	return handler(ctx, req)
}

func isAuthorized(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == validBearerToken
}
