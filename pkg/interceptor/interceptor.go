package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const msgAccessDenied = "access denied"

type accessChecker interface {
	CheckAccess(fullMethod string, roles []string) bool
}

type roleGetter interface {
	Roles() []string
}

type rolesAccessorOptions struct{}

type Option func(*rolesAccessorOptions)

func RolesAccessor(roleGetter roleGetter, opts ...Option) grpc.UnaryServerInterceptor {
	options := &rolesAccessorOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		switch v := info.Server.(type) {
		case accessChecker:
			if v.CheckAccess(info.FullMethod, roleGetter.Roles()) {
				return handler(ctx, req)
			}
		}

		return nil, status.Error(codes.PermissionDenied, msgAccessDenied)
	}
}
