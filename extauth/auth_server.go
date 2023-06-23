package main

import (
	"context"
	"log"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
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
			Status: &status.Status{
				Code: int32(code.Code_UNAUTHENTICATED),
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
		Status: &status.Status{
			Code: int32(code.Code_OK),
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
