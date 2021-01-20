package main

import (
	"context"
	"log"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/gogo/googleapis/google/rpc"
)

type authorizationServer struct{}

func (a *authorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	authHeader, _ := req.Attributes.Request.Http.Headers["authorization"]

	log.Println("Request received")
	log.Println("Auth header: " + authHeader)

	var userID string
	switch {
	case authHeader == "Bearer user-1-token":
		userID = "user-1"
	case authHeader == "Bearer user-2-token":
		userID = "user-2"
	default:
		return &auth.CheckResponse{
			Status: &rpc.Status{
				Code: int32(rpc.UNAUTHENTICATED),
			},
			HttpResponse: &auth.CheckResponse_DeniedResponse{
				DeniedResponse: &auth.DeniedHttpResponse{
					Status: &envoy_type.HttpStatus{
						Code: envoy_type.StatusCode_Unauthorized,
					},
					Body: "You need one of our super secure tokens to access",
				},
			},
		}, nil
	}

	return &auth.CheckResponse{
		Status: &rpc.Status{
			Code: int32(rpc.OK),
		},
		HttpResponse: &auth.CheckResponse_OkResponse{
			OkResponse: &auth.OkHttpResponse{
				// inject a header that can be used for future rate limiting
				// and backend processing
				Headers: []*core.HeaderValueOption{
					{
						Header: &core.HeaderValue{
							Key:   "x-user-id",
							Value: userID,
						},
					},
				},
			},
		},
	}, nil
}
