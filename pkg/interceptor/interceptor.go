package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const msgAccessDenied = "access denied"

type rolesAccessChecker interface {
	CheckRolesAccess(fullMethod string, roles []string) bool
}

type rolesProvider interface {
	Roles() []string
}

type rolesAccessorOptions struct{}

type Option func(*rolesAccessorOptions)

func RolesAccessor(rolesProvider rolesProvider, opts ...Option) grpc.UnaryServerInterceptor {
	options := &rolesAccessorOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		switch v := info.Server.(type) {
		case rolesAccessChecker:
			if v.CheckRolesAccess(info.FullMethod, rolesProvider.Roles()) {
				return handler(ctx, req)
			}
		}

		return nil, status.Error(codes.PermissionDenied, msgAccessDenied)
	}
}
