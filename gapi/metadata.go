package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (s *Server) extractMetaData(ctx context.Context) *Metadata {
	mtdt := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {

		if t := md.Get("content-type"); len(t) > 0 {
			if userAgents := md.Get("user-agent"); len(userAgents) > 0 {
				if p, ok := peer.FromContext(ctx); ok {
					mtdt.ClientIP = p.Addr.String()
				}
				mtdt.UserAgent = userAgents[0]
			}
		}
		if t := md.Get("grpcgateway-content-type"); len(t) > 0 {
			if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
				mtdt.UserAgent = userAgents[0]
			}
			if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
				mtdt.ClientIP = clientIPs[0]
			}
		}

	}
	return mtdt
}
