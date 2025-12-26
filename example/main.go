package main

import (
	"context"
	"errors"
	"log"
	"net"
	"strings"

	"github.com/casnerano/protoc-gen-go-rbac/example/pb"
	"github.com/casnerano/protoc-gen-go-rbac/pkg/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OneServer struct {
	pb.UnimplementedExampleServiceOneServer
}

func (s *OneServer) Stats(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *OneServer) Update(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	authResolver := func(ctx context.Context) (*interceptor.AuthContext, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("incoming md error")
		}

		authContext := interceptor.AuthContext{}

		if authorization, exists := md["authorization"]; exists {
			authContext.Authenticated = len(authorization) == 1
		}

		if roles, exists := md["roles"]; exists && len(roles) == 1 {
			authContext.Roles = strings.Split(roles[0], ",")
		}

		return &authContext, nil
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.RbacUnary(authResolver)))

	pb.RegisterExampleServiceOneServer(server, &OneServer{})

	reflection.Register(server)

	if err = server.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
