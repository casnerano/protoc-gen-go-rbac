package interceptor

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const msgAccessDenied = "access denied"

var ErrRoleExtractor = errors.New("failed role extractor")

type checkable interface {
	CheckAccess(fullMethod string, roles []string) bool
}

type AuthContext struct {
	Authenticated bool
	Roles         []string
}

type AuthContextResolver func(ctx context.Context) (*AuthContext, error)

func RolesAccessor(authContextResolver AuthContextResolver, opts ...Option) grpc.UnaryServerInterceptor {
	rbacOptions := &options{}

	for _, opt := range opts {
		opt(rbacOptions)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		authContext, err := authContextResolver(ctx)
		if err != nil {
			return nil, errors.Join(ErrRoleExtractor, err)
		}

		if rbacOptions.debug {
			slog.DebugContext(ctx, "RBAC: interceptor invoked",
				slog.String("method", info.FullMethod),
				slog.Any("user_roles", authContext.Roles),
			)
		}

		switch v := info.Server.(type) {
		case checkable:
			hasAccess := v.CheckAccess(info.FullMethod, authContext.Roles)

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
