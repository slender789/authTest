package auth_library

import (
	"context"
	"fmt"

	fgaSdk "github.com/openfga/go-sdk"
	"google.golang.org/grpc"
)

type Config struct {
	Rules         map[string]func(interface{}) []Rules
	OpenFGAConfig *fgaSdk.Configuration
}

type RequestOwner struct {
	User         string
	UsersProfile string
}

type CheckResponse struct {
	Allowed    bool   `json:"allowed"`
	Resolution string `json:"resolution"`
}

func AuthLibraryInterceptor(config Config, openfgaClient PermissionApi) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Printf("Received request for endpoint: %v", info.FullMethod)

		md, err := GetMetadata(ctx)
		if err != nil {
			fmt.Printf("Metadata error: %v", err)
			return nil, err
		}

		fmt.Printf("Request Owner %v %v", md.UserUuid, md.ProfileId)

		requestOwner := RequestOwner{
			User:         md.UserUuid,
			UsersProfile: md.ProfileId,
		}

		rules := config.Rules[info.FullMethod]
		for _, rule := range rules(req) {
			err := rule.Authorize(openfgaClient, requestOwner)
			if err != nil {
				fmt.Printf("Denied for %v: %v", info.FullMethod, err)
				return nil, err
			}
		}
		fmt.Printf("Authorized %v", info.FullMethod)

		resp, err := handler(ctx, req)
		if err != nil {
			fmt.Printf("Handler error for %v: %v", info.FullMethod, err)
			return nil, err
		}
		return resp, nil
	}
}
