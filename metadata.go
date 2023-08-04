package auth_library

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

type Metadata struct {
	UserUuid  string
	ProfileId string
}

func GetMetadata(ctx context.Context) (*Metadata, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &Metadata{}, fmt.Errorf("no data in incoming context")
	}
	meta := &Metadata{}

	for key, values := range md {
		for _, value := range values {
			switch key {
			case "user_uuid":
				meta.UserUuid = value
			case "profile_id":
				meta.ProfileId = value
			default:
				continue
			}
		}
	}
	return meta, nil
}
