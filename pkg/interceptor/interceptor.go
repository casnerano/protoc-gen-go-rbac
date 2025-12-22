package interceptor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const msgAccessDenied = "access denied"

type checkable interface {
	CheckAccess(fullMethod string, roles []string) bool
}

type roleGetter interface {
	Roles() []string
}

type options struct {
	debug bool
}

type Option func(*options)

func WithDebug() Option {
	return func(options *options) {
		options.debug = true
	}
}

func RolesAccessor(roleGetter roleGetter, opts ...Option) grpc.UnaryServerInterceptor {
	rbacOptions := &options{}

	for _, opt := range opts {
		opt(rbacOptions)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if rbacOptions.debug {
			slog.DebugContext(ctx, "RBAC: interceptor invoked",
				slog.String("method", info.FullMethod),
				slog.Any("user_roles", roleGetter.Roles()),
			)
		}

		switch v := info.Server.(type) {
		case checkable:
			hasAccess := v.CheckAccess(info.FullMethod, roleGetter.Roles())

			if rbacOptions.debug {
				slog.DebugContext(ctx, "RBAC: the rules for the method access",
					slog.Bool("has_access", hasAccess),
				)
			}

			if hasAccess {
				return handler(ctx, req)
			}
		}

		if rbacOptions.debug {
			slog.DebugContext(ctx, "RBAC: rules for the method are not defined")
		}

		return nil, status.Error(codes.PermissionDenied, msgAccessDenied)
	}
}
