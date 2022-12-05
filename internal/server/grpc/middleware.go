package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var ErrPeerFromContext = status.Error(codes.Internal, "getting peer fail")

func loggingMiddleware(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (response interface{}, err error) {
		p, ok := peer.FromContext(ctx)
		if !ok {
			logger.Error(ErrPeerFromContext.Error())

			return response, ErrPeerFromContext
		}

		remoteAddr := p.Addr.String()
		start := time.Now()
		response, err = handler(ctx, req)

		logger.Info(fmt.Sprintf(`%s [%s] %s %s (%.2fs) gRPC-Call"`,
			remoteAddr,
			start.Format("2006-01-02 15:04:05"),
			info.FullMethod,
			req,
			time.Since(start).Seconds(),
		))

		return response, err
	}
}
